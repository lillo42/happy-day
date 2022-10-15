package application

import (
	"happy_day/infrastructure"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteReservationHandler(t *testing.T) {
	repo := &infrastructure.MockReservationRepository{}
	repo.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)

	handler := DeleteReservationHandler{repository: repo}
	err := handler.Handle(nil, uuid.New())
	assert.Nil(t, err)
}
