export interface Product {
  id: string;
  name: string;
  price: number;
  products: InnerProduct[];
  createAt: Date;
  modifyAt: Date;
}

export interface InnerProduct {
  id: string;
  quantity: number;
}

export enum ProductSort {
  IdAsc = "id_asc",
  IdDesc = "id_desc",
  NameAsc = "name_asc",
  NameDesc = "name_desc",
  PriceAsc = "price_asc",
  PriceDesc = "price_desc",
}
