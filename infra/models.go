package infra

type (
	ProblemDetails struct {
		Status  int    `json:"status"`
		Title   string `json:"title"`
		Details string `json:"details,omitempty"`
		Type    string `json:"type,omitempty"`
	}

	Page[T any] struct {
		Items      []T   `json:"items"`
		TotalPage  int64 `json:"totalPage"`
		TotalItems int64 `json:"totalItems"`
	}
)
