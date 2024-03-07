import {DataSource} from "@angular/cdk/collections";
import {AfterViewInit, Component, ElementRef, OnDestroy, ViewChild} from '@angular/core';
import {Router} from "@angular/router";

import {MatDialog} from "@angular/material/dialog";
import {MatPaginator} from "@angular/material/paginator";
import {MatSelect} from "@angular/material/select";
import {MatSnackBar} from "@angular/material/snack-bar";
import {NgxMaskService} from "ngx-mask";

import {BehaviorSubject, debounceTime, Observable, Subscription} from "rxjs";

import {CustomersService} from "../customers.service";
import {CustomerDeleteComponent} from "../customer-delete/customer-delete.component";

@Component({
  selector: 'app-customers-list',
  templateUrl: './customers-list.component.html',
  styleUrls: ['./customers-list.component.scss']
})
export class CustomersListComponent implements AfterViewInit, OnDestroy {
  displayedColumns: string[] = ['id', 'name', 'comment', 'phones', 'actions'];
  dataSource: CustomerDataSource;

  @ViewChild(MatPaginator) paginator: MatPaginator | null = null;
  @ViewChild('selectField') field: MatSelect | null = null;
  @ViewChild('inputFilter') filter: ElementRef | null = null;

  private paginatorSubscription: Subscription | null = null;

  constructor(private router: Router,
              private dialog: MatDialog,
              private snack: MatSnackBar,
              customersService: CustomersService,
              ngMask: NgxMaskService) {
    this.dataSource = new CustomerDataSource(customersService, ngMask);
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
    this.dialog.open(CustomerDeleteComponent, {data: {id: id}})
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
    let phone = null;
    let comment = null;

    if (field === 'name') {
      name = value;
    } else if (field === 'phone') {
      phone = value;
    } else if (field === 'comment') {
      comment = value;
    }

    const page = this.paginator?.pageIndex || 0;
    const size = this.paginator?.pageSize || 50;
    this.dataSource.load(name, phone, comment, page, size,
      err => this.snack.open(err.error.message, 'OK', {duration: 5000}));
  }
}

interface CustomerElement {
  id: string;
  name: string;
  comment: string;
  phones: string;
}

class CustomerDataSource implements DataSource<CustomerElement> {
  public totalElements$: Observable<number>;
  private items: BehaviorSubject<CustomerElement[]>;
  private totalItems: BehaviorSubject<number>;
  private isLoading: BehaviorSubject<boolean>;

  constructor(private customersService: CustomersService,
              private ngMask: NgxMaskService) {
    this.items = new BehaviorSubject<CustomerElement[]>([]);
    this.isLoading = new BehaviorSubject<boolean>(false);
    this.totalItems = new BehaviorSubject<number>(0);
    this.totalElements$ = this.totalItems.asObservable();
  }

  connect(): Observable<CustomerElement[]> {
    return this.items.asObservable();
  }

  disconnect(): void {
    this.items.complete();
    this.isLoading.complete();
    this.totalItems.complete();
  }

  load(name: string, phone: string, comment: string, page: number, size: number, onError: (err: any) => void): void {
    this.isLoading.next(true);
    this.customersService.get(name, phone, comment, page, size)
      .pipe(debounceTime(1000))
      .subscribe({
        next: page => {
          this.totalItems.next(page.totalItems);
          this.items.next(page.items?.map(customer => <CustomerElement>{
            id: customer.id,
            name: customer.name,
            comment: customer.comment,
            phones: customer.phones?.map(phone => this.ngMask.applyMask(phone, '(00) 00000-0000||(00) 0000-0000')).join(', ') || ''
          }) || []);
        },
        error: err => {
          onError(err);
          this.items.next([]);
        },
        complete: () => this.isLoading.next(false)
      });
  }
}
