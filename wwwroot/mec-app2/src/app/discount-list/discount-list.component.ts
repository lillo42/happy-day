import {DataSource} from "@angular/cdk/collections";
import {AfterViewInit, Component, ElementRef, OnDestroy, ViewChild} from "@angular/core";
import {CommonModule} from "@angular/common";
import {Router, RouterLink} from "@angular/router";

import {MatButtonModule} from "@angular/material/button";
import {MatOptionModule} from "@angular/material/core";
import {MatDialog} from "@angular/material/dialog";
import {MatFormFieldModule} from "@angular/material/form-field";
import {MatIconModule} from "@angular/material/icon";
import {MatInputModule} from "@angular/material/input";
import {MatPaginator, MatPaginatorModule} from "@angular/material/paginator";
import {MatSelect, MatSelectModule} from "@angular/material/select";
import {MatSnackBar} from "@angular/material/snack-bar";
import {MatTableModule} from "@angular/material/table";

import {BehaviorSubject, debounceTime, Observable, Subscription} from "rxjs";

import {DiscountDeleteComponent} from "../discount-delete/discount-delete.component";
import {DiscountsService} from "../discounts.service";

@Component({
  selector: 'app-discount-list',
  standalone: true,
  imports: [CommonModule, MatButtonModule, MatFormFieldModule, MatIconModule, MatInputModule, MatOptionModule, MatPaginatorModule, MatSelectModule, MatTableModule, RouterLink],
  templateUrl: './discount-list.component.html',
  styleUrl: './discount-list.component.scss'
})
export class DiscountListComponent implements AfterViewInit, OnDestroy {
  displayedColumns: string[] = ['id', 'name', 'price', 'actions'];
  dataSource: DiscountDataSource;

  @ViewChild(MatPaginator) paginator: MatPaginator | null = null;
  @ViewChild('selectField') field: MatSelect | null = null;
  @ViewChild('inputFilter') filter: ElementRef | null = null;

  private paginatorSubscription: Subscription | null = null;

  constructor(discountsService: DiscountsService,
              private router: Router,
              private snack: MatSnackBar,
              private dialog: MatDialog) {
    this.dataSource = new DiscountDataSource(discountsService);
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

    this.dataSource.load(name, page, size, err => this.snack.open(err.message, 'OK'));
  }
}

interface DiscountElement {
  id: string;
  name: string;
  price: number;
}

class DiscountDataSource implements DataSource<DiscountElement> {
  public totalElements$: Observable<number>;
  private items: BehaviorSubject<DiscountElement[]>;
  private totalItems: BehaviorSubject<number>;
  private isLoading: BehaviorSubject<boolean>;

  constructor(private discountsService: DiscountsService) {
    this.items = new BehaviorSubject<DiscountElement[]>([]);
    this.totalItems = new BehaviorSubject<number>(0);
    this.isLoading = new BehaviorSubject<boolean>(false);
    this.totalElements$ = this.totalItems.asObservable();
  }

  connect(): Observable<DiscountElement[]> {
    return this.items.asObservable();
  }

  disconnect(): void {
    this.items.complete();
    this.totalItems.complete();
    this.isLoading.complete();
  }

  load(name: string | null, page: number, size: number, error: (err: any) => void): void {
    this.isLoading.next(true);
    this.discountsService.get(name, page, size)
      .pipe(debounceTime(1000))
      .subscribe({
        next: page => {
          if (page.items === null) {
            this.items.next([]);
            this.totalItems.next(0);
            return;
          }

          this.items.next(page.items.map(discount => {
            return <DiscountElement>{
              id: discount.id,
              name: discount.name,
              price: discount.price,
            }
          }));
          this.totalItems.next(page.totalPages);
        },
        error: err => {
          error(err);
          this.items.next([]);
          this.totalItems.next(0);
        },
        complete: () => this.isLoading.next(false)
      });
  }
}
