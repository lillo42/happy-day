package apis

import (
	"github.com/labstack/echo/v4"
	"happy_day/infrastructure"
	"net/http"
	"strconv"

	"happy_day/application/customer"
)

func MapCustomerEndpoints(e *echo.Echo) {
	e.GET("/api/v1/customers", getAllCustomers)
	e.POST("/api/v1/customers", createCustomer)

	e.GET("/api/v1/customers/:id", getCustomerById)
	e.PUT("/api/v1/customers/:id", updateCustomer)
	e.DELETE("/api/v1/customers/:id", deleteCustomer)
}

func getAllCustomers(ctx echo.Context) error {
	var filter infrastructure.CustomerFilter
	filter.Name = ctx.QueryParam("name")
	filter.Comment = ctx.QueryParam("comment")
	filter.Size, _ = strconv.ParseInt(ctx.QueryParam("size"), 10, 64)
	filter.Page, _ = strconv.ParseInt(ctx.QueryParam("page"), 10, 64)
	filter.SortBy = infrastructure.CustomerNameAsc

	params := ctx.QueryParams()
	if params.Has("sort") {
		filter.SortBy = infrastructure.CustomerSortBy(params.Get("sort"))
	}

	res, err := initializeGetAllCustomerHandler().Handle(ctx.Request().Context(), filter)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func createCustomer(ctx echo.Context) error {
	var req customer.ChangeOrCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	req.ID = 0
	res, err := initializeChangeOrCreateCustomerHandler().Handle(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, res)
}

func updateCustomer(ctx echo.Context) error {
	var req customer.ChangeOrCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return infrastructure.ErrCustomerNotFound
	}

	req.ID = uint(id)
	res, err := initializeChangeOrCreateCustomerHandler().Handle(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func getCustomerById(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return infrastructure.ErrCustomerNotFound
	}

	res, err := initializeGetCustomerByIdHandler().Handle(ctx.Request().Context(), uint(id))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func deleteCustomer(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return infrastructure.ErrCustomerNotFound
	}

	req := customer.DeleteRequest{Id: uint(id)}
	err = initializeDeleteCustomerHandler().Handle(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
