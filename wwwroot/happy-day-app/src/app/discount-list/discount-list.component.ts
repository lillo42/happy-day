import { AfterViewInit, Component, ElementRef, ViewChild } from "@angular/core";
import { CommonModule } from "@angular/common";
import { Router, RouterLink } from "@angular/router";

import { MatButtonModule } from "@angular/material/button";
import { MatOptionModule } from "@angular/material/core";
import { MatDialog } from "@angular/material/dialog";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatIconModule } from "@angular/material/icon";
import { MatInputModule } from "@angular/material/input";
import { MatPaginator, MatPaginatorModule } from "@angular/material/paginator";
import { MatSelect, MatSelectModule } from "@angular/material/select";
import { MatSnackBar } from "@angular/material/snack-bar";
import { MatTableDataSource, MatTableModule } from "@angular/material/table";

import { debounceTime } from "rxjs";

import { DiscountDeleteComponent } from "../discount-delete/discount-delete.component";
import { DiscountsService } from "../discounts.service";
import { ProductElement } from "../product-list/product-list.component";

@Component({
  selector: 'app-discount-list',
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatFormFieldModule, MatIconModule, MatInputModule, MatOptionModule, MatPaginatorModule, MatSelectModule, MatTableModule, RouterLink],
  templateUrl: './discount-list.component.html',
  styleUrl: './discount-list.component.scss'
})
export class DiscountListComponent implements AfterViewInit {
  dataSourceLength = 0;
  displayedColumns: string[] = ['id', 'name', 'price', 'actions'];
  dataSource: MatTableDataSource<DiscountElement>;

  @ViewChild(MatPaginator) paginator: MatPaginator | null = null;
  @ViewChild('selectField') field: MatSelect | null = null;
  @ViewChild('inputFilter') filter: ElementRef | null = null;

  constructor(private discountsService: DiscountsService,
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
    this.dialog.open(DiscountDeleteComponent, {data: {id: id}})
      .afterClosed()
      .subscribe(() => this.load());
  }

  edit(id: string): Promise<boolean> {
    return this.router.navigateByUrl(`/discounts/${id}`);
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

    this.discountsService.get(name, page, size)
      .pipe(debounceTime(1000))
      .subscribe({
        next: page => {
          if (page.items === null) {
            this.dataSourceLength = 0;
            this.dataSource.data = [];
            return;
          }

          this.dataSource.data = page.items.map(discount => {
            return <DiscountElement>{
              id: discount.id,
              name: discount.name,
              price: discount.price,
            }
          });

          this.dataSourceLength = page.totalPages;
        },
        error: err => this.snack.open(err.message, 'OK')
      });
  }
}

export interface DiscountElement {
  id: string;
  name: string;
  price: number;
}
