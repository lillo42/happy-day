import {Injectable} from '@angular/core';
import {HttpClient, HttpParams} from "@angular/common/http";
import {Observable} from "rxjs";
import {environment} from "../environments/environment";
import {Page} from "./common";

@Injectable({
  providedIn: 'root'
})
export class CustomersService {

  constructor(private httpClient: HttpClient) {
  }

  get(name: string | null,
      phone: string | null,
      comment: string | null,
      page: number,
      size: number): Observable<Page<Customer>> {
    let params = new HttpParams()
      .set('page', page)
      .set('size', size);

    if (name != null && name.trim().length > 0) {
      params = params.set('name', name.trim());
    }

    if (phone != null && phone.trim().length > 0) {
      params = params.set('phone', phone.trim());
    }

    if (comment != null && comment.trim().length > 0) {
      params = params.set('comment', comment.trim());
    }

    return this.httpClient.get<Page<Customer>>(`${environment.api}/api/customers`, {params});
  }

  getById(id: string): Observable<Customer> {
    return this.httpClient.get<Customer>(`${environment.api}/api/customers/${id}`);
  }

  create(customer: Customer): Observable<Customer> {
    return this.httpClient.post<Customer>(`${environment.api}/api/customers`, customer);
  }

  update(id: string, customer: Customer): Observable<Customer> {
    return this.httpClient.put<Customer>(`${environment.api}/api/customers/${id}`, customer);
  }

  delete(id: string): Observable<void> {
    return this.httpClient.delete<void>(`${environment.api}/api/customers/${id}`);
  }
}

export interface Customer {
  id: string;
  name: string;
  comment: string;
  phones: string[];
  pix: string;
  createAt: Date;
  updateAt: Date;
}
