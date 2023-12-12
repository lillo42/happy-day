import { DatePipe } from "@angular/common";
import { HttpErrorResponse } from "@angular/common/http";
import { Component, computed, OnInit, signal } from "@angular/core";
import { FormBuilder, FormGroup, Validators } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";

import { MatSnackBar } from "@angular/material/snack-bar";

import { of, switchMap } from "rxjs";

import { Product, ProductsService } from "../products.service";
import { ProblemDetails } from "../common";

@Component({
  selector: 'app-product-details',
  templateUrl: './product-details.component.html',
  styleUrls: ['./product-details.component.scss']
})
export class ProductDetailsComponent implements OnInit {
  form: FormGroup;
  id = signal<string | null>(null);
  hasFound = signal(false);
  isNew = computed(() => this.id() === null);
  isLoading = signal(false);

  constructor(private activatedRoute: ActivatedRoute,
              private datePipe: DatePipe,
              private builder: FormBuilder,
              private router: Router,
              private snackBar: MatSnackBar,
              private productsService: ProductsService) {
    this.form = this.builder.group({
      name: [null, [Validators.required, Validators.maxLength(255)]],
      price: [null, [Validators.required, Validators.min(0)]],
      createAt: [{value: null, disabled: true}, null],
      updateAt: [{value: null, disabled: true}, null],
    });
  }

  ngOnInit(): void {
    this.activatedRoute.paramMap
      .pipe(switchMap(params => {
        const id = params.get('id');
        this.isLoading.set(true);
        if (id !== null && id !== 'new') {
          this.id.set(id)
          return this.productsService.getById(id);
        } else {
          const empty: Product = {
            id: '',
            name: '',
            price: 0,
            createAt: new Date(),
            updateAt: new Date(),
          };
          return of(empty);
        }
      }))
      .subscribe({
        next: product => {
          this.updateForm(product);
          this.hasFound.set(true);
        },
        error: error => {
          if (error instanceof HttpErrorResponse) {
            if (error.status === 404) {
              this.hasFound.set(false);
            } else {
              const problemDetails: ProblemDetails = JSON.parse(error.message);
              this.snackBar.open($localize `an unexpected error happen: ${problemDetails.message}`, 'OK', {duration: 10000});
            }
            return;
          }

          this.snackBar.open($localize `an unexpected error happen: ${error.toString()}`, 'OK', {duration: 10000});
        },
        complete: () => this.isLoading.set(false)
      });
  }

  cancel(): Promise<boolean> {
    return this.router.navigateByUrl('/products');
  }

  save(): void {
    this.form.markAllAsTouched();
    this.form.markAsDirty();
    if (!this.form.valid) {
      return;
    }

    const product = <Product>{
      ...this.form.value
    };

    const save$ = this.isNew() ? this.productsService.create(product) : this.productsService.update(this.id()!, product);
    save$.subscribe({
      next: _ => this.router.navigateByUrl('/products'),
      error: error => this.handleError(error)
    });
  }

  private handleError(error: HttpErrorResponse): void {
    if (error.status == 400) {
      this.form.markAllAsTouched();
      this.form.markAsDirty();
      return;
    }

    if (error.status == 0) {
      this.snackBar.open($localize `an unexpected error happen: ${error.message}`, 'OK', {duration: 10000});
      return;
    }

    const problemDetails: ProblemDetails = error.error;
    if (problemDetails.type === 'product-name-is-empty') {
      this.form.get('name')?.setErrors({required: true});
    } else if (problemDetails.type === 'product-name-is-too-large') {
      this.form.get('name')?.setErrors({maxlength: true});
    } else if (problemDetails.type === 'product-price-is-invalid') {
      this.form.get('price')?.setErrors({required: true, min: true});
    } else if (problemDetails.type === 'product-conflict') {
      this.snackBar.open($localize `product update conflict, please reload the page`, 'OK', {duration: 10000});
    } else {
      this.snackBar.open($localize `an unexpected error happen: ${problemDetails.message}`, 'OK', {duration: 10000});
    }
  }

  private updateForm(product: Product) {
    this.form.patchValue({...product});

    this.form.get("createAt")!.setValue(this.datePipe.transform(product.createAt, 'dd/MM/yyyy HH:mm:ss'));
    this.form.get("updateAt")!.setValue(this.datePipe.transform(product.updateAt, 'dd/MM/yyyy HH:mm:ss'));
  }
}
