export interface Discount {
  id: string
  name: string
  products: DiscountProduct[]
  price: number
  createAt: Date
  updateAt: Date
}

export interface DiscountProduct {
  id: string
  name: string
  quantity: number
}
