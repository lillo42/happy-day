package common

type (
	Link struct {
		Href   string `json:"href,omitempty"`
		Method string `json:"method,omitempty"`
		Rel    string `json:"rel,omitempty"`
	}

	Hateoas struct {
		Links []Link `json:"links,omitempty"`
	}

	Page[T any] struct {
		TotalElements int64
		TotalPages    int64
		Items         []T
	}
)
