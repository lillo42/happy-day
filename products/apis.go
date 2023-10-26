package products

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"happyday/infra"
	"net/http"
	"strconv"
)

func Map(router *gin.RouterGroup) {
	db := GormFactory()

	db.AutoMigrate(&infra.Product{}, &infra.BoxProduct{})
	r := router.Group("/products")

	r.DELETE("/:id", func(context *gin.Context) {
		value := context.Param("id")
		id, err := uuid.Parse(value)

		if err != nil {
			context.Status(http.StatusNoContent)
			return
		}

		command := createCommand()
		err = command.Delete(context, id)
		if err != nil {
			writeError(context, err)
		} else {
			context.Status(http.StatusNoContent)
		}
	})

	r.POST("", func(context *gin.Context) {
		var req CreateOrChange
		if err := context.BindJSON(&req); err != nil {
			context.JSON(BadRequest.Status, &BadRequest)
			return
		}

		command := createCommand()
		product, err := command.CreateOrChange(context, req)
		if err != nil {
			writeError(context, err)
		} else {
			context.JSON(http.StatusCreated, &product)
		}
	})

	r.PUT("/:id", func(context *gin.Context) {
		val := context.Param("id")
		id, err := uuid.Parse(val)

		if err != nil {
			context.JSON(NotFound.Status, &NotFound)
			return
		}

		var req CreateOrChange
		if err := context.BindJSON(&req); err != nil {
			context.JSON(BadRequest.Status, &BadRequest)
			return
		}

		req.ID = id
		command := createCommand()
		product, err := command.CreateOrChange(context, req)
		if err != nil {
			writeError(context, err)
		} else {
			context.JSON(http.StatusOK, &product)
		}
	})

	r.GET("/:id", func(context *gin.Context) {
		val := context.Param("id")
		id, err := uuid.Parse(val)

		if err != nil {
			context.JSON(NotFound.Status, &NotFound)
			return
		}

		repo := createRepository()
		prod, err := repo.GetOrCreate(context, id)
		if err != nil {
			writeError(context, err)
		} else if prod.Version == 0 {
			context.JSON(NotFound.Status, &NotFound)
		} else {
			context.JSON(http.StatusOK, &prod)
		}
	})

	r.GET("", func(context *gin.Context) {
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

		repository := createRepository()
		res, err := repository.GetAll(context, filter)

		if err != nil {
			writeError(context, err)
		} else {
			context.JSON(http.StatusOK, res)
		}
	})
}

func writeError(context *gin.Context, err error) {
	problem := mapErrorToProblemDetails(err)
	context.JSON(problem.Status, problem)
}

func mapErrorToProblemDetails(err error) infra.ProblemDetails {
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

	if errors.Is(err, ErrBoxProductNotExists) {
		return BoxProductNotExists
	}

	Logger.Error("unmapped exception", zap.Error(err))
	return InternalErrorServer
}

func createCommand() *Command {
	return &Command{
		repository: &GormProductRepository{
			db: GormFactory(),
		},
	}
}

func createRepository() ProductRepository {
	return &GormProductRepository{
		db: GormFactory(),
	}
}

var (
	GormFactory func() *gorm.DB

	Logger *zap.Logger

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

	BoxProductNotExists = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "product-box-product-not-exists",
		Title:  "Products box product not exists",
	}
)
