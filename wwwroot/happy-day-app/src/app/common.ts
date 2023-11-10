export class Page<T> {
  items: T[] | null = null;
  totalItems: number = 0;
  totalPages: number = 0;
}

export interface  ProblemDetails {
  type: string;
  message: string;
  status: number;
}
