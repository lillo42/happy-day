import {CommonModule, DatePipe} from '@angular/common';
import {HttpErrorResponse} from "@angular/common/http";
import {AfterViewInit, Component, computed, ElementRef, OnInit, signal, ViewChild} from "@angular/core";
import {AbstractControl, FormArray, FormBuilder, FormGroup, ReactiveFormsModule, Validators} from "@angular/forms";
import {$localize} from "@angular/localize/init";
import {MatSelectModule} from "@angular/material/select";
import {ActivatedRoute, Router} from "@angular/router";

import {MatAutocompleteModule} from "@angular/material/autocomplete";
import {MatButtonModule} from "@angular/material/button";
import {MatOptionModule} from "@angular/material/core";
import {MatDatepickerModule} from "@angular/material/datepicker";
import {MatFormFieldModule} from "@angular/material/form-field";
import {MatIconModule} from "@angular/material/icon";
import {MatInputModule} from "@angular/material/input";
import {MatSnackBar} from "@angular/material/snack-bar";
import {MatTooltipModule} from "@angular/material/tooltip";

import {NgxMaskDirective} from "ngx-mask";
import {debounceTime, of, switchMap} from "rxjs";

import {ProblemDetails} from "../common";

import {Customer, CustomersService} from "../customers.service";
import {
  Order,
  OrderCreateOrChange,
  OrderCustomer,
  OrderPayment,
  OrderProduct,
  OrderQuote,
  OrdersService
} from "../orders.service";
import {ProductsService} from "../products.service";

@Component({
  selector: 'app-order-details',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, MatFormFieldModule, MatInputModule, MatButtonModule, MatTooltipModule, MatAutocompleteModule, MatOptionModule, MatIconModule, NgxMaskDirective, MatDatepickerModule, MatSelectModule],
  templateUrl: './order-details.component.html',
  styleUrl: './order-details.component.scss'
})
export class OrderDetailsComponent implements OnInit, AfterViewInit {
  form: FormGroup;

  id = signal<string | null>(null);
  hasFound = signal(false);
  isNew = computed(() => this.id() === null);
  isLoading = signal(false);


  filteredProducts = signal<OrderProduct[]>([]);
  filteredCustomers = signal<OrderCustomer[]>([]);

  @ViewChild("productInput") productInput: ElementRef | null = null;
  @ViewChild("customerInput") customerInput: ElementRef | null = null;

  constructor(private activatedRoute: ActivatedRoute,
              private datePipe: DatePipe,
              private builder: FormBuilder,
              private router: Router,
              private productsService: ProductsService,
              private customersService: CustomersService,
              private ordersService: OrdersService,
              private snackBar: MatSnackBar) {
    this.form = this.builder.group({
      address: [null, [Validators.required, Validators.maxLength(1000)]],
      comment: [null, null],

      deliveryAt: [null, [Validators.required]],
      pickUpAt: [null, [Validators.required]],

      totalPrice: [{value: null, disabled: true}, [Validators.required, Validators.min(0)]],
      discount: [null, [Validators.min(0)]],
      finalPrice: [null, [Validators.required, Validators.min(0)]],

      customer: this.builder.group({
        id: [{value: null, disabled: true}, null],
        name: [null, [Validators.required]],
        phones: this.builder.array([]),
        comment: [null, null],
      }),

      payments: this.builder.array([]),
      products: this.builder.array([]),

      createAt: [{value: null, disabled: true}, null],
      updateAt: [{value: null, disabled: true}, null],
    });
  }

  get products(): FormArray {
    return this.form.get('products') as FormArray;
  }

  get customer(): FormGroup {
    return this.form.get('customer') as FormGroup;
  }

  get customerPhones(): FormArray {
    return this.form.get('customer.phones') as FormArray;
  }

  get discount(): number {
    return this.form.get('discount')!.value as number;
  }

  get totalPrice(): number {
    return this.form.get('totalPrice')!.value as number;
  }

  get payments(): FormArray {
    return this.form.get('payments') as FormArray;
  }

  ngOnInit(): void {
    this.activatedRoute.paramMap
      .pipe(switchMap(params => {
        const id = params.get('id');
        if (id !== null && id !== 'new') {
          this.id.set(id);
          return this.ordersService.getById(id);
        } else {
          const empty: Order = {
            id: '',
            address: '',
            comment: '',

            deliveryAt: new Date(),
            pickUpAt: new Date(),

            totalPrice: 0,
            discount: 0,
            finalPrice: 0,

            customer: {
              id: '',
              name: '',
              phones: [],
              comment: '',
            },
            products: [],
            payments: [],

            createAt: new Date(),
            updateAt: new Date(),
          };
          return of(empty);
        }
      }))
      .subscribe({
        next: order => {
          this.updateForm(order);
          this.hasFound.set(true);
        },
        error: error => {
          if (error instanceof HttpErrorResponse) {
            if (error.status == 404) {
              this.hasFound.set(false);
            } else {
              const problemDetails: ProblemDetails = JSON.parse(error.message);
              this.snackBar.open(`an unexpected error happen: ${problemDetails.message}`, 'OK', {duration: 10000});
            }
            return;
          }

          this.snackBar.open(`an unexpected error happen: ${error.toString()}`, 'OK', {duration: 10000});
        },
        complete: () => this.isLoading.set(false)
      })
    ;
  }

  ngAfterViewInit(): void {
    this.loadProducts();
    this.loadCustomers();
  }

  deleteProduct(index: number): void {
    this.products.removeAt(index);
    this.quote();
  }

  addProduct(product: OrderProduct): void {
    if (this.productInput !== null) {
      this.productInput.nativeElement.value = '';
    }

    this.products.push(this.builder.group({
      id: [{value: product?.id, disabled: true}, [Validators.required]],
      name: [{value: product?.name, disabled: true}, null],
      quantity: [product?.quantity, [Validators.required, Validators.min(1)]],
      price: [{value: product?.price, disabled: true}, null],
    }));
  }

  displayProduct(product: OrderProduct | null): string {
    return product?.name || '';
  }

  updateCustomer(customer: OrderCustomer | null): void {
    this.customer.patchValue({...customer});
    this.customerPhones.clear();
    customer?.phones.forEach(phone => this.addCustomerPhone(phone));
  }

  displayCustomer(customer: OrderCustomer | null): string {
    return customer?.name || '';
  }

  deleteCustomerPhone(index: number): void {
    this.customerPhones.removeAt(index);
  }

  addCustomerPhone(phone: string | null = null): void {
    this.customerPhones.push(this.builder.control(phone, [
      Validators.required,
      Validators.minLength(8),
      Validators.maxLength(11),
      Validators.pattern('[- +()0-9]+')
    ]));
  }

  deletePayment(index: number): void {
    this.payments.removeAt(index);
  }

  addPayment(payment: OrderPayment | null = null): void {
    this.payments.push(this.builder.group({
      amount: [payment?.amount, [Validators.required, Validators.min(0)]],
      at: [payment?.at, [Validators.required]],
      method: [payment?.method, [Validators.required]],
    }));
  }

  loadProducts(): void {
    if (this.productInput === null) {
      return;
    }

    const name = this.productInput.nativeElement.value;
    this.productsService.get(name, 0, 100)
      .pipe(debounceTime(1000))
      .subscribe({
        next: page => {
          if (page.items === null) {
            this.filteredProducts.set([]);
            return;
          }

          this.filteredProducts.set(page.items.map(prod => {
            return <OrderProduct>{
              id: prod.id,
              name: prod.name,
              quantity: 0,
              price: prod.price,
            }
          }));
        },
        error: err => this.snackBar.open(err.message, 'OK')
      });
  }

  loadCustomers(): void {
    if (this.customerInput === null) {
      return;
    }

    const name = this.customerInput.nativeElement.value;
    this.customersService.get(name, null, null, 0, 100)
      .pipe(debounceTime(1000))
      .subscribe({
        next: page => {
          if (page.items === null) {
            this.filteredCustomers.set([]);
            return;
          }

          this.filteredCustomers.set(page.items.map(customer => {
            return <OrderCustomer>{
              id: customer.id,
              name: customer.name,
              phones: customer.phones,
              comment: customer.comment,
            }
          }));
        },
        error: err => this.snackBar.open(err.message, 'OK')
      });
  }

  asFormGroup(form: AbstractControl | null): FormGroup {
    return form as FormGroup;
  }

  cancel(): Promise<boolean> {
    return this.router.navigateByUrl('/orders');
  }

  updateFinalPrice(): void {
    this.form.patchValue({...{finalPrice: this.totalPrice - this.discount}});
  }

  quote(): void {
    const order = <Order>{...this.form.getRawValue()}
    const req = <OrderQuote>{
      products: order.products
    };

    this.ordersService.quote(req)
      .subscribe({
        next: res => {
          this.form.patchValue({...{totalPrice: res.totalPrice}});
          this.updateFinalPrice();
        },
        error: err => this.handlerError(err)
      });
  }

  save(): void {
    this.form.markAllAsTouched();
    this.form.markAsDirty();
    if (this.form.invalid) {
      return;
    }

    const order = <Order>{
      ...this.form.getRawValue()
    };

    let customer$ = of(<Customer>{
      id: order.customer.id,
      name: order.customer.name,
      phones: order.customer.phones,
      comment: order.customer.comment,
    });

    if (order.customer.id === null || order.customer.id === '') {
      const customer = <Customer>{
        id: '',
        name: order.customer.name,
        phones: order.customer.phones,
        comment: order.customer.comment,
      };

      customer$ = this.customersService.create(customer);
    }

    customer$
      .pipe(switchMap(customer => {
        const req = <OrderCreateOrChange>{
          address: order.address,
          comment: order.comment,

          deliveryAt: order.deliveryAt,
          pickUpAt: order.pickUpAt,

          totalPrice: order.totalPrice,
          discount: order.discount,
          finalPrice: order.finalPrice,
          customer: order.customer,

          customerId: customer.id,
          products: order.products,
          payments: order.payments,
        };

        if (this.isNew()) {
          return this.ordersService.create(req);
        } else {
          return this.ordersService.update(this.id()!, req);
        }
      }))
      .subscribe({
        next: _ => this.router.navigateByUrl('/orders'),
        error: error => this.handlerError(error)
      });
  }

  private handlerError(error: HttpErrorResponse) {
    if (error.status === 400) {
      this.form.markAllAsTouched();
      this.form.markAsDirty();
      return;
    }

    if (error.status == 0) {
      this.snackBar.open($localize`an unexpected error happen: ${error.message}`, 'OK', {duration: 10000});
      return;
    }

    const problemDetails: ProblemDetails = error.error;
    if (problemDetails.type === 'order-address-is-empty') {
      this.form.get('address')!.setErrors({required: true});
    } else if (problemDetails.type === 'order-address-is-too-large') {
      this.form.get('address')!.setErrors({maxlength: true});
    } else if (problemDetails.type === '"order-delivery-at-is-invalid') {
      this.form.get('deliveryAt')!.setErrors({matStartDateInvalid: true});
      this.form.get('pickUpAt')!.setErrors({matStartDateInvalid: true});
    } else if (problemDetails.type === 'order-final-price-at-is-invalid') {
      this.form.get('finalPrice')!.setErrors({min: true});
    } else if (problemDetails.type === 'order-payment-value-is-invalid') {
      this.payments.controls.forEach(control => control.setErrors({min: true}));
    } else if (problemDetails.type === 'order-conflict') {
      this.snackBar.open($localize`order update conflict, please reload the page`, 'OK', {duration: 10000});
    } else {
      this.snackBar.open($localize`an unexpected error happen: ${problemDetails.message}`, 'OK', {duration: 10000});
    }
  }

  private updateForm(order: Order): void {
    this.form.patchValue({...order});

    this.products.clear();
    order.products.forEach(product => this.addProduct(product));

    this.payments.clear();
    order.payments.forEach(payment => this.addPayment(payment));

    this.form.get("createAt")!.setValue(this.datePipe.transform(order.createAt, 'dd/MM/yyyy HH:mm:ss'));
    this.form.get("updateAt")!.setValue(this.datePipe.transform(order.updateAt, 'dd/MM/yyyy HH:mm:ss'));
  }
}
