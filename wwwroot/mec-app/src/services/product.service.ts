import type { ApiResponse, Page, ProblemDetails } from '@/models/common'
import type { Product } from '@/models/product'

export class ProductService {
  private readonly api: string

  constructor() {
    this.api = `${import.meta.env.VITE_MEC_API}/api/products`
  }

  public async getAll(
    name: string | null,
    page: number | null,
    size: number | null
  ): Promise<ApiResponse<Page<Product>>> {
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

    const products: Page<Product> = await response.json()
    return {
      data: products,
      error: null,
      success: true,
      response: response
    }
  }

  public async get(id: string): Promise<ApiResponse<Product>> {
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

    const product: Product = await response.json()
    return {
      data: product,
      error: null,
      success: true,
      response: response
    }
  }

  public async create(product: Product): Promise<ApiResponse<Product>> {
    const response = await fetch(this.api, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(product)
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

    const entity: Product = await response.json()
    return {
      data: entity,
      error: null,
      success: true,
      response: response
    }
  }

  public async update(product: Product): Promise<ApiResponse<Product>> {
    const response = await fetch(`${this.api}/${product.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(product)
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

    const entity: Product = await response.json()
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