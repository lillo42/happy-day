package apis

import (
	"net/http"
	"strconv"

	"happy_day/application/reservation"
	"happy_day/infrastructure"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func MapReservationEndpoints(e *echo.Echo) {
	e.GET("/api/v1/reservations", getAllReservation)
	e.POST("/api/v1/reservations", createReservation)

	e.POST("/api/v1/reservations/quote", quoteReservation)

	e.GET("/api/v1/reservations/:id", getReservationById)
	e.PUT("/api/v1/reservations/:id", updateReservation)
	e.DELETE("/api/v1/reservations/:id", deleteReservation)
}

func createReservation(ctx echo.Context) error {
	var req reservation.CreateRequest

	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	res, err := initializeCreateReservationHandler().Handle(ctx.Request().Context(), req)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, res)
}

func updateReservation(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrReservationNotFound
	}

	var req reservation.ChangeRequest

	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	req.Id = id
	res, err := initializeChangeReservationHandler().Handle(ctx.Request().Context(), req)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func deleteReservation(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrReservationNotFound
	}

	err = initializeDeleteReservationHandler().Handle(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func getReservationById(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrReservationNotFound
	}

	res, err := initializeGetReservationByIdHandler().Handle(ctx.Request().Context(), id)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func getAllReservation(ctx echo.Context) error {
	var filter infrastructure.ReservationFilter
	filter.Text = ctx.Param("text")
	filter.Size, _ = strconv.ParseInt(ctx.Param("size"), 10, 64)
	filter.Page, _ = strconv.ParseInt(ctx.Param("page"), 10, 64)
	filter.SortBy = infrastructure.ReservationOrderBy(ctx.Param("orderBy"))

	res, err := initializeGetAllReservationHandler().Handle(ctx.Request().Context(), filter)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func quoteReservation(ctx echo.Context) error {
	var req reservation.QuoteRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	res, err := initializeQuoteReservationHandler().Handler(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}
