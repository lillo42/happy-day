import {AfterViewInit, Component, ElementRef, ViewChild} from '@angular/core';
import {CommonModule} from '@angular/common';
import {Router, RouterLink} from "@angular/router";

import {MatButtonModule} from "@angular/material/button";
import {MatDialog} from "@angular/material/dialog";
import {MatOptionModule} from "@angular/material/core";
import {MatFormFieldModule} from "@angular/material/form-field";
import {MatIconModule} from "@angular/material/icon";
import {MatInputModule} from "@angular/material/input";
import {MatPaginator, MatPaginatorModule} from "@angular/material/paginator";
import {MatSelect, MatSelectModule} from "@angular/material/select";
import {MatSnackBar} from "@angular/material/snack-bar";
import {MatTableDataSource, MatTableModule} from "@angular/material/table";
import {debounceTime} from "rxjs";

import {OrderDeleteComponent} from "../order-delete/order-delete.component";
import {OrdersService} from "../orders.service";


@Component({
  selector: 'app-order-list',
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatFormFieldModule, MatIconModule, MatInputModule, MatOptionModule, MatPaginatorModule, MatSelectModule, MatTableModule, RouterLink],
  templateUrl: './order-list.component.html',
  styleUrl: './order-list.component.scss'
})
export class OrderListComponent implements AfterViewInit {
  dataSourceLength = 0;
  displayedColumns: string[] = ['id', 'customerName', 'deliveryAt', 'pickUpAt', 'finalPrice', 'actions'];
  dataSource: MatTableDataSource<OrderElement>;

  @ViewChild(MatPaginator) paginator: MatPaginator | null = null;
  @ViewChild('selectField') field: MatSelect | null = null;
  @ViewChild('inputFilter') filter: ElementRef | null = null;

  constructor(private orderService: OrdersService,
              private router: Router,
              private snack: MatSnackBar,
              private dialog: MatDialog) {
    this.dataSource = new MatTableDataSource<OrderElement>([]);
  }

  ngAfterViewInit(): void {
    this.dataSource.paginator = this.paginator;
    this.load();
  }

  delete(id: string): void {
    this.dialog.open(OrderDeleteComponent, {data: {id: id}})
      .afterClosed()
      .subscribe(() => this.load());
  }

  edit(id: string): Promise<boolean> {
    return this.router.navigateByUrl(`/orders/${id}`);
  }

  load(): void {
    const field = this.field?.value || '';
    const value = this.filter?.nativeElement.value || '';

    let name = null;
    let address = null;
    let customerName = null;
    let customerPhone = null;

    if (field === 'name') {
      name = value;
    } else if (field === 'address') {
      address = value;
    } else if (field === 'customerName') {
      customerName = value;
    } else if (field === 'customerPhone') {
      customerPhone = value;
    }

    const page = this.paginator?.pageIndex || 0;
    const size = this.paginator?.pageSize || 50;

    this.orderService.get(name, address, customerName, customerPhone, page, size)
      .pipe(debounceTime(1000))
      .subscribe({
        next: page => {
          if (page.items === null) {
            this.dataSourceLength = 0;
            this.dataSource.data = [];
            return;
          }

          this.dataSourceLength = page.totalPages;
          this.dataSource.data = page.items.map(order => {
            return <OrderElement>{
              id: order.id,
              customerName: order.customer.name,
              deliveryAt: order.deliveryAt,
              pickUpAt: order.pickUpAt,
              finalPrice: order.finalPrice
            };
          });
        },
        error: err => this.snack.open(err.message, 'OK', {duration: 5000})
      })
  }
}

export interface OrderElement {
  id: string;
  customerName: string;
  deliveryAt: Date;
  pickUpAt: Date;
  finalPrice: number;
}
