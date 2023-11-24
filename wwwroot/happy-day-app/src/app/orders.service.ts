import {HttpClient, HttpParams} from "@angular/common/http";
import {Injectable} from '@angular/core';

import {Observable} from "rxjs";

import {environment} from "../environments/environment";
import {Page} from "./common";

@Injectable({
  providedIn: 'root'
})
export class OrdersService {

  constructor(private httpClient: HttpClient) {
  }

  get(name: string | null,
      address: string | null,
      customerName: string | null,
      customerPhone: string | null,
      page: number,
      size: number): Observable<Page<Order>> {
    let params = new HttpParams()
      .set('page', page)
      .set('size', size);

    if (name != null && name.trim().length > 0) {
      params = params.set('name', name.trim());
    }

    if (address != null && address.trim().length > 0) {
      params = params.set('address', address.trim());
    }

    if (customerName != null && customerName.trim().length > 0) {
      params = params.set('customerName', customerName.trim());
    }

    if (customerPhone != null && customerPhone.trim().length > 0) {
      params = params.set('customerPhone', customerPhone.trim());
    }

    return this.httpClient.get<Page<Order>>(`${environment.api}/api/orders`, {params});
  }

  getById(id: string): Observable<Order> {
    return this.httpClient.get<Order>(`${environment.api}/api/orders/${id}`);
  }

  create(req: OrderCreateOrChange): Observable<Order> {
    return this.httpClient.post<Order>(`${environment.api}/api/orders`, req);
  }

  update(id: string, req: OrderCreateOrChange): Observable<Order> {
    return this.httpClient.put<Order>(`${environment.api}/api/orders/${id}`, req);
  }

  delete(id: string): Observable<void> {
    return this.httpClient.delete<void>(`${environment.api}/api/orders/${id}`);
  }

  quote(quote: OrderQuote): Observable<OrderQuoteResponse> {
    return this.httpClient.post<OrderQuoteResponse>(`${environment.api}/api/orders/quote`, quote);
  }
}

export interface Order {
  id: string;
  address: string;
  comment: string;

  deliveryAt: Date;
  pickUpAt: Date;

  totalPrice: number;
  discount: number;
  finalPrice: number;

  customer: OrderCustomer;
  products: OrderProduct[];
  payments: OrderPayment[];

  createAt: Date;
  updateAt: Date;
}

export interface OrderCustomer {
  id: string;
  name: string;
  comment: string;
  phones: string[];
}

export interface OrderProduct {
  id: string;
  name: string;
  quantity: number;
  price: number;
}

export interface OrderPayment {
  info: string;
  value: number;
  at: Date;
  method: "pix" | "bank-transfer" | "cash";
}

export interface OrderQuote {
  products: OrderProduct[];
}

export interface OrderQuoteResponse {
  totalPrice: number;
}

export interface OrderCreateOrChange {
  address: string;
  comment: string;

  deliveryAt: Date;
  pickUpAt: Date;

  totalPrice: number;
  discount: number;
  finalPrice: number;

  customerId: string;

  products: OrderProduct[];
  payments: OrderPayment[];
}

