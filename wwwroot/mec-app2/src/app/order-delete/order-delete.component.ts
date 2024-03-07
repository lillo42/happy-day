import {HttpErrorResponse} from "@angular/common/http";
import {Component, Inject} from '@angular/core';
import {CommonModule} from '@angular/common';

import {MatButtonModule} from "@angular/material/button";
import {
  MAT_DIALOG_DATA,
  MatDialogActions,
  MatDialogContent,
  MatDialogRef,
  MatDialogTitle
} from "@angular/material/dialog";

import {OrdersService} from "../orders.service";

@Component({
  selector: 'app-order-delete',
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatDialogActions, MatDialogContent, MatDialogTitle],
  templateUrl: './order-delete.component.html',
  styleUrl: './order-delete.component.scss'
})
export class OrderDeleteComponent {
  error: string | null = null;

  constructor(public dialogRef: MatDialogRef<OrderDeleteView>,
              @Inject(MAT_DIALOG_DATA) public data: OrderDeleteView,
              private ordersService: OrdersService) {
  }

  cancel(): void {
    this.dialogRef.close();
  }

  delete(): void {
    this.ordersService
      .delete(this.data.id)
      .subscribe({
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


export interface OrderDeleteView {
  id: string;
}
