import { AfterViewInit, Component, ElementRef, ViewChild } from '@angular/core';
import { Router } from "@angular/router";

import { MatDialog } from "@angular/material/dialog";
import { MatPaginator } from "@angular/material/paginator";
import { MatSelect } from "@angular/material/select";
import { MatSnackBar } from "@angular/material/snack-bar";
import { MatTableDataSource } from "@angular/material/table";

import { debounceTime } from "rxjs";

import { ProductsService } from "../products.service";
import { ProductDeleteComponent } from "../product-delete/product-delete.component";

@Component({
  selector: 'app-product-list',
  templateUrl: './product-list.component.html',
  styleUrls: ['./product-list.component.scss']
})
export class ProductListComponent implements AfterViewInit {
  dataSourceLength = 0;
  displayedColumns: string[] = ['id', 'name', 'price', 'actions'];
  dataSource: MatTableDataSource<ProductElement>;

  @ViewChild(MatPaginator) paginator: MatPaginator | null = null;
  @ViewChild('selectField') field: MatSelect | null = null;
  @ViewChild('inputFilter') filter: ElementRef | null = null;

  constructor(private productsService: ProductsService,
              private router: Router,
              private snack: MatSnackBar,
              private dialog: MatDialog) {
    this.dataSource = new MatTableDataSource<ProductElement>([]);
  }

  ngAfterViewInit(): void {
    this.dataSource.paginator = this.paginator;
    this.load();
  }

  delete(id: string): void {
    this.dialog.open(ProductDeleteComponent, {data: {id: id}})
      .afterClosed()
      .subscribe(() => this.load());
  }

  edit(id: string): Promise<boolean> {
    return this.router.navigateByUrl(`/customers/${id}`);
  }

  load(): void {
    const field = this.field?.value || '';
    const value = this.filter?.nativeElement.value || '';

    let name = null;

    if (field === 'name') {
      name = value;
    }

    const page = this.paginator?.pageIndex || 0;
    const size = this.paginator?.pageSize || 50;

    this.productsService.get(name, page, size)
      .pipe(debounceTime(1000))
      .subscribe({
        next: page => {
          if (page.items === null) {
            this.dataSource.data = [];
            return;
          }

          this.dataSource.data = page.items.map(product => {
            return <ProductElement>{
              id: product.id,
              name: product.name,
              price: product.price,
            }
          });

          this.dataSourceLength = page.totalPages;
        },
        error: err => this.snack.open(err.message, 'OK')
      });
  }
}

export interface ProductElement {
  id: string;
  name: string;
  price: number;
}
