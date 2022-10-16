import { Injectable } from '@angular/core';
import {Customer, CustomerSort} from "../models/customer";
import {Observable} from "rxjs";
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {Page} from "../models/page";

@Injectable({
  providedIn: 'root'
})
export class CustomerService {

  constructor(private httpClient: HttpClient) {}

  getAll(page: number, size: number, text: string, sort: CustomerSort | null): Observable<Page<Customer>> {
    let query = `page=${page}&size=${size}`;
    if(sort === null) {
      query += `&sort=${sort}`
    }

    if(text !== "") {
      query += `&text=${text}`
    }

    return this.httpClient.get<Page<Customer>>(`${environment.api}/api/v1/customers?${query}`)
  }

  get(id: string): Observable<Customer> {
    return this.httpClient.get<Customer>(`${environment.api}/api/v1/customers/${id}`)
  }

  create(customer: Customer): Observable<Customer> {
    return this.httpClient.post<Customer>(`${environment.api}/api/v1/customers`, customer)
  }

  update(id: string, customer: Customer): Observable<Customer> {
    return this.httpClient.put<Customer>(`${environment.api}/api/v1/customers/${id}`, customer)
  }

  delete(id: string): Observable<any> {
    return this.httpClient.delete<any>(`${environment.api}/api/v1/customers/${id}`)
  }
}
