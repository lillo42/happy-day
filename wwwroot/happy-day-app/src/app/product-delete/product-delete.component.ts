import { Component, Inject } from "@angular/core";
import { HttpErrorResponse } from "@angular/common/http";

import { MAT_DIALOG_DATA, MatDialogRef } from "@angular/material/dialog";

import { CustomerDeleteView } from "../customer-delete/customer-delete.component";
import { ProductsService } from "../products.service";

@Component({
  selector: 'app-product-delete',
  templateUrl: './product-delete.component.html',
  styleUrls: ['./product-delete.component.scss']
})
export class ProductDeleteComponent {
  error: string | null = null;

  constructor(public dialogRef: MatDialogRef<ProductDeleteView>,
              @Inject(MAT_DIALOG_DATA) public data: CustomerDeleteView,
              private productService: ProductsService) {
  }

  cancel(): void {
    this.dialogRef.close();
  }

  delete(): void {
    this.productService
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

export interface ProductDeleteView {
  id: string;
}
