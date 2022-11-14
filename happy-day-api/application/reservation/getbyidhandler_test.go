package reservation

import (
	"context"
	"testing"

	"happy_day/domain/reservation"
	"happy_day/infrastructure"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetReservationById(t *testing.T) {
	repo := &infrastructure.MockReservationRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(reservation.State{}, nil)

	handler := GetByIdHandler{
		repository: repo,
	}

	_, err := handler.Handle(context.Background(), uuid.New())
	assert.Nil(t, err)
}
