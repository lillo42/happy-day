import { AfterViewInit, Component, ElementRef, ViewChild } from '@angular/core';
import { Router } from "@angular/router";

import { MatDialog } from "@angular/material/dialog";
import { MatPaginator } from "@angular/material/paginator";
import { MatSelect } from "@angular/material/select";
import { MatSnackBar } from "@angular/material/snack-bar";
import { MatTableDataSource } from "@angular/material/table";

import { debounceTime } from "rxjs";

import { CustomersService } from "../customers.service";
import { CustomerDeleteComponent } from "../customer-delete/customer-delete.component";

@Component({
  selector: 'app-list-customers',
  templateUrl: './list-customers.component.html',
  styleUrls: ['./list-customers.component.scss']
})
export class ListCustomersComponent implements AfterViewInit {
  displayedColumns: string[] = ['id', 'name', 'comment', 'phones', 'pix', 'actions'];
  dataSource: MatTableDataSource<CustomerElement>;

  @ViewChild(MatPaginator) paginator: MatPaginator | null = null;
  @ViewChild('selectField') field: MatSelect | null = null;
  @ViewChild('inputFilter') filter: ElementRef | null = null;

  constructor(private customersService: CustomersService,
              private router: Router,
              private snack: MatSnackBar,
              private dialog: MatDialog) {
    this.dataSource = new MatTableDataSource<CustomerElement>([]);
  }

  ngAfterViewInit(): void {
    this.dataSource.paginator = this.paginator;
    this.loadCustomers();
  }

  delete(id: string): void {
    this.dialog.open(CustomerDeleteComponent, {data: { id: id }})
      .afterClosed()
      .subscribe(() => this.loadCustomers());
  }

  edit(id: string): Promise<boolean> {
    return this.router.navigateByUrl(`/customers/${id}`);
  }

  loadCustomers(): void {
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

    this.customersService.get(name, phone, comment, page, size)
      .pipe(debounceTime(1000))
      .subscribe({
        next: page => {
          if (page.items === null) {
            this.dataSource.data = [];
            return;
          }

          this.dataSource.data = page.items.map(customer => {
            return <CustomerElement>{
              id: customer.id,
              name: customer.name,
              comment: customer.comment,
              pix: customer.pix,
              phones: customer.phones.join(', ')
            }
          });
        },
        error: err => this.snack.open(err.message, 'OK')
      });
  }
}

export interface CustomerElement {
  id: string;
  name: string;
  comment: string;
  pix: string;
  phones: string;
}
