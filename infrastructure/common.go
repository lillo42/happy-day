package infrastructure

const Database = "happy-day"

type Page[T any] struct {
	Items         []T   `json:"items"`
	TotalPages    int64 `json:"totalPages"`
	TotalElements int64 `json:"totalElements"`
}
