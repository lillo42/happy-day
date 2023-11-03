package customers

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
	r := router.Group("/customers")
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

		filter := CustomerFilter{
			Name:    context.Query("name"),
			Comment: context.Query("comment"),
			Phone:   context.Query("phones"),
			Page:    int(page),
			Size:    int(size),
		}

		repository := CreateRepository(context)

		logger.InfoContext(context, "going to get all customer",
			slog.String("name", filter.Name),
			slog.String("comment", filter.Comment),
			slog.String("phone", filter.Phone),
			slog.Uint64("page", page),
			slog.Uint64("size", size),
		)

		res, err := repository.GetAll(context, filter)

		if err != nil {
			logger.WarnContext(context, "error during get all customer", slog.Any("err", err),
				slog.String("name", filter.Name),
				slog.String("comment", filter.Comment),
				slog.String("phone", filter.Phone),
				slog.Uint64("page", page),
				slog.Uint64("size", size),
			)
			writeError(context, err)
			return
		}

		logger.InfoContext(context, "get all customer with success")
		context.JSON(http.StatusOK, res)
	})

	r.GET("/:id", func(context *gin.Context) {
		logger := infra.ResolverLogger(context)

		value := context.Param("id")
		id, err := uuid.Parse(value)

		if err != nil {
			context.JSON(NotFound.Status, &NotFound)
			return
		}

		repo := CreateRepository(context)

		logger.InfoContext(context, "going to get customer", slog.String("id", id.String()))
		customer, err := repo.GetOrCreate(context, id)
		if err != nil {
			logger.WarnContext(context, "error during get customer", slog.Any("err", err))
			writeError(context, err)
			return
		}

		if customer.Version == 0 {
			logger.InfoContext(context, "customer not found", slog.String("id", id.String()))
			context.JSON(NotFound.Status, &NotFound)
			return
		}

		logger.InfoContext(context, "customer get with success", slog.String("id", id.String()))
		context.JSON(http.StatusOK, customer)
	})

	r.POST("", func(context *gin.Context) {
		logger := infra.ResolverLogger(context)

		var req CreateOrChangeCustomer
		if err := context.BindJSON(&req); err != nil {
			logger.WarnContext(context, "error during json bind", slog.Any("err", err))
			context.JSON(BadRequest.Status, BadRequest)
			return
		}

		req.ID = uuid.New()
		command := createCommand(context)

		logger.InfoContext(context, "going to create customer")
		customer, err := command.CreateOrChange(context, req)
		if err != nil {
			logger.WarnContext(context, "error during customer creating", slog.Any("err", err))
			writeError(context, err)
			return
		}

		logger.InfoContext(context, "customer create with success", slog.String("id", customer.ID.String()))
		context.JSON(http.StatusCreated, &customer)
	})

	r.PUT("/:id", func(context *gin.Context) {
		logger := infra.ResolverLogger(context)

		value := context.Param("id")
		id, err := uuid.Parse(value)

		if err != nil {
			context.JSON(NotFound.Status, &NotFound)
			return
		}

		var req CreateOrChangeCustomer
		if err := context.BindJSON(&req); err != nil {
			logger.WarnContext(context, "error during json bind",
				slog.Any("err", err),
				slog.String("id", id.String()),
			)
			context.JSON(BadRequest.Status, BadRequest)
			return
		}

		req.ID = id
		command := createCommand(context)

		logger.InfoContext(context, "going to change customer", slog.String("id", id.String()))
		customer, err := command.CreateOrChange(context, req)
		if err != nil {
			logger.WarnContext(context, "error during change customer", slog.Any("err", err))
			writeError(context, err)
			return
		}

		logger.InfoContext(context, "customer change with success", slog.String("id", id.String()))
		context.JSON(http.StatusOK, &customer)
	})

	r.DELETE("/:id", func(context *gin.Context) {
		logger := infra.ResolverLogger(context)

		value := context.Param("id")
		id, err := uuid.Parse(value)

		if err != nil {
			context.Status(http.StatusNoContent)
			return
		}

		command := createCommand(context)

		logger.InfoContext(context, "going to delete customer", slog.String("id", id.String()))
		err = command.Delete(context, id)
		if err != nil {
			logger.WarnContext(context, "error during delete",
				slog.String("id", id.String()),
				slog.Any("err", err),
			)

			writeError(context, err)
			return
		}

		logger.InfoContext(context, "delete customer with success", slog.String("id", id.String()))
		context.Status(http.StatusNoContent)
	})
}

var (
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

func createCommand(ctx context.Context) *Command {
	return &Command{
		repository: CreateRepository(ctx),
	}
}

func CreateRepository(ctx context.Context) CustomerRepository {
	return &GormCustomerRepository{
		db: infra.GormFactory(ctx),
	}
}

func mapErrorToProblemDetails(context *gin.Context, err error) infra.ProblemDetails {
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

	logger := infra.ResolverLogger(context)
	logger.ErrorContext(context, "Unmapped exception", slog.Any("err", err))

	return InternalErrorServer
}

func writeError(context *gin.Context, err error) {
	problem := mapErrorToProblemDetails(context, err)
	context.JSON(problem.Status, problem)
}
