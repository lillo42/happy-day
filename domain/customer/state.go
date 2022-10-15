package customer

import (
	"time"

	"github.com/google/uuid"
)

type (
	State struct {
		Id       uuid.UUID `json:"id,omitempty" bson:"id"`
		Name     string    `json:"name" bson:"name"`
		Comment  string    `json:"comment,omitempty" bson:"comment"`
		Phones   []Phone   `json:"phones" bson:"phones"`
		CreateAt time.Time `json:"createAt" bson:"createAt"`
		ModifyAt time.Time `json:"modifyAt" bson:"modifyAt"`
	}

	Phone struct {
		Number string `json:"number" bson:"number"`
	}
)
