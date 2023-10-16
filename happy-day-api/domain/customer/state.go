package customer

import (
	"time"
)

type (
	State struct {
		ID        uint      `json:"id"`
		Name      string    `json:"name"`
		Comment   string    `json:"comment,omitempty"`
		Phones    []Phone   `json:"phones"`
		CreatedAt time.Time `json:"createdAt"`
		UpdateAt  time.Time `json:"updatedAt"`
	}

	Phone struct {
		Number string `json:"number"`
	}
)
