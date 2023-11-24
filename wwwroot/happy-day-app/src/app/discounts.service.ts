import {Injectable} from '@angular/core';
import {HttpClient, HttpParams} from "@angular/common/http";

import {Observable} from "rxjs";

import {environment} from "../environments/environment";
import {Page} from "./common";

@Injectable({
  providedIn: 'root'
})
export class DiscountsService {

  constructor(private httpClient: HttpClient) {
  }

  get(name: string | null,
      page: number,
      size: number): Observable<Page<Discount>> {
    let params = new HttpParams()
      .set('page', page)
      .set('size', size);

    if (name != null && name.trim().length > 0) {
      params = params.set('name', name.trim());
    }

    return this.httpClient.get<Page<Discount>>(`${environment.api}/api/discounts`, {params});
  }

  getById(id: string): Observable<Discount> {
    return this.httpClient.get<Discount>(`${environment.api}/api/discounts/${id}`);
  }

  create(discount: Discount): Observable<Discount> {
    return this.httpClient.post<Discount>(`${environment.api}/api/discounts`, discount);
  }

  update(id: string, discount: Discount): Observable<Discount> {
    return this.httpClient.put<Discount>(`${environment.api}/api/discounts/${id}`, discount);
  }

  delete(id: string): Observable<void> {
    return this.httpClient.delete<void>(`${environment.api}/api/discounts/${id}`);
  }

}

export interface Discount {
  id: string;
  name: string;
  price: number;
  products: DiscountProduct[];
  createAt: Date;
  updateAt: Date;
}

export interface DiscountProduct {
  id: string;
  name: string;
  quantity: number;
}
