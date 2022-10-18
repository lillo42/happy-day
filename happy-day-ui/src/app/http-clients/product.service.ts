import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {Product, ProductSort} from "../models/product";
import {Page} from "../models/page";
import {environment} from "../../environments/environment";


@Injectable({
  providedIn: 'root'
})
export class ProductService {

  constructor(private httpClient: HttpClient) { }

  getAll(page: number, size: number, text: string, sort: ProductSort | null): Observable<Page<Product>> {
    let query = `page=${page}&size=${size}`;
    if(sort === null) {
      query += `&sort=${sort}`
    }

    if(text !== "") {
      query += `&text=${text}`
    }

    return this.httpClient.get<Page<Product>>(`${environment.api}/api/v1/products?${query}`)
  }

  get(id: string): Observable<Product> {
    return this.httpClient.get<Product>(`${environment.api}/api/v1/products/${id}`);
  }

  create(customer: Product): Observable<Product> {
    return this.httpClient.post<Product>(`${environment.api}/api/v1/products`, customer);
  }

  update(id: string, customer: Product): Observable<Product> {
    return this.httpClient.put<Product>(`${environment.api}/api/v1/products/${id}`, customer);
  }

  delete(id: string): Observable<any> {
    return this.httpClient.delete<any>(`${environment.api}/api/v1/products/${id}`);
  }
}
