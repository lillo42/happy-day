import {DatePipe} from "@angular/common";
import {Component, Inject, OnInit} from '@angular/core';
import {FormArray, FormBuilder, FormGroup, Validators} from "@angular/forms";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {BehaviorSubject, debounceTime, Observable, switchMap, tap} from "rxjs";
import {Product, ProductSort} from "../models/product";
import {ReservationService} from "../http-clients/reservation.service";
import {ProductService} from "../http-clients/product.service";
import {CustomerService} from "../http-clients/customer.service";
import {PaymentInstallment, Product as RProduct, Quote} from "../models/reservation";
import {Phone} from "../models/customer";

@Component({
  selector: 'app-reservation',
  templateUrl: './reservation.component.html',
  styleUrls: ['./reservation.component.scss'],
  providers: [DatePipe]
})
export class ReservationComponent implements OnInit {
  CREATE_TITLE = "Criação de um reserva";
  CHANGE_TITLE = "Atualização de um reserva";
  DELETE_TITLE = "Remover um reserva";

  formGroup: FormGroup;
  filterProducts = new BehaviorSubject<Product[]>([]);
  filterProducts$: Observable<Product[]>;

  constructor(private dialogRef: MatDialogRef<ReservationComponent>,
              @Inject(MAT_DIALOG_DATA) private data: ReservationData,
              private builder: FormBuilder,
              private reservationService: ReservationService,
              private productService: ProductService,
              private customerService: CustomerService,
              private datePipe: DatePipe) {
    this.filterProducts$ = this.filterProducts.asObservable();
    this.formGroup = builder.group({
      id: [{value: null, disabled: true}],
      products: builder.array([]),
      price: [{value: null, disabled: true}],
      discount: [{value: 0}, []],
      finalPrice: [{value: null, disabled: true}],
      comment: [null, []],
      deliveryAt: [null, [Validators.required]],
      deliveryBy: builder.array([]),
      pickupAt: [null, [Validators.required]],
      pickUpBy: builder.array([]),
      address: this.builder.group({
        street: [null, [Validators.required]],
        number: [null, [Validators.required]],
        complement: [null, []],
        neighborhood: [null, [Validators.required]],
        city: [null, [Validators.required]],
      }),
      customer: this.builder.group({
        id: [{value: null, disabled: true}, [Validators.required]],
        name: [null, [Validators.required]],
        comment: [null, []],
        phones: this.builder.array([]),
      }),
      paymentInstallments: this.builder.array([]),
      createdAt: [{value: null, disabled: true}],
      modifiedAt: [{value: null, disabled: true}],
    });
  }

  ngOnInit(): void {
    if (this.data.behavior === ReservationBehavior.Create) {
      this.addProduct();
      this.formGroup.patchValue({
        price: this.price,
        discount: 0,
        finalPrice: 0
      });
      return;
    }
  }

  isCreateMode(): boolean {
    return this.data.behavior === ReservationBehavior.Create;
  }

  isChangeMode(): boolean {
    return this.data.behavior === ReservationBehavior.Change;
  }

  isDeleteMode(): boolean {
    return this.data.behavior === ReservationBehavior.Delete;
  }

  // Products
  get products(): FormArray {
    return this.formGroup.controls["products"] as FormArray;
  }

  deleteProduct(index: number): void {
    this.products.removeAt(index);
    if(this.products.length === 0) {
      this.addProduct();
    }
  }

  addProduct(product: RProduct | null = null): void {
    this.products.push(this.builder.group({
      id: [{value: null, disabled: true}, []],
      name: [{value: null, disabled: !this.isCreateMode()}, [Validators.required]],
      quantity: [{value: null, disabled: !this.isCreateMode()}, [Validators.required, Validators.min(1)]],
    }));

    let productControl = this.products.controls[this.products.length - 1];
    if (product !== null) {
      this.productService.get(product.id)
        .pipe(tap(p => productControl.patchValue({name: p.name})))
        .subscribe();
    }

    productControl.get("name")!.valueChanges
      .pipe(
        debounceTime(1000),
        switchMap(name => this.productService.getAll(1, 1000, name, ProductSort.NameAsc)),
        tap(products => this.filterProducts.next(products.items))
      )
      .subscribe();
  }

  productSelected(product: Product, index: number): void {
    let productControl = this.products.controls[index];
    productControl.patchValue({id: product.id});
    this.filterProducts.next([]);
  }

  // Price
  get price(): number {
    return this.formGroup.get("price")!.getRawValue();
  }

  updateFinalPrice(): void {
    this.formGroup.get("finalPrice")!.setValue(this.price - this.formGroup.get("discount")!.getRawValue());
  }

  calculatePrice(): void {
    let quote = <Quote>{
      products: this.products.getRawValue().map((p: any) => {
        return <RProduct>{
          id: p.id,
          quantity: p.quantity
        }
      })
    };

    this.reservationService.quote(quote)
      .pipe(
        tap(res => this.formGroup.get("price")!.setValue(res.price)),
        tap(() => this.updateFinalPrice())
      )
      .subscribe();
  }

  // Phones
  get phones(): FormArray {
    return this.formGroup.get("customer")!.get("phones") as FormArray;
  }

  deletePhone(index: number): void {
    this.phones.removeAt(index);
    if(this.phones.length === 0) {
      this.addPhone();
    }
  }

  addPhone(phone: Phone | null = null): void {
    this.phones.push(this.builder.group({
      number: [{value: phone?.number, disabled: this.isDeleteMode() }, [Validators.required, Validators.pattern('[- +()0-9]+')]]
    }));
  }

  // Delivery and PickUp by
  get deliveryBy(): FormArray {
    return this.formGroup.get("deliveryBy") as FormArray;
  }

  get pickUpBy(): FormArray {
    return this.formGroup.get("pickUpBy") as FormArray;
  }

  addDeliveryBy(employee: string | null = null): void {
    this.addBy(this.deliveryBy, employee);
  }

  addPickUpBy(employee: string | null = null): void {
    this.addBy(this.pickUpBy, employee);
  }

  deleteDeliveryBy(index: number): void {
    this.deliveryBy.removeAt(index);
  }

  deletePickUpBy(index: number): void {
    this.pickUpBy.removeAt(index);
  }

  private addBy(array: FormArray, name: string | null) {
    array.push(this.builder.group({
      name: [{value: name, disabled: this.isDeleteMode()}, [Validators.required]]
    }));
  }

  // Payment Installments
  get paymentInstallments(): FormArray {
    return this.formGroup.get("paymentInstallments") as FormArray;
  }

  addPaymentInstallment(installment: PaymentInstallment | null = null): void {
    this.paymentInstallments.push(this.builder.group({
      at: [{value: installment?.at, disabled: this.isDeleteMode()}, [Validators.required]],
      amount: [{value: installment?.amount, disabled: this.isDeleteMode()}, [Validators.required, Validators.min(1)]],
      method: [{value: installment?.method, disabled: this.isDeleteMode()}, [Validators.required]],
    }));
  }

  deletePaymentInstallment(index: number): void {
    this.paymentInstallments.removeAt(index);
  }
}

export enum ReservationBehavior {
  Create,
  Change,
  Delete
}

export interface ReservationData {
  behavior: ReservationBehavior;
  id: string;
}
