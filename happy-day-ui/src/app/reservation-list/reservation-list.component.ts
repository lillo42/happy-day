import {AfterViewInit, Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {DataSource} from "@angular/cdk/collections";
import {BehaviorSubject, catchError, debounceTime, finalize, Observable, of, Subject, tap} from "rxjs";
import {ReservationService} from "../http-clients/reservation.service";
import {ReservationSort} from "../models/reservation";
import {MatDialog} from "@angular/material/dialog";
import {MatPaginator} from "@angular/material/paginator";
import {MatSort, Sort} from "@angular/material/sort";
import {ReservationBehavior, ReservationComponent, ReservationData} from "../reservation/reservation.component";

@Component({
  selector: 'app-reservation-list',
  templateUrl: './reservation-list.component.html',
  styleUrls: ['./reservation-list.component.scss']
})
export class ReservationListComponent implements OnInit, AfterViewInit {
  displayedColumns = ["id", "customerName", "price", "deliveryAt", "pickUpAt", "actions"];
  dataSource: ReservationDataSource;

  private sortBy = ReservationSort.IdAsc;
  private textChanged = new Subject<string>();

  @ViewChild("filter") filter: ElementRef | null = null;
  @ViewChild(MatPaginator) paginator: MatPaginator | null = null;
  @ViewChild(MatSort) sort: MatSort | null = null;

  constructor(private dialog: MatDialog,
              reservationService: ReservationService) {
    this.dataSource = new ReservationDataSource(reservationService);
  }

  ngOnInit(): void {
    this.textChanged
      .pipe(debounceTime(1000))
      .subscribe(() => this.reload());
  }

  ngAfterViewInit(): void {
    if(this.paginator == null) {
      return;
    }

    this.paginator.page
      .pipe(tap(() => this.reload()))
      .subscribe();
  }

  applyFilter(key: string): void {
    this.textChanged.next(key);
  }

  sortChanged(sort: Sort): void {
    if(this.paginator !== null) {
      this.paginator.pageIndex = 0;
    }

    if(sort.direction === "asc") {
      if(sort.active === "id") {
        this.sortBy = ReservationSort.IdAsc;
      }
    } else {
      if(sort.active === "id") {
        this.sortBy = ReservationSort.IdDesc;
      }
    }
  }

  openDialog(id: string, behavior: ReservationBehavior): void {
    const dialogRef = this.dialog.open(ReservationComponent, {
      width: '80%',
      data: <ReservationData>{
        behavior: behavior,
        id: id,
      },
    });

    dialogRef.afterClosed()
      .pipe(tap(() => this.reload()))
      .subscribe();
  }

  private reload(): void {
    if(this.paginator !== null && this.sort !== null && this.filter !== null) {
      this.dataSource.load(this.paginator.pageIndex, this.paginator.pageSize,
        this.filter?.nativeElement.value, this.sortBy);
    }
  }
}

interface ReservationViewModel {
  id: string;
  customerName: string;
  price: number;
  deliveryAt: Date;
  pickUpAt: Date;
}

class ReservationDataSource implements DataSource<ReservationViewModel> {
  private items = new BehaviorSubject<ReservationViewModel[]>([]);
  private totalElements = new BehaviorSubject<number>(0);
  private isLoading = new BehaviorSubject<boolean>(false);

  public totalElements$ = this.totalElements.asObservable();
  constructor(private reservationService: ReservationService) { }

  connect(): Observable<ReservationViewModel[]> {
    return this.items.asObservable();
  }

  disconnect(): void {
    this.items.complete();
    this.totalElements.complete();
    this.isLoading.complete();
  }

  load(pageIndex: number, pageSize: number, sortBy: ReservationSort | null, filter: string): void {
    this.isLoading.next(true);
    this.reservationService.getAll(pageIndex, pageSize, filter, sortBy)
      .pipe(
        tap(page => {
          this.totalElements.next(page.totalElements);
          this.items.next(page.items.map(reservation => <ReservationViewModel>{
            id: reservation.id,
            customerName: reservation.customer.name,
            price: reservation.price,
            deliveryAt: reservation.delivery.at,
            pickUpAt: reservation.pickUp.at,
          }));
        }),
        catchError(() => of([])),
        finalize(() => this.isLoading.next(false))
      )
      .subscribe();
  }
}
