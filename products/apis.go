package products

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"happyday/infra"
	"log/slog"
	"net/http"
	"strconv"
)

func Map(router *gin.RouterGroup) {
	r := router.Group("/products")

	r.DELETE("/:id", func(context *gin.Context) {
		logger := infra.ResolverLogger(context)

		value := context.Param("id")
		id, err := uuid.Parse(value)

		if err != nil {
			context.Status(http.StatusNoContent)
			return
		}

		command := CreateCommand(context)

		logger.InfoContext(context, "going to delete product", slog.String("id", id.String()))
		err = command.Delete(context, id)
		if err != nil {
			logger.WarnContext(context, "error during delete",
				slog.String("id", id.String()),
				slog.Any("err", err),
			)

			writeError(context, err)
			return
		}

		logger.InfoContext(context, "delete product with success", slog.String("id", id.String()))
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

		command := CreateCommand(context)
		logger.InfoContext(context, "going to create product")
		product, err := command.CreateOrChange(context, req)
		if err != nil {
			logger.WarnContext(context, "error during product creating", slog.Any("err", err))
			writeError(context, err)
			return
		}

		logger.InfoContext(context, "product created with success", slog.String("id", product.ID.String()))
		context.JSON(http.StatusCreated, &product)
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
		command := CreateCommand(context)

		logger.InfoContext(context, "going to change product", slog.String("id", id.String()))
		product, err := command.CreateOrChange(context, req)
		if err != nil {
			logger.WarnContext(context, "error during change product", slog.Any("err", err))
			writeError(context, err)
			return
		}

		logger.InfoContext(context, "product change with success", slog.String("id", id.String()))
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
			logger.WarnContext(context, "error during get product", slog.Any("err", err))
			writeError(context, err)
			return
		}

		if prod.Version == 0 {
			logger.InfoContext(context, "product not found", slog.String("id", id.String()))
			context.JSON(NotFound.Status, &NotFound)
			return
		}

		logger.InfoContext(context, "product get with success", slog.String("id", id.String()))
		context.JSON(http.StatusOK, &prod)
	})

	r.GET("", func(context *gin.Context) {
		logger := infra.ResolverLogger(context)

		val := context.Query("page")
		page, _ := strconv.ParseUint(val, 10, 64)
		if page > 0 {
			page = page - 1
		}

		val = context.Query("size")
		size, _ := strconv.ParseUint(val, 10, 64)
		if size == 0 {
			size = 50
		}

		filter := ProductFilter{
			Name: context.Query("name"),
			Page: int(page),
			Size: int(size),
		}

		repository := CreateRepository(context)
		logger.InfoContext(context, "going to get all product",
			slog.String("name", filter.Name),
			slog.Uint64("page", page),
			slog.Uint64("size", size),
		)
		res, err := repository.GetAll(context, filter)

		if err != nil {
			logger.WarnContext(context, "error during get all customer", slog.Any("err", err),
				slog.String("name", filter.Name),
				slog.Uint64("page", page),
				slog.Uint64("size", size),
			)
			writeError(context, err)
			return
		}

		logger.InfoContext(context, "get all product with success")
		context.JSON(http.StatusOK, res)
	})
}

func writeError(context *gin.Context, err error) {
	problem := mapErrorToProblemDetails(context, err)
	context.JSON(problem.Status, problem)
}

func mapErrorToProblemDetails(context *gin.Context, err error) infra.ProblemDetails {
	if errors.Is(err, ErrProductNotExists) {
		return NotFound
	}

	if errors.Is(err, ErrNameIsEmpty) {
		return NameIsEmpty
	}

	if errors.Is(err, ErrNameTooLarge) {
		return NameIsTooLarge
	}

	if errors.Is(err, ErrPriceIsInvalid) {
		return PriceIsInvalid
	}

	if errors.Is(err, ErrConcurrencyUpdate) {
		return ConcurrencyIssue
	}

	logger := infra.ResolverLogger(context)
	logger.ErrorContext(context, "unmapped exception", slog.Any("err", err))
	return InternalErrorServer
}

func CreateCommand(ctx context.Context) *Command {
	return &Command{
		repository: CreateRepository(ctx),
	}
}

func CreateRepository(ctx context.Context) ProductRepository {
	return &GormProductRepository{
		db: infra.GormFactory(ctx),
	}
}

var (
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
		Type:   "product-not-found",
		Title:  "Products not found",
	}

	ConcurrencyIssue = infra.ProblemDetails{
		Status: http.StatusConflict,
		Type:   "product-conflict",
		Title:  "Products update conflict",
	}

	NameIsEmpty = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "product-name-is-empty",
		Title:  "Products name is empty",
	}

	NameIsTooLarge = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "product-name-is-too-large",
		Title:  "Products name is too large",
	}

	PriceIsInvalid = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "product-price-is-invalid",
		Title:  "Products price is invalid",
	}
)
