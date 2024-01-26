import type { ApiResponse, Page, ProblemDetails } from '@/models/common'
import type { Product } from '@/models/product'

export class ProductService {
  public async getAll(
    name: string | null,
    page: number | null,
    size: number | null
  ): Promise<ApiResponse<Page<Product>>> {
    let query = ''
    if (page !== null) {
      query += `page=${page}`
    }

    if (size !== null) {
      if (query.length > 0) {
        query += '&'
      }
      query += `size=${size}`
    }

    if (name !== null && name.length > 0) {
      if (query.length > 0) {
        query += '&'
      }
      query += `name=${name}`
    }

    if (query.length > 0) {
      query = `?${query}`
    }

    const response = await fetch(`${import.meta.env.VITE_MEC_API}/api/products${query}`)

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
    const response = await fetch(`${import.meta.env.VITE_MEC_API}/api/products/${id}`)
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
    const response = await fetch(`${import.meta.env.MEC_API}/api/products`, {
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
    const response = await fetch(`${import.meta.env.MEC_API}/api/products/${product.id}`, {
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
    const response = await fetch(`${import.meta.env.MEC_API}/api/products/${id}`, {
      method: 'DELETE'
    })

    const error: ProblemDetails = await response.json()

    return {
      data: null,
      error: error,
      success: response.ok,
      response: response
    }
  }
}
