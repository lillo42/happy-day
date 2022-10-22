import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {Observable} from "rxjs";
import {Page} from "../models/page";
import {environment} from "../../environments/environment";
import {Quote, QuoteResponse, Reservation, ReservationSort} from "../models/reservation";

@Injectable({
  providedIn: 'root'
})
export class ReservationService {

  constructor(private httpClient: HttpClient) { }

  getAll(page: number, size: number, text: string, sort: ReservationSort | null): Observable<Page<Reservation>> {
    let query = `page=${page}&size=${size}`;
    if (text !== "") {
      query += `&text=${text}`;
    }

    if(sort !== null) {
      query += `&sort=${sort}`;
    }

    return this.httpClient.get<Page<Reservation>>(`${environment.api}/api/v1/reservations?${query}`);
  }

  get(id: string): Observable<Reservation> {
    return this.httpClient.get<Reservation>(`${environment.api}/api/v1/reservations/${id}`);
  }

  create(reservation: Reservation): Observable<Reservation> {
    return this.httpClient.post<Reservation>(`${environment.api}/api/v1/reservations`, reservation);
  }

  update(id: string, reservation: Reservation): Observable<Reservation> {
    return this.httpClient.put<Reservation>(`${environment.api}/api/v1/reservations/${id}`, reservation);
  }

  delete(id: string): Observable<any> {
    return this.httpClient.delete<any>(`${environment.api}/api/v1/reservations/${id}`);
  }

  quote(quote: Quote): Observable<QuoteResponse> {
    return this.httpClient.post<QuoteResponse>(`${environment.api}/api/v1/reservations/quote`, quote);
  }
}
