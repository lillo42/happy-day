import type { ApiResponse, Page, ProblemDetails } from '@/models/common'
import type { Discount } from '@/models/discount'

export class DiscountService {
  private readonly api: string

  constructor() {
    this.api = `${import.meta.env.VITE_MEC_API}/api/discounts`
  }

  public async fetchAll(
    name: string | null,
    page: number | null,
    size: number | null
  ): Promise<ApiResponse<Page<Discount>>> {
    const params = new URLSearchParams()
    if (page !== null) {
      params.set('page', page.toString())
    }

    if (size !== null) {
      params.set('size', size.toString())
    }

    if (name !== null && name.length > 0) {
      params.set('name', name)
    }

    const response = await fetch(`${this.api}?${params.toString()}`)

    if (!response.ok) {
      const error: ProblemDetails = await response.json()
      return {
        data: null,
        success: false,
        error: error,
        response: response
      }
    }

    const discounts: Page<Discount> = await response.json()
    return {
      data: discounts,
      error: null,
      success: true,
      response: response
    }
  }

  public async fetch(id: string): Promise<ApiResponse<Discount>> {
    const response = await fetch(`${this.api}/${id}`)
    if (!response.ok) {
      const error: ProblemDetails = await response.json()
      return {
        data: null,
        error: error,
        success: false,
        response: response
      }
    }

    const discount: Discount = await response.json()
    return {
      data: discount,
      error: null,
      success: true,
      response: response
    }
  }

  public async create(discount: Discount): Promise<ApiResponse<Discount>> {
    const response = await fetch(this.api, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(discount)
    })

    if (!response.ok) {
      const error: ProblemDetails = await response.json()
      return {
        data: null,
        error: error,
        success: false,
        response: response
      }
    }

    discount = await response.json()
    return {
      data: discount,
      error: null,
      success: true,
      response: response
    }
  }

  public async update(discount: Discount): Promise<ApiResponse<Discount>> {
    const response = await fetch(`${this.api}/${discount.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(discount)
    })

    if (!response.ok) {
      const error: ProblemDetails = await response.json()
      return {
        data: null,
        error: error,
        success: false,
        response: response
      }
    }

    discount = await response.json()
    return {
      data: discount,
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
}
