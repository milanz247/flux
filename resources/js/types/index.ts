// TypeScript mirrors of the Go DTOs (app/dto). Keep these in sync — every
// req.View() call serializes one of the Go DTOs into these shapes.

export interface UserDTO {
  id: number
  name: string
  email: string
  verified: boolean
  createdAt: string
}

export interface PaginationDTO {
  page: number
  perPage: number
  total: number
  totalPages: number
}

export interface DashboardStatDTO {
  label: string
  value: string
  hint: string
}
