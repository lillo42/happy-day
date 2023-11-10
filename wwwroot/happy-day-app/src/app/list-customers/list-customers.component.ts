import {AfterViewInit, Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import { MatTableDataSource } from "@angular/material/table";
import { MatPaginator } from "@angular/material/paginator";
import {CustomersService} from "../customers.service";
import {MatSelect} from "@angular/material/select";
import {MatInput} from "@angular/material/input";
import {debounce, debounceTime, merge, mergeAll} from "rxjs";
import {Router} from "@angular/router";

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
              private router: Router) {
    this.dataSource = new MatTableDataSource<CustomerElement>([]);
  }

  ngAfterViewInit(): void {
    this.dataSource.paginator = this.paginator;
    this.loadCustomers();
  }

  delete(id: string): void {
  }

  edit(id: string): Promise<boolean> {
    return this.router.navigateByUrl(`/customers/${id}`);
  }

  add(): void {

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
      .subscribe((page) => {
        if(page.items === null || page.items.length == 0) {
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
