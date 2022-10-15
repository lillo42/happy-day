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
	getAllHandler  application.GetAllReservationHandler
	getByIdHandler application.GetReservationByIdHandler
	createHandler  application.CreateReservationHandler
	changeHandler  application.ChangeReservationHandler
	deleteHandler  application.DeleteReservationHandler
	quoteHandler   application.QuoteReservationHandler
}

func (controller ReservationController) Routes(e *echo.Echo) {
	e.GET("/api/v1/reservations", controller.getAll)
	e.POST("/api/v1/reservations", controller.create)

	e.POST("/api/v1/reservations/quote", controller.quote)

	e.GET("/api/v1/reservations/:id", controller.get)
	e.PUT("/api/v1/reservations/:id", controller.update)
	e.DELETE("/api/v1/reservations/:id", controller.delete)
}

func (controller ReservationController) create(ctx echo.Context) error {
	var req application.CreateReservationRequest

	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	res, err := controller.createHandler.Handle(ctx.Request().Context(), req)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, res)
}

func (controller ReservationController) update(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrReservationNotFound
	}

	var req application.ChangeReservationRequest

	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	req.Id = id
	res, err := controller.changeHandler.Handle(ctx.Request().Context(), req)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func (controller ReservationController) delete(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrReservationNotFound
	}

	err = controller.deleteHandler.Handle(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (controller ReservationController) get(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrReservationNotFound
	}

	res, err := controller.getByIdHandler.Handle(ctx.Request().Context(), id)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func (controller ReservationController) getAll(ctx echo.Context) error {
	var filter infrastructure.ReservationFilter
	filter.Text = ctx.Param("text")
	filter.Size, _ = strconv.ParseInt(ctx.Param("size"), 10, 64)
	filter.Page, _ = strconv.ParseInt(ctx.Param("page"), 10, 64)
	filter.SortBy = infrastructure.ReservationOrderBy(ctx.Param("orderBy"))

	res, err := controller.getAllHandler.Handle(ctx.Request().Context(), filter)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func (controller ReservationController) quote(ctx echo.Context) error {
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
