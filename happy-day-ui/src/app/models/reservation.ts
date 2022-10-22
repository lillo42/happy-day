import {Customer} from "./customer";

export interface Reservation {
  id: string;
  price: number;
  discount: number;
  finalPrice: number;
  comment: string;
  products: Product[];
  delivery: DeliveryOrPickUp;
  pickUp: DeliveryOrPickUp;
  paymentInstallments: PaymentInstallment[];
  customer: Customer;
  address: Address;
  createdAt: Date;
  modifiedAt: Date;
}

export interface Product {
  id: string;
  price: number;
  quantity: number;
}

export interface DeliveryOrPickUp {
  at: Date;
  by: string[];
}

export interface PaymentInstallment {
  at: Date;
  amount: number;
  method: PaymentMethod;
}

export enum PaymentMethod {
  Pix = "pix",
  BankTransfer = "bankTransfer",
  Cash = "cash",
}

export interface Address {
  street: string;
  number: string;
  neighborhood: string;
  complement: string;
  postalCode: string;
  city: string;
}

export enum ReservationSort {
  IdAsc = "id_asc",
  IdDesc = "id_desc",
}

export interface Quote {
  products: Product[];
}

export interface QuoteResponse {
  price: number;
}
