package apis

import (
	"happy_day/application"
	"happy_day/infrastructure"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReservationController struct {
	quoteHandler application.QuoteReservationHandler
}

func (controller ReservationController) Routes(e *echo.Echo) {
	e.GET("/api/v1/reservations", getAllReservations)
	e.POST("/api/v1/reservations", createReservation)
	e.POST("/api/v1/reservations/quote", controller.quoteReservation)
	e.GET("/api/v1/reservations/:id", getReservationById)
	e.PUT("/api/v1/reservations/:id", changeReservation)
	e.DELETE("/api/v1/reservations/:id", deleteReservation)
}

func createReservation(ctx echo.Context) error {
	var req application.CreateReservationRequest

	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	handler := initCreateReservationHandler()
	res, err := handler.Handler(ctx.Request().Context(), req)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, res)
}

func changeReservation(ctx echo.Context) error {
	id, _ := uuid.Parse(ctx.Param("id"))

	var req application.ChangeReservationRequest

	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	req.Id = id
	handler := initChangeReservationHandler()
	res, err := handler.Handler(ctx.Request().Context(), req)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func deleteReservation(ctx echo.Context) error {
	id, _ := uuid.Parse(ctx.Param("id"))

	handler := initDeleteReservationHandler()
	err := handler.Handler(ctx.Request().Context(), id)

	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func getReservationById(ctx echo.Context) error {
	id, _ := uuid.Parse(ctx.Param("id"))

	handler := initGetReservationByIdHandler()
	res, err := handler.Handler(ctx.Request().Context(), id)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func getAllReservations(ctx echo.Context) error {
	var filter infrastructure.ReservationFilter
	filter.Text = ctx.Param("text")
	filter.Size, _ = strconv.ParseInt(ctx.Param("size"), 10, 64)
	filter.Page, _ = strconv.ParseInt(ctx.Param("page"), 10, 64)
	filter.OrderBy = infrastructure.ReservationOrderBy(ctx.Param("orderBy"))

	handler := initGetAllReservationHandler()
	res, err := handler.Handler(ctx.Request().Context(), filter)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func (controller ReservationController) quoteReservation(ctx echo.Context) error {
	var req application.QuoteReservationRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	res, err := controller.quoteHandler.Handler(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}
