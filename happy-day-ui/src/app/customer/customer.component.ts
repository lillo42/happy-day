import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogModule, MatDialogRef} from "@angular/material/dialog";
import {Form, FormArray, FormBuilder, FormControl, FormGroup, Validator, Validators} from "@angular/forms";
import {CustomerService} from "../http-clients/customer.service";
import {Observable, tap} from "rxjs";
import {createUrlTreeFromSnapshot} from "@angular/router";
import {Customer, Phone} from "../models/customer";

@Component({
  selector: 'app-customer',
  templateUrl: './customer.component.html',
  styleUrls: ['./customer.component.scss']
})
export class CustomerComponent implements OnInit {

  CREATE_TITLE = "Criação de um cliente";
  CHANGE_TITLE = "Atualização de um cliente";
  DELETE_TITLE = "Remover um cliente";

  formGroup: FormGroup;

  constructor(private dialogRef: MatDialogRef<CustomerComponent>,
              @Inject(MAT_DIALOG_DATA)private data: CustomerData,
              private builder: FormBuilder,
              private customerService: CustomerService) {
    this.formGroup = builder.group({
      id: [{value: null, disabled: true}],
      name: [null, [Validators.required]],
      comment: [null, null],
      phones: builder.array([])
    })
  }

  ngOnInit(): void {
    if(this.data.behavior === CustomerBehavior.Create){
      return;
    }

    this.customerService.get(this.data.id)
      .pipe(
        tap(customer => {
          for(let i = 0; i < customer.phones.length; i++) {
            this.phones.push(this.createPhone(customer.phones[i]));
          }
          this.formGroup.patchValue({...customer});
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
    this.phones.removeAt(index)
  }

  addPhone(): void {
    this.phones.push(this.createPhone());
  }

  createPhone(phone: Phone | null = null): FormGroup {
    return this.builder.group({
      number: [{value: phone?.number, disabled: this.isDeleteMode() }, [Validators.required, Validators.pattern('[- +()0-9]+')]]
    });
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
    let customer = <Customer>{...this.formGroup.value};
    let observable: Observable<Customer> | null = null;
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
