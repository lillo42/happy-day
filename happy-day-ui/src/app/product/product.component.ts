import {Component, Inject, OnInit} from '@angular/core';
import {FormArray, FormBuilder, FormGroup, Validators} from "@angular/forms";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {ProductService} from "../http-clients/product.service";
import {DatePipe} from "@angular/common";
import {BehaviorSubject, debounceTime, Observable, switchMap, tap} from "rxjs";
import {InnerProduct, Product, ProductSort} from "../models/product";

@Component({
  selector: 'app-product',
  templateUrl: './product.component.html',
  styleUrls: ['./product.component.scss'],
  providers: [DatePipe]
})
export class ProductComponent implements OnInit {
  CREATE_TITLE = "Criação de um produto";
  CHANGE_TITLE = "Atualização de um produto";
  DELETE_TITLE = "Remover um produto";

  formGroup: FormGroup;
  filterProducts = new BehaviorSubject<Product[]>([]);
  filterProducts$: Observable<Product[]>;

  constructor(private dialogRef: MatDialogRef<ProductComponent>,
              @Inject(MAT_DIALOG_DATA)private data: ProductData,
              private builder: FormBuilder,
              private productService: ProductService,
              private datePipe: DatePipe) {
    this.filterProducts$ = this.filterProducts.asObservable();

    this.formGroup = builder.group({
      id: [{value: null, disabled: true}],
      name: [null, [Validators.required]],
      price: [null, [Validators.required, Validators.min(0)]],
      products: builder.array([]),
      createAt: [{value: null, disabled: true}],
      modifyAt: [{value: null, disabled: true}],
    });
  }

  ngOnInit(): void {
    if(this.data.behavior == ProductBehavior.Create){
      return;
    }

    this.productService.get(this.data.id)
      .pipe(
        tap(product => {
          for(let i = 0; i < product.products.length; i++) {
            this.addProduct(product.products[i]);
          }
          this.formGroup.patchValue({...product});
          this.formGroup.get("createAt")!.setValue(this.datePipe.transform(product.createAt, 'dd/MM/yyyy HH:mm:ss'));
          this.formGroup.get("modifyAt")!.setValue(this.datePipe.transform(product.modifyAt, 'dd/MM/yyyy HH:mm:ss'));
        })
      )
      .subscribe();
  }

  get products(): FormArray {
    return this.formGroup.controls["products"] as FormArray;
  }

  isCreateMode(): boolean {
    return this.data.behavior === ProductBehavior.Create;
  }

  isChangeMode(): boolean {
    return this.data.behavior === ProductBehavior.Change;
  }

  isDeleteMode(): boolean {
    return this.data.behavior === ProductBehavior.Delete;
  }

  deleteProduct(index: number): void {
    this.products.removeAt(index);
    if(this.products.length == 0) {
      this.addProduct();
    }
  }

  addProduct(product: InnerProduct | null = null): void {
    this.products.push(this.builder.group({
      id: [{value: product?.id, disabled: true}],
      name: [{value: null, disabled: this.isDeleteMode()}, [Validators.required]],
      quantity: [{value: product?.quantity, disabled: this.isDeleteMode()}, [Validators.required, Validators.min(1)]],
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

  cancel(): void {
    this.dialogRef.close();
  }

  delete(): void {
    this.productService.delete(this.data.id)
      .pipe(tap(() => this.dialogRef.close()))
      .subscribe();
  }

  save(): void {
    this.formGroup.markAllAsTouched();
    this.formGroup.markAsDirty();

    if (!this.formGroup.valid) {
      return;
    }

    let product = <Product>{...this.formGroup.value};
    product.products = this.products.getRawValue();
    let observable: Observable<Product> | null;
    if(this.data.behavior == ProductBehavior.Create) {
      observable = this.productService.create(product);
    } else {
      observable = this.productService.update(this.data.id, product);
    }

    observable.pipe(tap(() => this.dialogRef.close()))
      .subscribe();
  }
}

export enum ProductBehavior {
  Create,
  Change,
  Delete
}

export interface ProductData {
  behavior: ProductBehavior;
  id: string;
}
