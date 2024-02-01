import type { ApiResponse, Page, ProblemDetails } from '@/models/common'
import type { Customer } from '@/models/customer'

export class CustomerService {
  private readonly api: string

  constructor() {
    this.api = `${import.meta.env.VITE_MEC_API}/api/customers`
  }

  public async fetchAll(
    name: string | null,
    comment: string | null,
    phone: string | null,
    page: number | null,
    size: number | null
  ): Promise<ApiResponse<Page<Customer>>> {
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
    if (comment !== null && comment.length > 0) {
      params.set('comment', comment)
    }
    if (phone !== null && phone.length > 0) {
      params.set('phone', phone)
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
    const customers: Page<Customer> = await response.json()
    return {
      data: customers,
      error: null,
      success: true,
      response: response
    }
  }

  public async fetch(id: string): Promise<ApiResponse<Customer>> {
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
    const customer: Customer = await response.json()
    return {
      data: customer,
      error: null,
      success: true,
      response: response
    }
  }

  public async create(customer: Customer): Promise<ApiResponse<Customer>> {
    const response = await fetch(this.api, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(customer)
    })

    if (!response.ok) {
      const error: ProblemDetails = await response.json()
      return {
        data: null,
        success: false,
        error: error,
        response: response
      }
    }

    const entity: Customer = await response.json()
    return {
      data: entity,
      error: null,
      success: true,
      response: response
    }
  }

  public async update(customer: Customer): Promise<ApiResponse<Customer>> {
    const response = await fetch(`${this.api}/${customer.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(customer)
    })

    if (!response.ok) {
      const error: ProblemDetails = await response.json()
      return {
        data: null,
        success: false,
        error: error,
        response: response
      }
    }

    const entity: Customer = await response.json()
    return {
      data: entity,
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
