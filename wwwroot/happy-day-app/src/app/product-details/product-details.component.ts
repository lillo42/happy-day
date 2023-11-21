import {Component, OnInit} from '@angular/core';
import {DatePipe} from "@angular/common";
import {HttpErrorResponse} from "@angular/common/http";
import {FormBuilder, FormGroup, Validators} from "@angular/forms";
import {ActivatedRoute, Router} from "@angular/router";

import {MatSnackBar} from "@angular/material/snack-bar";

import {Product, ProductsService} from "../products.service";
import {of, switchMap} from "rxjs";
import {ProblemDetails} from "../common";
import {Customer} from "../customers.service";

@Component({
  selector: 'app-product-details',
  templateUrl: './product-details.component.html',
  styleUrls: ['./product-details.component.scss']
})
export class ProductDetailsComponent implements OnInit {
  form: FormGroup;
  id: string | null = null;
  hasFound: boolean = true;
  isNew: boolean = true;

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
        if (id !== null && id !== 'new') {
          this.id = id;
          this.isNew = false;
          return this.productsService.getById(id);
        } else {
          const empty = {
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
        next: product => this.updateForm(product),
        error: error => {
          if (error instanceof HttpErrorResponse) {
            if (error.status === 404) {
              this.hasFound = false;
            } else {
              const problemDetails: ProblemDetails = JSON.parse(error.message);
              this.snackBar.open(`an unexpected error happen: ${problemDetails.message}`, 'OK', {duration: 10000});
            }
            return;
          }

          this.snackBar.open(`an unexpected error happen: ${error.toString()}`, 'OK', {duration: 10000});
        }
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

    if(this.isNew) {
      this.productsService
        .create(product)
        .subscribe({
          next: product => {
            this.updateForm(product);
            this.isNew = false;
          },
          error: error => this.handleError(error)
        });
    } else {
      this.productsService
        .update(this.id!, product)
        .subscribe({
          next: product => this.form.patchValue(product),
          error: error => this.handleError(error)
        });
    }
  }

  private handleError(error: HttpErrorResponse): void {
    if (error.status == 400) {
      this.form.markAllAsTouched();
      this.form.markAsDirty();
      return;
    }

    if (error.status == 0) {
      this.snackBar.open(`an unexpected error happen: ${error.message}`, 'OK', {duration: 10000});
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
      this.snackBar.open('product update conflict, please reload the page', 'OK', {duration: 10000});
    } else {
      this.snackBar.open(`an unexpected error happen: ${problemDetails.message}`, 'OK', {duration: 10000});
    }
  }

  private updateForm(product: Product) {
    this.form.patchValue({...product});

    this.form.get("createAt")!.setValue(this.datePipe.transform(product.createAt, 'dd/MM/yyyy HH:mm:ss'));
    this.form.get("updateAt")!.setValue(this.datePipe.transform(product.updateAt, 'dd/MM/yyyy HH:mm:ss'));
  }
}
