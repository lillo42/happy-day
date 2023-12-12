import { CommonModule, DatePipe } from '@angular/common';
import { HttpErrorResponse } from "@angular/common/http";
import { AfterViewInit, Component, computed, ElementRef, OnInit, signal, ViewChild } from '@angular/core';
import { AbstractControl, FormArray, FormBuilder, FormGroup, ReactiveFormsModule, Validators } from "@angular/forms";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { ActivatedRoute, Router } from "@angular/router";

import { MatAutocompleteModule } from "@angular/material/autocomplete";
import { MatButtonModule } from "@angular/material/button";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatIconModule } from "@angular/material/icon";
import { MatInputModule } from "@angular/material/input";
import { MatSnackBar } from "@angular/material/snack-bar";
import { MatTooltipModule } from "@angular/material/tooltip";

import { NgxMaskDirective } from "ngx-mask";
import { of, switchMap } from "rxjs";

import { ProblemDetails } from "../common";
import { Discount, DiscountProduct, DiscountsService } from "../discounts.service";
import { ProductsService } from "../products.service";

@Component({
  selector: 'app-discount-details',
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatFormFieldModule, MatIconModule, MatInputModule, MatTooltipModule, NgxMaskDirective, ReactiveFormsModule, MatAutocompleteModule, MatProgressSpinnerModule],
  templateUrl: './discount-details.component.html',
  styleUrl: './discount-details.component.scss'
})
export class DiscountDetailsComponent implements OnInit, AfterViewInit {
  form: FormGroup;
  id = signal<string | null>(null);
  hasFound = signal(false);
  isNew = computed(() => this.id() === null);
  isLoading = signal(false);
  isProductsMissing = signal(false);

  filteredProducts = signal<DiscountProduct[]>([]);

  @ViewChild("productInput") productInput: ElementRef | null = null;

  constructor(private activatedRoute: ActivatedRoute,
              private datePipe: DatePipe,
              private builder: FormBuilder,
              private router: Router,
              private discountsService: DiscountsService,
              private productsService: ProductsService,
              private snackBar: MatSnackBar) {
    this.form = this.builder.group({
      name: [null, [Validators.required, Validators.maxLength(255)]],
      price: [null, [Validators.required, Validators.min(1)]],
      products: this.builder.array([]),
      createAt: [{value: null, disabled: true}, null],
      updateAt: [{value: null, disabled: true}, null],
    });
  }

  ngOnInit(): void {
    this.activatedRoute.paramMap
      .pipe(switchMap(params => {
        const id = params.get('id');
        if (id !== null && id !== 'new') {
          this.id.set(id);
          return this.discountsService.getById(id);
        } else {
          const empty: Discount = {
            id: '',
            name: '',
            price: 0,
            products: [],
            createAt: new Date(),
            updateAt: new Date(),
          };
          return of(empty);
        }
      }))
      .subscribe({
        next: discount => {
          this.updateForm(discount);
          this.hasFound.set(true);
        },
        error: error => {
          if (error instanceof HttpErrorResponse) {
            if (error.status == 404) {
              this.hasFound.set(true);
            } else {
              const problemDetails: ProblemDetails = JSON.parse(error.message);
              this.snackBar.open($localize`an unexpected error happen: ${problemDetails.message}`, 'OK', {duration: 10000});
            }
            return;
          }

          this.snackBar.open($localize`an unexpected error happen: ${error.toString()}`, 'OK', {duration: 10000});
        }
      });
  }

  ngAfterViewInit(): void {
    this.loadProducts();
  }

  get products(): FormArray {
    return this.form.get('products') as FormArray;
  }

  deleteProduct(index: number): void {
    this.products.removeAt(index);
  }

  addProduct(product: DiscountProduct): void {
    this.isProductsMissing.set(false);
    this.products.push(this.builder.group({
      id: [{value: product?.id, disabled: true}, [Validators.required]],
      name: [{value: product?.name, disabled: true}, null],
      quantity: [product?.quantity, [Validators.required, Validators.min(1)]],
    }));

    this.productInput!.nativeElement.value = '';
  }

  displayProduct(product: DiscountProduct | null): string {
    return product?.name || '';
  }

  asFormGroup(control: AbstractControl): FormGroup {
    return control as FormGroup;
  }

  cancel(): Promise<boolean> {
    return this.router.navigateByUrl('/discounts');
  }

  save(): void {
    this.form.markAllAsTouched();
    this.form.markAsDirty();
    if (this.form.invalid) {
      return;
    }

    const discount = <Discount>{
      ...this.form.getRawValue()
    };


    const save$ = this.isNew() ? this.discountsService.create(discount) : this.discountsService.update(this.id()!, discount);
    save$.subscribe({
      next: _ => this.router.navigateByUrl('/discounts'),
      error: error => this.handleError(error)
    });
  }

  loadProducts(): void {
    if (this.productInput === null) {
      return;
    }

    const name = this.productInput.nativeElement.value;
    this.productsService.get(name, 0, 100)
      .subscribe({
        next: page => {
          if (page.items === null) {
            this.filteredProducts.set([]);
            return;
          }

          this.filteredProducts.set(page.items.map(prod => {
            return <DiscountProduct>{
              id: prod.id,
              name: prod.name,
              quantity: 0,
            }
          }));
        },
        error: err => this.snackBar.open(err.message, 'OK')
      });
  }

  private handleError(error: HttpErrorResponse) {
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
    if (problemDetails.type === 'discount-name-is-empty') {
      this.form.get('name')?.setErrors({required: true});
    } else if (problemDetails.type === 'discount-name-is-too-large') {
      this.form.get('name')?.setErrors({maxlength: true});
    } else if (problemDetails.type === 'discount-price-is-invalid') {
      this.form.get('price')?.setErrors({required: true, min: true});
    } else if (problemDetails.type === 'discount-products-is-missing') {
      this.isProductsMissing.set(true);
    } else if (problemDetails.type === 'discount-conflict') {
      this.snackBar.open($localize`discount update conflict, please reload the page`, 'OK', {duration: 10000});
    } else {
      this.snackBar.open($localize`an unexpected error happen: ${problemDetails.message}`, 'OK', {duration: 10000});
    }
  }

  private updateForm(discount: Discount): void {
    this.form.patchValue({...discount});

    this.products.clear();
    discount.products.forEach(product => this.addProduct(product));
    this.form.get("createAt")!.setValue(this.datePipe.transform(discount.createAt, 'dd/MM/yyyy HH:mm:ss'));
    this.form.get("updateAt")!.setValue(this.datePipe.transform(discount.updateAt, 'dd/MM/yyyy HH:mm:ss'));
  }
}
