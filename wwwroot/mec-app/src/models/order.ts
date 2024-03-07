export interface Order {
  id: string
  address: string
  comment: string | null

  deliveryAt: Date
  pickUpAt: Date

  totalPrice: number
  discount: number
  finalPrice: number

  customer: OrderCustomer
  products: OrderProduct[]
  payments: OrderPayment[]

  createAt: Date
  updateAt: Date
}

export interface OrderProduct {
  id: string
  name: string
  quantity: number
  price: number
}

export interface OrderPayment {
  amount: number
  at: Date
  method: 'pix' | 'bank-transfer' | 'cash'
  info: string | null
}

export interface OrderCustomer {
  id: string
  name: string
  comment: string | null
  phones: string[]
}
