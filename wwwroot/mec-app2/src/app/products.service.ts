import {HttpClient, HttpParams} from "@angular/common/http";
import {Injectable} from '@angular/core';

import {Observable} from "rxjs";

import {Page} from "./common";
import {environment} from "../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class ProductsService {

  constructor(private httpClient: HttpClient) {
  }

  get(name: string | null,
      page: number,
      size: number): Observable<Page<Product>> {
    let params = new HttpParams()
      .set('page', page)
      .set('size', size);

    if (name != null && name.trim().length > 0) {
      params = params.set('name', name.trim());
    }

    return this.httpClient.get<Page<Product>>(`${environment.api}/api/products`, {params});
  }

  getById(id: string): Observable<Product> {
    return this.httpClient.get<Product>(`${environment.api}/api/products/${id}`);
  }

  create(product: Product): Observable<Product> {
    return this.httpClient.post<Product>(`${environment.api}/api/products`, product);
  }

  update(id: string, product: Product): Observable<Product> {
    return this.httpClient.put<Product>(`${environment.api}/api/products/${id}`, product);
  }

  delete(id: string): Observable<void> {
    return this.httpClient.delete<void>(`${environment.api}/api/products/${id}`);
  }
}

export interface Product {
  id: string;
  name: string;
  price: number;
  createAt: Date;
  updateAt: Date;
}
