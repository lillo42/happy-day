package application

import (
	"context"
	"happy_day/domain/reservation"
	"happy_day/infrastructure"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetReservationById(t *testing.T) {
	repo := &infrastructure.MockReservationRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(reservation.State{}, nil)

	handler := GetReservationByIdHandler{
		repository: repo,
	}

	_, err := handler.Handle(context.Background(), uuid.New())
	assert.Nil(t, err)
}
