package customers

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

	db.AutoMigrate(&infra.Customer{})

	r := router.Group("/customers")
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

		filter := CustomerFilter{
			Name:    context.Query("name"),
			Comment: context.Query("comment"),
			Phone:   context.Query("phones"),
			Page:    int(page),
			Size:    int(size),
		}

		repository := createRepository()
		res, err := repository.GetAll(context, filter)

		if err != nil {
			writeError(context, err)
		} else {
			context.JSON(http.StatusOK, res)
		}
	})

	r.GET("/:id", func(context *gin.Context) {
		value := context.Param("id")
		id, err := uuid.Parse(value)

		if err != nil {
			context.JSON(NotFound.Status, &NotFound)
			return
		}

		repo := createRepository()
		customer, err := repo.GetOrCreate(context, id)
		if err != nil {
			writeError(context, err)
			return
		}

		if customer.Version == 0 {
			context.JSON(NotFound.Status, &NotFound)
			return
		}

		context.JSON(http.StatusOK, customer)
	})

	r.POST("", func(context *gin.Context) {
		var req CreateOrChangeCustomer
		if err := context.BindJSON(&req); err != nil {
			context.JSON(BadRequest.Status, BadRequest)
			return
		}

		req.ID = uuid.New()

		command := createCommand()
		customer, err := command.CreateOrChange(context, req)
		if err != nil {
			writeError(context, err)
		} else {
			context.JSON(http.StatusCreated, &customer)
		}
	})

	r.PUT("/:id", func(context *gin.Context) {
		value := context.Param("id")
		id, err := uuid.Parse(value)

		if err != nil {
			context.JSON(NotFound.Status, &NotFound)
			return
		}

		var req CreateOrChangeCustomer
		if err := context.BindJSON(&req); err != nil {
			context.JSON(BadRequest.Status, BadRequest)
			return
		}

		req.ID = id
		command := createCommand()
		customer, err := command.CreateOrChange(context, req)
		if err != nil {
			writeError(context, err)
		} else {
			context.JSON(http.StatusOK, &customer)
		}
	})

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
}

var (
	GormFactory func() *gorm.DB

	Logger *zap.Logger

	NotFound = infra.ProblemDetails{
		Status: http.StatusNotFound,
		Type:   "customer-not-found",
		Title:  "Customer not found",
	}

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

	ConcurrencyIssue = infra.ProblemDetails{
		Status: http.StatusConflict,
		Type:   "customer-conflict",
		Title:  "Customer update conflict",
	}

	NameIsEmpty = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "customer-name-is-empty",
		Title:  "Customer name is empty",
	}

	NameIsTooLarge = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "customer-name-is-too-large",
		Title:  "Customer name is too large",
	}

	PixIsTooLarge = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "customer-pix-is-too-large",
		Title:  "Customer pix is too large",
	}

	PhoneNumberIsInvalid = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "customer-phone-number-is-invalid",
		Title:  "Customer phone number is invalid",
	}
)

func createCommand() *Command {
	return &Command{
		repository: &GormCustomerRepository{
			db: GormFactory(),
		},
	}
}

func createRepository() CustomerRepository {
	return &GormCustomerRepository{
		db: GormFactory(),
	}
}

func mapErrorToProblemDetails(err error) infra.ProblemDetails {
	if errors.Is(err, ErrNotFound) {
		return NotFound
	}

	if errors.Is(err, ErrConcurrencyUpdate) {
		return ConcurrencyIssue
	}

	if errors.Is(err, ErrNameIsEmpty) {
		return NameIsEmpty
	}

	if errors.Is(err, ErrNameIsTooLarge) {
		return NameIsTooLarge
	}

	if errors.Is(err, ErrPixIsTooLarge) {
		return PixIsTooLarge
	}

	if errors.Is(err, ErrPhoneNumberIsInvalid) {
		return PhoneNumberIsInvalid
	}

	Logger.Error("Unmapped exception", zap.Error(err))

	return InternalErrorServer
}

func writeError(context *gin.Context, err error) {
	problem := mapErrorToProblemDetails(err)
	context.JSON(problem.Status, problem)
}
