import {AfterViewInit, Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {CollectionViewer, DataSource} from "@angular/cdk/collections";
import {BehaviorSubject, catchError, debounceTime, finalize, Observable, of, Subject, tap} from "rxjs";
import {ProductService} from "../http-clients/product.service";
import {ProductSort} from "../models/product";
import {MatDialog} from "@angular/material/dialog";
import {MatPaginator} from "@angular/material/paginator";
import {MatSort, Sort} from "@angular/material/sort";
import {ProductBehavior, ProductComponent, ProductData} from "../product/product.component";

@Component({
  selector: 'app-product-list',
  templateUrl: './product-list.component.html',
  styleUrls: ['./product-list.component.scss']
})
export class ProductListComponent implements OnInit, AfterViewInit {

  displayedColumns = ["id",  "name", "price", "actions"];
  dataSource: ProductDataSource;

  private sortBy = ProductSort.NameAsc;
  private textChanged = new Subject<string>();

  @ViewChild("filter") filter: ElementRef | null = null;
  @ViewChild(MatPaginator) paginator: MatPaginator | null = null;
  @ViewChild(MatSort) sort: MatSort | null = null;

  constructor(private dialog: MatDialog,
              productService: ProductService) {
    this.dataSource = new ProductDataSource(productService);
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

  openDialog(id: string, behavior: ProductBehavior): void {
    const dialogRef = this.dialog.open(ProductComponent, {
      width: '80%',
      data: <ProductData>{
        behavior: behavior,
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
        this.sortBy = ProductSort.IdAsc;
      } else {
        this.sortBy = ProductSort.IdDesc;
      }
    } else if (sort.active == "name") {
      if (sort.direction == "asc") {
        this.sortBy = ProductSort.NameAsc;
      } else {
        this.sortBy = ProductSort.NameDesc;
      }
    } else if (sort.active == "price") {
      if (sort.direction == "asc") {
        this.sortBy = ProductSort.PriceAsc;
      } else {
        this.sortBy = ProductSort.PriceDesc;
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

interface ProductViewModel {
  id: string;
  name: string;
  price: number;
}

class ProductDataSource implements DataSource<ProductViewModel> {
  private items = new BehaviorSubject<ProductViewModel[]>([]);
  private totalElements = new BehaviorSubject<number>(0);
  private isLoading = new BehaviorSubject<boolean>(false);

  public totalElements$ = this.totalElements.asObservable();
  constructor(private customerService: ProductService) { }

  connect(): Observable<ProductViewModel[]> {
    return this.items.asObservable();
  }

  disconnect(collectionViewer: CollectionViewer): void {
    this.items.complete();
    this.totalElements.complete();
    this.isLoading.complete();
  }

  get(page: number, size: number, text: string, sort: ProductSort | null): void {
    this.isLoading.next(true);
    this.customerService.getAll(page, size, text, sort)
      .pipe(
        tap(page => {
          this.totalElements.next(page.totalElements);
          this.items.next(page.items.map(product => <ProductViewModel>{
            id: product.id,
            name: product.name,
            price: product.price,
          }));
        }),
        catchError(() => of([])),
        finalize(() => this.isLoading.next(false))
      ).subscribe();
  }


}
