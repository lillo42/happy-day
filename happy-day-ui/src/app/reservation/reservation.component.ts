import {DatePipe} from "@angular/common";
import {Component, Inject, OnInit} from '@angular/core';
import {AbstractControl, FormArray, FormBuilder, FormGroup, Validators} from "@angular/forms";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {BehaviorSubject, debounceTime, Observable, switchMap, tap} from "rxjs";
import {Product, ProductSort} from "../models/product";
import {ReservationService} from "../http-clients/reservation.service";
import {ProductService} from "../http-clients/product.service";
import {CustomerService} from "../http-clients/customer.service";
import {DeliveryOrPickUp, PaymentInstallment, Product as RProduct, Quote, Reservation} from "../models/reservation";
import {Customer, CustomerSort, Phone} from "../models/customer";

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

  filteredProducts = new BehaviorSubject<Product[]>([]);
  filteredProducts$: Observable<Product[]>;

  filteredCustomers = new BehaviorSubject<Customer[]>([]);
  filteredCustomers$: Observable<Customer[]>;

  constructor(private dialogRef: MatDialogRef<ReservationComponent>,
              @Inject(MAT_DIALOG_DATA) private data: ReservationData,
              private builder: FormBuilder,
              private reservationService: ReservationService,
              private productService: ProductService,
              private customerService: CustomerService) {
    this.filteredProducts$ = this.filteredProducts.asObservable();
    this.filteredCustomers$ = this.filteredCustomers.asObservable();
    this.formGroup = this.builder.group({
      id: [{value: null, disabled: true}],
      products: this.builder.array([]),
      price: [{value: null, disabled: true}],
      discount: [{value: 0, disabled: this.isDeleteMode()}, []],
      finalPrice: [{value: null, disabled: true}],
      comment: [{value: null, disabled: this.isDeleteMode()}, []],
      deliveryAt: [{value: null, disabled: this.isDeleteMode()}, [Validators.required]],
      deliveryBy: this.builder.array([]),
      pickUpAt: [{value: null, disabled: this.isDeleteMode()}, [Validators.required]],
      pickUpBy: this.builder.array([]),
      address: this.builder.group({
        street: [{value: null, disabled: this.isDeleteMode()}, [Validators.required]],
        number: [{value: null, disabled: this.isDeleteMode()}, [Validators.required]],
        complement: [{value: null, disabled: this.isDeleteMode()}, []],
        neighborhood: [{value: null, disabled: this.isDeleteMode()}, [Validators.required]],
        postalCode: [{value: null, disabled: this.isDeleteMode()}, [Validators.required]],
        city: [{value: null, disabled: this.isDeleteMode()}, [Validators.required]],
      }),
      customer: this.builder.group({
        id: [{value: null, disabled: true}, [Validators.required]],
        name: [{value: null, disabled: this.isDeleteMode()}, [Validators.required]],
        comment: [{value: null, disabled: this.isDeleteMode()}, []],
        phones: this.builder.array([]),
      }),
      paymentInstallments: this.builder.array([]),
      createdAt: [{value: null, disabled: true}],
      modifiedAt: [{value: null, disabled: true}],
    });
  }

  ngOnInit(): void {
    this.formGroup.get("customer")!.get("name")!.valueChanges
      .pipe(
        debounceTime(1000),
        switchMap(name => this.customerService.getAll(1, 1000, name, CustomerSort.NameAsc)),
        tap(customers => this.filteredCustomers.next(customers.items))
      )
      .subscribe();

    if (this.data.behavior === ReservationBehavior.Create) {
      this.addProduct();
      this.formGroup.patchValue({
        price: this.price,
        discount: 0,
        finalPrice: 0
      });
      return;
    }

    this.reservationService.get(this.data.id)
      .pipe(
        tap(reservation => {
          this.formGroup.patchValue(reservation);
          this.formGroup.get("deliveryAt")!.setValue(new Date(reservation.delivery.at));
          reservation.delivery.by.forEach(by => this.addDeliveryBy(by));

          this.formGroup.get("pickUpAt")!.setValue(new Date(reservation.pickUp.at));
          reservation.pickUp.by.forEach(by => this.addPickUpBy(by));

          reservation.products.forEach(product => this.addProduct(product));

          reservation.customer.phones.forEach(phone => this.addPhone(phone));

        })
      )
      .subscribe();
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
      id: [{value: product?.id, disabled: true}, []],
      name: [{value: null, disabled: !this.isCreateMode()}, [Validators.required]],
      quantity: [{value: product?.quantity, disabled: !this.isCreateMode()}, [Validators.required, Validators.min(1)]],
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
        tap(products => this.filteredProducts.next(products.items))
      )
      .subscribe();
  }

  productSelected(product: Product, index: number): void {
    let productControl = this.products.controls[index];
    productControl.patchValue({id: product.id});
    this.filteredProducts.next([]);
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

  // Address
  get address(): AbstractControl {
    return this.formGroup.get("address")!;
  }

  // Customer
  get customer(): AbstractControl {
    return this.formGroup.get("customer")!;
  }

  customerSelected(customer: Customer): void {
    this.customer.patchValue({id: customer.id, comment: customer.comment});
    for(let i = 0; i < customer.phones.length; i++) {
      this.addPhone(customer.phones[i]);
    }

    this.filteredCustomers.next([]);
  }

  save(): void {
    this.formGroup.markAllAsTouched();
    this.formGroup.markAsDirty();
    if (this.formGroup.invalid) {
      return;
    }

    let reservation = <Reservation>{};
    reservation.id = this.data.id;
    reservation.price = this.price;
    reservation.discount = this.formGroup.get("discount")!.getRawValue();
    reservation.finalPrice = this.formGroup.get("finalPrice")!.getRawValue();
    reservation.comment = this.formGroup.get("comment")!.getRawValue();
    reservation.customer = this.formGroup.get("customer")!.getRawValue();
    reservation.address = this.formGroup.get("address")!.getRawValue();
    reservation.products = this.products.getRawValue();
    reservation.paymentInstallments = this.paymentInstallments.getRawValue();
    reservation.delivery = <DeliveryOrPickUp>{
      at: this.formGroup.get("deliveryAt")!.getRawValue(),
      by: this.deliveryBy.getRawValue().map(x => x.name),
    };

    reservation.pickUp = <DeliveryOrPickUp>{
      at: this.formGroup.get("pickUpAt")!.getRawValue(),
      by: this.pickUpBy.getRawValue().map(x => x.name),
    };

    let obs: Observable<Reservation>;
    if (this.isCreateMode()) {
      obs = this.reservationService.create(reservation)
        .pipe(
          tap(r => reservation.id = r.id),
          switchMap(() => this.reservationService.update(reservation.id!, reservation))
        );
    } else {
      obs =this.reservationService.update(this.data.id, reservation);
    }

    obs
      .pipe(tap(_ => this.dialogRef.close()))
      .subscribe();
  }

  cancel(): void {
    this.dialogRef.close();
  }

  delete(): void {
    this.reservationService.delete(this.data.id)
      .pipe(tap(() => this.dialogRef.close()))
      .subscribe()
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
