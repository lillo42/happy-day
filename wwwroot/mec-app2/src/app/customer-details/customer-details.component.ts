import { DatePipe } from '@angular/common';
import { HttpErrorResponse } from '@angular/common/http';
import { Component, computed, OnInit, signal } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { $localize } from '@angular/localize/init';
import { ActivatedRoute, Router } from '@angular/router';

import { MatSnackBar } from '@angular/material/snack-bar';

import { Observable, of, switchMap, tap } from 'rxjs';

import { Customer, CustomersService } from '../customers.service';
import { ProblemDetails } from '../common';

@Component({
  selector: 'app-customer-details',
  templateUrl: './customer-details.component.html',
  styleUrls: ['./customer-details.component.scss'],
})
export class CustomerDetailsComponent implements OnInit {
  form: FormGroup;
  id: string | null = null;

  isNew = computed(() => this.id === null);
  hasFound = signal(false);
  isLoading = signal(true);

  constructor(
    private activatedRoute: ActivatedRoute,
    private datePipe: DatePipe,
    private builder: FormBuilder,
    private router: Router,
    private customersService: CustomersService,
    private snackBar: MatSnackBar,
  ) {
    this.form = this.builder.group({
      name: [null, [Validators.required, Validators.maxLength(255)]],
      comment: [null, null],
      phones: this.builder.array([]),
      createAt: [{ value: null, disabled: true }, null],
      updateAt: [{ value: null, disabled: true }, null],
    });
  }

  get phones(): FormArray {
    return this.form.get('phones') as FormArray;
  }

  ngOnInit(): void {
    this.activatedRoute.paramMap
      .pipe(
        switchMap((params) => {
          const id = params.get('id');
          if (id !== null && id !== 'new') {
            this.id = id;
            return this.customersService.getById(id);
          } else {
            const empty: Customer = {
              id: '',
              name: '',
              comment: '',
              phones: [],
              createAt: new Date(),
              updateAt: new Date(),
            };
            return of(empty);
          }
        }),
      )
      .pipe(tap(() => this.isLoading.set(false)))
      .subscribe({
        next: (customer) => {
          this.updateForm(customer);
          this.hasFound.set(true);
        },
        error: (error) => {
          if (error instanceof HttpErrorResponse) {
            if (error.status !== 404) {
              const problemDetails: ProblemDetails = JSON.parse(error.message);
              this.snackBar.open(
                $localize`an unexpected error happen: ${problemDetails.message}`,
                'OK',
                { duration: 10000 },
              );
            }
            return;
          }

          this.snackBar.open(
            $localize`an unexpected error happen: ${error.toString()}`,
            'OK',
            { duration: 10000 },
          );
        },
      });
  }

  addPhone(phone: string | null = null) {
    this.phones.push(
      this.builder.control(phone, [
        Validators.required,
        Validators.minLength(8),
        Validators.maxLength(11),
        Validators.pattern('[- +()0-9]+'),
      ]),
    );
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
      ...this.form.value,
    };

    const save$: Observable<Customer> = this.isNew()
      ? this.customersService.create(customer)
      : this.customersService.update(this.id!, customer);
    save$.subscribe({
      next: (_) => this.router.navigateByUrl('/customers'),
      error: (err) => this.handlerError(err),
    });
  }

  cancel(): Promise<boolean> {
    return this.router.navigateByUrl('/customers');
  }

  private updateForm(customer: Customer) {
    this.form.patchValue({ ...customer });

    this.phones.clear();
    customer.phones.forEach((phone) => this.addPhone(phone));

    this.form
      .get('createAt')!
      .setValue(
        this.datePipe.transform(customer.createAt, 'dd/MM/yyyy HH:mm:ss'),
      );
    this.form
      .get('updateAt')!
      .setValue(
        this.datePipe.transform(customer.updateAt, 'dd/MM/yyyy HH:mm:ss'),
      );
  }

  private handlerError(error: HttpErrorResponse) {
    if (error.status === 400) {
      this.form.markAllAsTouched();
      this.form.markAsDirty();
      return;
    }

    if (error.status == 0) {
      this.snackBar.open(`an unexpected error happen: ${error.message}`, 'OK', {
        duration: 10000,
      });
      return;
    }

    const problemDetails: ProblemDetails = error.error;
    if (problemDetails.type === 'customer-name-is-empty') {
      this.form.get('name')!.setErrors({ required: true });
    } else if (problemDetails.type === 'customer-name-is-too-large') {
      this.form.get('name')!.setErrors({ maxlength: true });
    } else if (problemDetails.type === 'customer-pix-is-too-large') {
      this.form.get('pix')!.setErrors({ maxlength: true });
    } else if (problemDetails.type === 'customer-phone-number-is-invalid') {
      this.phones.controls.forEach((control) =>
        control.setErrors({ pattern: true }),
      );
    } else if (problemDetails.type === 'customer-conflict') {
      this.snackBar.open(
        $localize`customer update conflict, please reload the page`,
        'OK',
        { duration: 10000 },
      );
    } else {
      this.snackBar.open(
        $localize`an unexpected error happen: ${problemDetails.message}`,
        'OK',
        { duration: 10000 },
      );
    }
  }
}
