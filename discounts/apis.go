package discounts

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"mec/infra"
	"net/http"
	"strconv"
)

func Map(router *gin.RouterGroup) {
	r := router.Group("/discounts")

	r.DELETE("/:id", func(context *gin.Context) {
		logger := infra.ResolverLogger(context)

		value := context.Param("id")
		id, err := uuid.Parse(value)

		if err != nil {
			context.Status(http.StatusNoContent)
			return
		}

		command := createCommand(context)
		logger.InfoContext(context, "going to delete discount", slog.String("id", id.String()))
		if err = command.Delete(context, id); err != nil {
			logger.WarnContext(context, "error during delete",
				slog.String("id", id.String()),
				slog.Any("err", err),
			)

			writeError(context, err)
			return
		}

		logger.InfoContext(context, "delete discount with success", slog.String("id", id.String()))
		context.Status(http.StatusNoContent)
	})

	r.POST("", func(context *gin.Context) {
		logger := infra.ResolverLogger(context)

		var req CreateOrChange
		if err := context.BindJSON(&req); err != nil {
			logger.WarnContext(context, "error during json bind", slog.Any("err", err))
			context.JSON(BadRequest.Status, &BadRequest)
			return
		}

		command := createCommand(context)
		logger.InfoContext(context, "going to create discount")
		discount, err := command.CreateOrChange(context, req)
		if err != nil {
			logger.WarnContext(context, "error during discount creating", slog.Any("err", err))
			writeError(context, err)
			return
		}

		logger.InfoContext(context, "discount created with success", slog.String("id", discount.ID.String()))
		context.JSON(http.StatusCreated, &discount)

	})

	r.PUT("/:id", func(context *gin.Context) {
		logger := infra.ResolverLogger(context)

		val := context.Param("id")
		id, err := uuid.Parse(val)

		if err != nil {
			context.JSON(NotFound.Status, &NotFound)
			return
		}

		var req CreateOrChange
		if err := context.BindJSON(&req); err != nil {
			logger.WarnContext(context, "error during json bind",
				slog.Any("err", err),
				slog.String("id", id.String()),
			)
			context.JSON(BadRequest.Status, &BadRequest)
			return
		}

		req.ID = id
		command := createCommand(context)

		logger.InfoContext(context, "going to change discount", slog.String("id", id.String()))
		product, err := command.CreateOrChange(context, req)
		if err != nil {
			logger.WarnContext(context, "error during change discount", slog.Any("err", err))
			writeError(context, err)
			return
		}

		logger.InfoContext(context, "discount change with success", slog.String("id", id.String()))
		context.JSON(http.StatusOK, &product)
	})

	r.GET("/:id", func(context *gin.Context) {
		logger := infra.ResolverLogger(context)

		val := context.Param("id")
		id, err := uuid.Parse(val)

		if err != nil {
			context.JSON(NotFound.Status, &NotFound)
			return
		}

		repo := CreateRepository(context)
		prod, err := repo.GetOrCreate(context, id)
		if err != nil {
			logger.WarnContext(context, "error during get discount", slog.Any("err", err))
			writeError(context, err)
			return
		}

		if prod.Version == 0 {
			logger.InfoContext(context, "discount not found", slog.String("id", id.String()))
			context.JSON(NotFound.Status, &NotFound)
			return
		}

		logger.InfoContext(context, "discount get with success", slog.String("id", id.String()))
		context.JSON(http.StatusOK, &prod)
	})

	r.GET("", func(context *gin.Context) {
		logger := infra.ResolverLogger(context)

		val := context.Query("page")
		page, _ := strconv.ParseUint(val, 10, 64)
		if page < 0 {
			page = 0
		}

		val = context.Query("size")
		size, _ := strconv.ParseUint(val, 10, 64)
		if size == 0 {
			size = 50
		}

		filter := DiscountFilter{
			Name: context.Query("name"),
			Page: int(page),
			Size: int(size),
		}

		repo := CreateRepository(context)

		logger.InfoContext(context, "going to get all discounts",
			slog.String("name", filter.Name),
			slog.Uint64("page", page),
			slog.Uint64("size", size),
		)

		res, err := repo.GetAll(context, filter)
		if err != nil {
			logger.WarnContext(context, "error during get all discounts", slog.Any("err", err),
				slog.String("name", filter.Name),
				slog.Uint64("page", page),
				slog.Uint64("size", size),
			)
			writeError(context, err)
			return
		}

		logger.InfoContext(context, "get all discount with success")
		context.JSON(http.StatusOK, res)
	})
}

func writeError(context *gin.Context, err error) {
	problem := mapErrorToProblemDetails(context, err)
	context.JSON(problem.Status, problem)
}

func mapErrorToProblemDetails(context *gin.Context, err error) infra.ProblemDetails {
	if errors.Is(err, ErrDiscountNotFound) {
		return NotFound
	}

	if errors.Is(err, ErrNameIsEmpty) {
		return NameIsEmpty
	}

	if errors.Is(err, ErrNameIsTooLarge) {
		return NameIsTooLarge
	}

	if errors.Is(err, ErrPriceIsInvalid) {
		return PriceIsInvalid
	}

	if errors.Is(err, ErrProductsIsMissing) {
		return ProductsIsMissing
	}

	if errors.Is(err, ErrProductNotFound) {
		return ProductNotFound
	}

	if errors.Is(err, ErrConcurrencyUpdate) {
		return ConcurrencyIssue
	}

	logger := infra.ResolverLogger(context)
	logger.ErrorContext(context, "unmapped exception", slog.Any("err", err))
	return InternalErrorServer
}

func createCommand(ctx context.Context) *Command {
	return &Command{
		productService: ProductServiceFactory(ctx),
		repository:     CreateRepository(ctx),
	}
}

func CreateRepository(ctx context.Context) DiscountRepository {
	return &GormDiscountRepository{
		db: infra.GormFactory(ctx),
	}
}

var (
	ProductServiceFactory func(ctx context.Context) ProductService

	InternalErrorServer = infra.ProblemDetails{
		Status: http.StatusInternalServerError,
		Type:   "internal-error-server",
		Title:  "Internal error server",
	}

	BadRequest = infra.ProblemDetails{
		Status: http.StatusBadRequest,
		Type:   "invalid-payload",
		Title:  "Invalid Payload",
	}

	NotFound = infra.ProblemDetails{
		Status: http.StatusNotFound,
		Type:   "discount-not-found",
		Title:  "Discount not found",
	}

	ConcurrencyIssue = infra.ProblemDetails{
		Status: http.StatusConflict,
		Type:   "discount-conflict",
		Title:  "Discount update conflict",
	}

	NameIsEmpty = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "discount-name-is-empty",
		Title:  "Discount name is empty",
	}

	NameIsTooLarge = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "discount-name-is-too-large",
		Title:  "Discount name is too large",
	}

	PriceIsInvalid = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "discount-price-is-invalid",
		Title:  "Discount price is invalid",
	}

	ProductsIsMissing = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "discount-products-is-missing",
		Title:  "Discount products list is missing",
	}

	ProductNotFound = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "discount-product-not-found",
		Title:  "Discount product not found",
	}
)
