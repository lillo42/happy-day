package reservation

import (
	"context"
	"testing"

	"happy_day/common"
	"happy_day/domain/reservation"
	"happy_day/infrastructure"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllReservations(t *testing.T) {
	req := infrastructure.ReservationFilter{
		Page: 0,
		Size: 0,
		Text: common.RandString(10),
	}

	repo := &infrastructure.MockReservationRepository{}
	repo.
		On("GetAll", mock.Anything, mock.Anything).
		Return(infrastructure.Page[reservation.State]{}, nil)

	handler := GetAllHandler{
		repository: repo,
	}

	_, err := handler.Handle(context.Background(), req)
	assert.Nil(t, err)
}
