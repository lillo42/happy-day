import {DataSource} from "@angular/cdk/collections";
import {AfterViewInit, Component, ElementRef, OnDestroy, ViewChild} from '@angular/core';
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
import {MatTableModule} from "@angular/material/table";
import {BehaviorSubject, debounceTime, Observable, Subscription} from "rxjs";

import {OrderDeleteComponent} from "../order-delete/order-delete.component";
import {OrdersService} from "../orders.service";


@Component({
  selector: 'app-order-list',
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatFormFieldModule, MatIconModule, MatInputModule, MatOptionModule, MatPaginatorModule, MatSelectModule, MatTableModule, RouterLink],
  templateUrl: './order-list.component.html',
  styleUrl: './order-list.component.scss'
})
export class OrderListComponent implements AfterViewInit, OnDestroy {
  displayedColumns: string[] = ['id', 'customerName', 'deliveryAt', 'pickUpAt', 'finalPrice', 'actions'];
  dataSource: OrderDataSource;

  @ViewChild(MatPaginator) paginator: MatPaginator | null = null;
  @ViewChild('selectField') field: MatSelect | null = null;
  @ViewChild('inputFilter') filter: ElementRef | null = null;

  private paginatorSubscription: Subscription | null = null;

  constructor(orderService: OrdersService,
              private router: Router,
              private snack: MatSnackBar,
              private dialog: MatDialog) {
    this.dataSource = new OrderDataSource(orderService);
  }

  ngAfterViewInit(): void {
    if (this.paginator === null) {
      return;
    }

    this.paginatorSubscription = this.paginator.page.subscribe(() => this.load());
    this.load();
  }

  ngOnDestroy(): void {
    this.paginatorSubscription?.unsubscribe();
    this.paginatorSubscription = null;
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

    this.dataSource.load(name, address, customerName, customerPhone, page, size, err => this.snack.open(err.message, 'OK'));
  }
}

interface OrderElement {
  id: string;
  customerName: string;
  deliveryAt: Date;
  pickUpAt: Date;
  finalPrice: number;
}

class OrderDataSource implements DataSource<OrderElement> {
  public totalElements$: Observable<number>;
  private items: BehaviorSubject<OrderElement[]>;
  private totalItems: BehaviorSubject<number>;
  private isLoading: BehaviorSubject<boolean>;

  constructor(private ordersService: OrdersService) {
    this.items = new BehaviorSubject<OrderElement[]>([]);
    this.totalItems = new BehaviorSubject<number>(0);
    this.isLoading = new BehaviorSubject<boolean>(false);
    this.totalElements$ = this.totalItems.asObservable();
  }

  connect(): Observable<OrderElement[]> {
    return this.items.asObservable();
  }

  disconnect(): void {
    this.items.complete();
    this.totalItems.complete();
    this.isLoading.complete();
  }

  load(name: string | null, address: string | null, customerName: string | null, customerPhone: string | null, page: number, size: number, error: (err: any) => void): void {
    this.isLoading.next(true);
    this.ordersService.get(name, address, customerName, customerPhone, page, size)
      .pipe(debounceTime(1000))
      .subscribe({
        next: page => {
          if (page.items === null) {
            this.items.next([]);
            this.totalItems.next(0);
            return;
          }

          this.items.next(page.items.map(order => {
            return <OrderElement>{
              id: order.id,
              customerName: order.customer.name,
              deliveryAt: order.deliveryAt,
              pickUpAt: order.pickUpAt,
              finalPrice: order.finalPrice
            };
          }));
          this.totalItems.next(page.totalItems);
        },
        error: err => error(err),
        complete: () => this.isLoading.next(false)
      })
  }
}
