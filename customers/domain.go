package customers

import (
	"github.com/google/uuid"
	"time"
)

type Customer struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Comment  string    `json:"comment,omitempty"`
	Pix      string    `json:"pix,omitempty"`
	Phones   []string  `json:"phones,omitempty"`
	CreateAt time.Time `json:"createAt"`
	UpdateAt time.Time `json:"updateAt"`
	Version  uint      `json:"-"`
}
