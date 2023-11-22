import { CommonModule } from "@angular/common";
import { HttpErrorResponse } from "@angular/common/http";
import { Component, Inject } from '@angular/core';

import { MatButtonModule } from "@angular/material/button";
import {
  MAT_DIALOG_DATA,
  MatDialogActions,
  MatDialogContent,
  MatDialogRef,
  MatDialogTitle
} from "@angular/material/dialog";

import { DiscountsService } from "../discounts.service";

@Component({
  selector: 'app-discount-delete',
  standalone: true,
    imports: [CommonModule, MatButtonModule, MatDialogActions, MatDialogContent, MatDialogTitle],
  templateUrl: './discount-delete.component.html',
  styleUrl: './discount-delete.component.scss'
})
export class DiscountDeleteComponent {
  error: string | null = null;

  constructor(public dialogRef: MatDialogRef<DiscountDeleteView>,
              @Inject(MAT_DIALOG_DATA) public data: DiscountDeleteView,
              private discountsService: DiscountsService) {
  }

  cancel(): void {
    this.dialogRef.close();
  }

  delete(): void {
    this.discountsService
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

export interface DiscountDeleteView {
  id: string;
}
