package apis

import (
	"happy_day/application"
	"happy_day/infrastructure"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ProductController struct {
	getAllHandler         application.GetAllProductsHandler
	getByIdHandler        application.GetProductByIdHandler
	createOrChangeHandler application.CreateOrChangeProductHandler
	deleteHandler         application.DeleteProductHandler
}

func (controller ProductController) Routes(e *echo.Echo) {
	e.GET("/api/v1/products", controller.getAll)
	e.POST("/api/v1/products", controller.create)

	e.GET("/api/v1/products/:id", controller.get)
	e.PUT("/api/v1/products/:id", controller.update)
	e.DELETE("/api/v1/products/:id", controller.delete)
}

func (controller ProductController) getAll(ctx echo.Context) error {
	var filter infrastructure.ProductFilter
	filter.Text = ctx.QueryParam("text")
	filter.Size, _ = strconv.ParseInt(ctx.Param("size"), 10, 64)
	filter.Page, _ = strconv.ParseInt(ctx.Param("page"), 10, 64)
	filter.SortBy = infrastructure.ProductNameAsc

	params := ctx.QueryParams()
	if params.Has("sort") {
		filter.SortBy = infrastructure.ProductSortBy(params.Get("sort"))
	}

	res, err := controller.getAllHandler.Handle(ctx.Request().Context(), filter)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func (controller ProductController) create(ctx echo.Context) error {
	var req application.CreateOrChangeProductRequest
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

func (controller ProductController) get(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrProductNotFound
	}

	res, err := controller.getByIdHandler.Handle(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func (controller ProductController) update(ctx echo.Context) error {
	var req application.CreateOrChangeProductRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrProductNotFound
	}

	req.Id = id
	res, err := controller.createOrChangeHandler.Handle(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func (controller ProductController) delete(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrProductNotFound
	}

	req := application.DeleteProductRequest{
		Id: id,
	}

	err = controller.deleteHandler.Handle(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
