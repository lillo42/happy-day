export interface Customer {
  id: string
  name: string
  comment: string | null
  phones: string[]
  createAt: Date
  updateAt: Date
}
