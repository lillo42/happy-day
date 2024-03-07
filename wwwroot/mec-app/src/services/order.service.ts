import type { ApiResponse, Page, ProblemDetails } from '@/models/common'
import type { Order, OrderPayment, OrderProduct } from '@/models/order'

export class OrderService {
  private readonly api: string

  constructor() {
    this.api = `${import.meta.env.VITE_MEC_API}/api/orders`
  }

  public async fetchAll(
    comment: string | null,
    address: string | null,
    customerName: string | null,
    customerPhone: string | null,
    page: number,
    size: number
  ): Promise<ApiResponse<Page<Order>>> {
    const params = new URLSearchParams()
    params.set('page', page.toString())
    params.set('size', size.toString())

    if (comment !== null && comment.length > 0) {
      params.set('comment', comment)
    }

    if (address !== null && address.length > 0) {
      params.set('address', address)
    }

    if (customerName !== null && customerName.length > 0) {
      params.set('customerName', customerName)
    }

    if (customerPhone !== null && customerPhone.length > 0) {
      params.set('customerPhone', customerPhone)
    }

    const response = await fetch(`${this.api}?${params.toString()}`)
    if (!response.ok) {
      const error = await response.json()
      return {
        data: null,
        success: false,
        error: error,
        response: response
      }
    }

    const orders: Page<Order> = await response.json()
    return {
      data: orders,
      error: null,
      success: true,
      response: response
    }
  }

  public async fetch(id: string): Promise<ApiResponse<Order>> {
    const response = await fetch(`${this.api}/${id}`)
    if (!response.ok) {
      const error = await response.json()
      return {
        data: null,
        success: false,
        error: error,
        response: response
      }
    }

    const orders: Order = await response.json()
    return {
      data: orders,
      error: null,
      success: true,
      response: response
    }
  }

  public async create(request: OrderCreateOrChangeRequest): Promise<ApiResponse<Order>> {
    const response = await fetch(this.api, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(request)
    })

    if (!response.ok) {
      const error = await response.json()
      return {
        data: null,
        success: false,
        error: error,
        response: response
      }
    }

    const order: Order = await response.json()
    return {
      data: order,
      error: null,
      success: true,
      response: response
    }
  }

  public async update(
    id: string,
    request: OrderCreateOrChangeRequest
  ): Promise<ApiResponse<Order>> {
    const response = await fetch(`${this.api}/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(request)
    })

    if (!response.ok) {
      const error = await response.json()
      return {
        data: null,
        success: false,
        error: error,
        response: response
      }
    }

    const order: Order = await response.json()
    return {
      data: order,
      error: null,
      success: true,
      response: response
    }
  }
  public async delete(id: string): Promise<ApiResponse<any>> {
    const response = await fetch(`${this.api}/${id}`, {
      method: 'DELETE'
    })

    const error: ProblemDetails | null = response.ok ? null : await response.json()
    return {
      data: null,
      error: error,
      success: response.ok,
      response: response
    }
  }

  public async quote(products: OrderProduct[]): Promise<ApiResponse<number>> {
    const response = await fetch(`${this.api}/quote`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ products: products })
    })

    if (!response.ok) {
      const error = await response.json()
      return {
        data: null,
        success: false,
        error: error,
        response: response
      }
    }
    const price: { totalPrice: number } = await response.json()
    return {
      data: price.totalPrice,
      error: null,
      success: true,
      response: response
    }
  }
}

export interface OrderCreateOrChangeRequest {
  address: string
  comment: string

  deliveryAt: Date
  pickUpAt: Date

  totalPrice: number
  discount: number
  finalPrice: number

  customerId: string

  products: OrderProduct[]
  payments: OrderPayment[]
}
