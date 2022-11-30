package apis

import (
	"net/http"
	"strconv"

	"happy_day/application/product"
	"happy_day/infrastructure"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func MapProductEndpoints(e *echo.Echo) {
	e.GET("/api/v1/products", getAllProducts)
	e.POST("/api/v1/products", createProduct)

	e.GET("/api/v1/products/:id", getProductById)
	e.PUT("/api/v1/products/:id", updateProduct)
	e.DELETE("/api/v1/products/:id", deleteProduct)
}

func getAllProducts(ctx echo.Context) error {
	var filter infrastructure.ProductFilter
	filter.Text = ctx.QueryParam("text")
	filter.Size, _ = strconv.ParseInt(ctx.Param("size"), 10, 64)
	filter.Page, _ = strconv.ParseInt(ctx.Param("page"), 10, 64)
	filter.SortBy = infrastructure.ProductNameAsc

	params := ctx.QueryParams()
	if params.Has("sort") {
		filter.SortBy = infrastructure.ProductSortBy(params.Get("sort"))
	}

	res, err := initializeGetAllProductHandler().Handle(ctx.Request().Context(), filter)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func createProduct(ctx echo.Context) error {
	var req product.ChangeOrCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	req.Id = uuid.Nil
	res, err := initializeChangeOrCreateProductHandler().Handle(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, res)
}

func getProductById(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrProductNotFound
	}

	res, err := initalizeGetProductByIdHandler().Handle(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func updateProduct(ctx echo.Context) error {
	var req product.ChangeOrCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrInvalidBody
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrProductNotFound
	}

	req.Id = id
	res, err := initializeChangeOrCreateProductHandler().Handle(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func deleteProduct(ctx echo.Context) error {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return infrastructure.ErrProductNotFound
	}

	req := product.DeleteRequest{
		Id: id,
	}

	err = initializeDeleteProductHandler().Handle(ctx.Request().Context(), req)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}
