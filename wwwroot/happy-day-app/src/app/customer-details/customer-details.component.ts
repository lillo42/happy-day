import { DatePipe } from "@angular/common";
import { HttpErrorResponse } from "@angular/common/http";
import { Component, OnInit } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, Validators } from "@angular/forms";
import { ActivatedRoute, Router } from "@angular/router";

import { MatSnackBar } from "@angular/material/snack-bar";

import { of, switchMap } from "rxjs";

import { Customer, CustomersService } from "../customers.service";
import { ProblemDetails } from "../common";

@Component({
  selector: 'app-customer-details',
  templateUrl: './customer-details.component.html',
  styleUrls: ['./customer-details.component.scss']
})
export class CustomerDetailsComponent implements OnInit {
  form: FormGroup;
  id: string | null = null;
  isNew: boolean = true;
  hasFound: boolean = true;

  constructor(private activatedRoute: ActivatedRoute,
              private datePipe: DatePipe,
              private builder: FormBuilder,
              private router: Router,
              private customersService: CustomersService,
              private snackBar: MatSnackBar) {
    this.form = this.builder.group({
      name: [null, [Validators.required, Validators.maxLength(255)]],
      comment: [null, null],
      phones: this.builder.array([]),
      pix: [null, [Validators.maxLength(255)]],
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
          return this.customersService.getById(id);
        } else {
          const empty: Customer = {
            id: '',
            name: '',
            comment: '',
            phones: [],
            pix: '',
            createAt: new Date(),
            updateAt: new Date(),
          };
          return of(empty);
        }
      }))
      .subscribe({
        next: customer => this.updateForm(customer),
        error: error => {
          if (error instanceof HttpErrorResponse) {
            if (error.status == 404) {
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

  get phones(): FormArray {
    return this.form.get('phones') as FormArray;
  }

  addPhone(phone: string | null = null) {
    this.phones.push(this.builder.control(phone, [
      Validators.required,
      Validators.minLength(8),
      Validators.maxLength(11),
      Validators.pattern('[- +()0-9]+')
    ]));
  }

  deletePhone(index: number) {
    this.phones.removeAt(index);
  }

  save(): void {
    this.form.markAllAsTouched();
    this.form.markAsDirty();
    if (this.form.invalid) {
      return;
    }

    const customer = <Customer>{
      ...this.form.value
    };

    if (this.isNew) {
      this.customersService
        .create(customer)
        .subscribe({
          next: customer => {
            this.updateForm(customer);
            this.isNew = false;
            this.id = customer.id;
          },
          error: error => this.handlerError(error)
        });
    } else {
      this.customersService
        .update(this.id!, customer)
        .subscribe({
          next: customer => this.updateForm(customer),
          error: error => this.handlerError(error)
        });
    }
  }

  cancel(): Promise<boolean> {
    return this.router.navigateByUrl('/customers');
  }

  private updateForm(customer: Customer) {
    this.form.patchValue({...customer});

    this.phones.clear();
    customer.phones.forEach(phone => this.addPhone(phone));

    this.form.get("createAt")!.setValue(this.datePipe.transform(customer.createAt, 'dd/MM/yyyy HH:mm:ss'));
    this.form.get("updateAt")!.setValue(this.datePipe.transform(customer.updateAt, 'dd/MM/yyyy HH:mm:ss'));
  }

  private handlerError(error: HttpErrorResponse) {
    if (error.status === 400) {
      this.form.markAllAsTouched();
      this.form.markAsDirty();
      return;
    }

    if (error.status == 0) {
      this.snackBar.open(`an unexpected error happen: ${error.message}`, 'OK', {duration: 10000});
      return;
    }

    const problemDetails: ProblemDetails = error.error;
    if (problemDetails.type === 'customer-name-is-empty') {
      this.form.get('name')!.setErrors({required: true});
    } else if (problemDetails.type === 'customer-name-is-too-large') {
      this.form.get('name')!.setErrors({maxlength: true});
    } else if (problemDetails.type === 'customer-pix-is-too-large') {
      this.form.get('pix')!.setErrors({maxlength: true});
    } else if (problemDetails.type === 'customer-phone-number-is-invalid') {
      this.phones.controls.forEach(control => control.setErrors({pattern: true}));
    } else if (problemDetails.type === 'customer-conflict') {
      this.snackBar.open('customer update conflict, please reload the page', 'OK', {duration: 10000});
    } else {
      this.snackBar.open(`an unexpected error happen: ${problemDetails.message}`, 'OK', {duration: 10000});
    }
  }
}
