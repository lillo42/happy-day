package apis

import (
	"happy_day/application"
	"happy_day/infrastructure"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CustomerController struct {
	getAllHandler         application.GetAllCustomersHandler
	getByIdHandler        application.GetCustomerByIdHandler
	createOrChangeHandler application.CreateOrChangeCustomerHandler
	deleteHandler         application.DeleteCustomerHandler
}

func (controller CustomerController) Routes(e *echo.Echo) {
	e.GET("/api/v1/customers", controller.getAll)
	e.POST("/api/v1/customers", controller.create)

	e.GET("/api/v1/customers/:id", controller.get)
	e.PUT("/api/v1/customers/:id", controller.update)
	e.DELETE("/api/v1/customers/:id", controller.delete)
}

func (controller CustomerController) getAll(ctx echo.Context) error {
	var filter infrastructure.CustomerFilter
	filter.Text = ctx.QueryParam("text")
	filter.Size, _ = strconv.ParseInt(ctx.Param("size"), 10, 64)
	filter.Page, _ = strconv.ParseInt(ctx.Param("page"), 10, 64)
	filter.SortBy = infrastructure.CustomerNameAsc

	params := ctx.QueryParams()
	if params.Has("sort") {
		filter.SortBy = infrastructure.CustomerSortBy(params.Get("sort"))
	}

	res, err := controller.getAllHandler.Handle(ctx.Request().Context(), filter)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func (controller CustomerController) create(ctx echo.Context) error {
	var req application.CreateOrChangeCustomerRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	req.Id = uuid.Nil
	res, err := controller.createOrChangeHandler.Handle(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, res)
}

func (controller CustomerController) update(ctx echo.Context) error {
	var req application.CreateOrChangeCustomerRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrCustomerNotFound
	}

	req.Id = id
	res, err := controller.createOrChangeHandler.Handle(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func (controller CustomerController) get(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrCustomerNotFound
	}

	res, err := controller.getByIdHandler.Handle(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func (controller CustomerController) delete(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrCustomerNotFound
	}

	req := application.DeleteCustomerRequest{Id: id}
	err = controller.deleteHandler.Handle(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
