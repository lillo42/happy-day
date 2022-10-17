import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef} from "@angular/material/dialog";
import { AbstractControl, FormArray, FormBuilder, FormControl, FormGroup, Validator, Validators } from "@angular/forms";
import { CustomerService } from "../http-clients/customer.service";
import { Observable, tap } from "rxjs";
import { Customer, Phone } from "../models/customer";
import {DatePipe} from "@angular/common";

@Component({
  selector: 'app-customer',
  templateUrl: './customer.component.html',
  styleUrls: ['./customer.component.scss'],
  providers: [DatePipe]
})
export class CustomerComponent implements OnInit {

  CREATE_TITLE = "Criação de um cliente";
  CHANGE_TITLE = "Atualização de um cliente";
  DELETE_TITLE = "Remover um cliente";

  formGroup: FormGroup;

  constructor(private dialogRef: MatDialogRef<CustomerComponent>,
              @Inject(MAT_DIALOG_DATA)private data: CustomerData,
              private builder: FormBuilder,
              private customerService: CustomerService,
              private datePipe: DatePipe) {
    this.formGroup = builder.group({
      id: [{value: null, disabled: true}],
      name: [null, [Validators.required]],
      comment: [null, null],
      phones: builder.array([]),
      createAt: [{value: null, disabled: true}],
      modifyAt: [{value: null, disabled: true}],
    })
  }

  ngOnInit(): void {
    if(this.data.behavior === CustomerBehavior.Create){
      this.addPhone();
      return;
    }

    this.customerService.get(this.data.id)
      .pipe(
        tap(customer => {
          for(let i = 0; i < customer.phones.length; i++) {
            this.addPhone(customer.phones[i]);
          }
          this.formGroup.patchValue({...customer});
          this.formGroup.get("createAt")!.setValue(this.datePipe.transform(customer.createAt, 'dd/MM/yyyy HH:mm:ss'));
          this.formGroup.get("modifyAt")!.setValue(this.datePipe.transform(customer.modifyAt, 'dd/MM/yyyy HH:mm:ss'));
        })
      )
      .subscribe();
  }

  get phones(): FormArray {
    return this.formGroup.controls["phones"] as FormArray;
  }

  isCreateMode(): boolean {
    return this.data.behavior === CustomerBehavior.Create;
  }

  isChangeMode(): boolean {
    return this.data.behavior === CustomerBehavior.Change;
  }

  isDeleteMode(): boolean {
    return this.data.behavior === CustomerBehavior.Delete;
  }

  deletePhone(index: number): void {
    this.phones.removeAt(index);
    if(this.phones.length == 0) {
      this.addPhone();
    }
  }

  addPhone(phone: Phone | null = null): void {
    this.phones.push(this.builder.group({
      number: [{value: phone?.number, disabled: this.isDeleteMode() }, [Validators.required, Validators.pattern('[- +()0-9]+')]]
    }));
  }

  cancel(): void {
    this.dialogRef.close()
  }

  delete(): void {
    this.customerService.delete(this.data.id)
      .pipe(tap(() => this.dialogRef.close()))
      .subscribe()
  }

  save(): void {
    this.formGroup.markAllAsTouched();
    this.formGroup.markAsDirty();

    if (!this.formGroup.valid) {
      return;
    }

    let customer = <Customer>{...this.formGroup.value};
    let observable: Observable<Customer> | null;
    if(this.data.behavior == CustomerBehavior.Create) {
      observable = this.customerService.create(customer);
    } else {
      observable = this.customerService.update(this.data.id, customer);
    }

    observable.pipe(tap(() => this.dialogRef.close()))
      .subscribe();
  }
}

export enum CustomerBehavior {
  Create,
  Change,
  Delete
}

export interface CustomerData {
  behavior: CustomerBehavior;
  id: string;
}
