import {AfterViewInit, Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {MatSort, Sort} from "@angular/material/sort";
import {MatPaginator} from "@angular/material/paginator";
import {CollectionViewer, DataSource} from "@angular/cdk/collections";
import {
  BehaviorSubject,
  catchError,
  debounceTime,
  finalize,
  Observable,
  of,
  Subject,
  tap
} from "rxjs";
import {CustomerSort} from "../models/customer";
import {CustomerService} from "../http-clients/customer.service";
import {MatDialog} from "@angular/material/dialog";
import {CustomerBehavior, CustomerComponent, CustomerData} from "../customer/customer.component";

@Component({
  selector: 'app-customer-lists',
  templateUrl: './customer-lists.component.html',
  styleUrls: ['./customer-lists.component.scss']
})
export class CustomerListsComponent implements OnInit, AfterViewInit {
  displayedColumns = ["name",  "comment", "phones", "actions"]
  dataSource: CustomerDataSource;

  private sortBy = CustomerSort.NameAsc;
  private textChanged = new Subject<string>();

  @ViewChild("filter") filter: ElementRef | null = null;
  @ViewChild(MatPaginator) paginator: MatPaginator | null = null;
  @ViewChild(MatSort) sort: MatSort | null = null;

  constructor(
    private dialog: MatDialog,
    customerService: CustomerService
  ) {
    this.dataSource = new CustomerDataSource(customerService);
  }

  ngOnInit(): void {
    this.textChanged
      .pipe(debounceTime(1000))
      .subscribe(() => this.reload());
  }

  ngAfterViewInit(): void {
    if(this.sort === null || this.paginator === null) {
      return;
    }

    this.paginator.page
      .pipe(tap(() => this.reload()))
      .subscribe();

    this.reload();
  }

  createOrUpdate(id: string) {
    const dialogRef = this.dialog.open(CustomerComponent, {
      width: '80%',
      data: <CustomerData>{
        behavior: id === "" ? CustomerBehavior.Create : CustomerBehavior.Change,
        id: id,
      },
    });

    dialogRef.afterClosed()
      .pipe(tap(() => this.reload()))
      .subscribe();
  }

  delete(id: string): void {
    const dialogRef = this.dialog.open(CustomerComponent, {
      width: '80%',
      data: <CustomerData>{
        behavior: CustomerBehavior.Delete,
        id: id,
      },
    });

    dialogRef.afterClosed()
      .pipe(tap(() => this.reload()))
      .subscribe();
  }

  applyFilter(key: string): void {
    this.textChanged.next(key);
  }

  sortChange(sort: Sort): void {
    if(this.paginator != null) {
      this.paginator.pageIndex = 0;
    }

    if (sort.active == "id") {
      if (sort.direction == "asc") {
        this.sortBy = CustomerSort.IdAsc;
      } else {
        this.sortBy = CustomerSort.IdDesc;
      }
    } else if (sort.active == "name") {
      if (sort.direction == "asc") {
        this.sortBy = CustomerSort.NameAsc;
      } else {
        this.sortBy = CustomerSort.NameDesc;
      }
    } else if (sort.active == "comment") {
      if (sort.direction == "asc") {
        this.sortBy = CustomerSort.CommentAsc;
      } else {
        this.sortBy = CustomerSort.CommentDesc;
      }
    }

    this.reload();
  }

  private reload(): void {
    if(this.paginator !== null && this.sort !== null && this.filter !== null) {
      this.dataSource.get(this.paginator.pageIndex + 1, this.paginator.pageSize,
        this.filter.nativeElement.value, this.sortBy);
    }
  }
}

interface CustomerViewModel {
  id: string;
  name: string;
  comment: string;
  phones: string;
}

class CustomerDataSource implements DataSource<CustomerViewModel> {
  private items = new BehaviorSubject<CustomerViewModel[]>([]);
  private totalElements = new BehaviorSubject<number>(0);
  private isLoading = new BehaviorSubject<boolean>(false);

  public totalElements$ = this.totalElements.asObservable();
  constructor(private customerService: CustomerService) { }

  connect(collectionViewer: CollectionViewer): Observable<CustomerViewModel[]> {
    return this.items.asObservable();
  }

  disconnect(collectionViewer: CollectionViewer): void {
    this.items.complete();
    this.totalElements.complete();
    this.isLoading.complete();
  }

  get(page: number, size: number, text: string, sort: CustomerSort | null): void {
    this.isLoading.next(true);
    this.customerService.getAll(page, size, text, sort)
      .pipe(
        tap(page => {
          this.totalElements.next(page.totalElements);
          this.items.next(page.items.map(customer => <CustomerViewModel>{
              id: customer.id,
              name: customer.name,
              comment: customer.comment,
              phones: customer.phones.map(x => x.number).join("\n")
            }));
        }),
        catchError(() => of([])),
        finalize(() => this.isLoading.next(false))
      ).subscribe();
  }
}
