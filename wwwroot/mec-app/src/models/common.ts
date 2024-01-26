export interface DropdownOption {
  label: string
  value: string
}

export interface Page<T> {
  items: T[] | null
  totalItems: number
  totalPages: number
}

export interface ApiResponse<T> {
  data: T | null
  error: ProblemDetails | null
  response: Response
  success: boolean
}

export interface ProblemDetails {
  type: string
  message: string
  status: number
}
