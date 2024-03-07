import {DataSource} from "@angular/cdk/collections";
import {AfterViewInit, Component, ElementRef, OnDestroy, ViewChild} from '@angular/core';
import {Router} from "@angular/router";

import {MatDialog} from "@angular/material/dialog";
import {MatPaginator} from "@angular/material/paginator";
import {MatSelect} from "@angular/material/select";
import {MatSnackBar} from "@angular/material/snack-bar";

import {BehaviorSubject, debounceTime, Observable, Subscription} from "rxjs";

import {ProductsService} from "../products.service";
import {ProductDeleteComponent} from "../product-delete/product-delete.component";

@Component({
  selector: 'app-product-list',
  templateUrl: './product-list.component.html',
  styleUrls: ['./product-list.component.scss']
})
export class ProductListComponent implements AfterViewInit, OnDestroy {
  displayedColumns: string[] = ['id', 'name', 'price', 'actions'];
  dataSource: ProductDataSource;

  @ViewChild(MatPaginator) paginator: MatPaginator | null = null;
  @ViewChild('selectField') field: MatSelect | null = null;
  @ViewChild('inputFilter') filter: ElementRef | null = null;

  private paginatorSubscription: Subscription | null = null;

  constructor(private router: Router,
              private snack: MatSnackBar,
              private dialog: MatDialog,
              productsService: ProductsService) {
    this.dataSource = new ProductDataSource(productsService);
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
    this.dialog.open(ProductDeleteComponent, {data: {id: id}})
      .afterClosed()
      .subscribe(() => this.load());
  }

  edit(id: string): Promise<boolean> {
    return this.router.navigateByUrl(`/products/${id}`);
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

    this.dataSource.load(name, page, size, (err: any) => {
      this.snack.open(err.message, 'OK', {duration: 5000});
    });
  }
}

interface ProductElement {
  id: string;
  name: string;
  price: number;
}

class ProductDataSource implements DataSource<ProductElement> {
  public totalElements$: Observable<number>;
  private items: BehaviorSubject<ProductElement[]>;
  private totalItems: BehaviorSubject<number>;
  private isLoading: BehaviorSubject<boolean>;

  constructor(private productsService: ProductsService) {
    this.items = new BehaviorSubject<ProductElement[]>([]);
    this.totalItems = new BehaviorSubject<number>(0);
    this.isLoading = new BehaviorSubject<boolean>(false);
    this.totalElements$ = this.totalItems.asObservable();
  }

  connect(): Observable<ProductElement[]> {
    return this.items.asObservable();
  }

  disconnect(): void {
    this.items.complete();
    this.totalItems.complete();
    this.isLoading.complete();
  }

  load(name: string | null,
       page: number,
       size: number,
       error: (err: any) => void): void {
    this.isLoading.next(true);

    this.productsService.get(name, page, size)
      .pipe(debounceTime(1000))
      .subscribe({
        next: page => {
          if (page.items === null) {
            this.items.next([]);
            this.totalItems.next(0);
            return;
          }

          this.items.next(page.items.map(product => {
            return <ProductElement>{
              id: product.id,
              name: product.name,
              price: product.price,
            }
          }));
          this.totalItems.next(page.totalItems);
        },
        error: err => error(err),
        complete: () => this.isLoading.next(false)
      });
  }
}
