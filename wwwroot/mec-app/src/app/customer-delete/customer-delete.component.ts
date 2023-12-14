import {Component, Inject} from '@angular/core';
import {HttpErrorResponse} from "@angular/common/http";

import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";

import {CustomersService} from "../customers.service";

@Component({
  selector: 'app-customer-delete',
  templateUrl: './customer-delete.component.html',
  styleUrls: ['./customer-delete.component.scss']
})
export class CustomerDeleteComponent {
  error: string | null = null;

  constructor(public dialogRef: MatDialogRef<CustomerDeleteComponent>,
              @Inject(MAT_DIALOG_DATA) public data: CustomerDeleteView,
              private customerService: CustomersService) {

  }

  cancel(): void {
    this.dialogRef.close();
  }

  delete(): void {
    this.customerService.delete(this.data.id).subscribe({
      next: () => this.dialogRef.close(),
      error: err => {
        if (err instanceof HttpErrorResponse && err.status > 0) {
          this.error = err.error.message;
        } else {
          this.error = err.message;
        }
      }
    });
  }
}

export interface CustomerDeleteView {
  id: string;
}
