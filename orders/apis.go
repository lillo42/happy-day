package orders

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"happyday/infra"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func Map(router *gin.RouterGroup) {
	r := router.Group("/orders")
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
		logger.InfoContext(context, "going to create order")
		order, err := command.CreateOrChange(context, req)
		if err != nil {
			logger.WarnContext(context, "error during order creating", slog.Any("err", err))
			writeError(context, err)
			return
		}

		logger.InfoContext(context, "order created with success", slog.String("id", order.ID.String()))
		context.JSON(http.StatusCreated, &order)
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

		logger.InfoContext(context, "going to change order", slog.String("id", id.String()))
		product, err := command.CreateOrChange(context, req)
		if err != nil {
			logger.WarnContext(context, "error during change order", slog.Any("err", err))
			writeError(context, err)
			return
		}

		logger.InfoContext(context, "order change with success", slog.String("id", id.String()))
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

		repo := createRepository(context)
		order, err := repo.GetOrCreate(context, id)
		if err != nil {
			logger.WarnContext(context, "error during get order", slog.Any("err", err))
			writeError(context, err)
			return
		}

		if order.Version == 0 {
			logger.InfoContext(context, "order not found", slog.String("id", id.String()))
			context.JSON(NotFound.Status, &NotFound)
			return
		}

		logger.InfoContext(context, "order get with success", slog.String("id", id.String()))
		context.JSON(http.StatusOK, &order)
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
		if size == 0 || size > 50 {
			size = 50
		}

		var deliveryBetween []time.Time
		at := context.Query("deliveryStartAt")
		date, err := time.Parse(time.RFC3339, at)
		if err == nil {
			deliveryBetween = append(deliveryBetween, date)
		}

		at = context.Query("deliveryEndAt")
		date, err = time.Parse(time.RFC3339, at)
		if err == nil {
			deliveryBetween = append(deliveryBetween, date)
		}

		filter := OrderFilter{
			Address:         context.Query("address"),
			Comment:         context.Query("comment"),
			CustomerName:    context.Query("customerName"),
			CustomerPhone:   context.Query("CustomerPhone"),
			DeliveryBetween: deliveryBetween,
			Page:            int(page),
			Size:            int(size),
		}

		repository := createRepository(context)
		logger.InfoContext(context, "going to get all orders",
			slog.String("address", filter.Address),
			slog.String("comment", filter.Comment),
			slog.String("customerName", filter.CustomerName),
			slog.String("customerPhone", filter.CustomerPhone),
			slog.Any("deliveryBetween", filter.DeliveryBetween),
			slog.Uint64("page", page),
			slog.Uint64("size", size),
		)
		res, err := repository.GetAll(context, filter)

		if err != nil {
			logger.WarnContext(context, "error during get all orders", slog.Any("err", err),
				slog.String("address", filter.Address),
				slog.String("comment", filter.Comment),
				slog.String("customerName", filter.CustomerName),
				slog.String("customerPhone", filter.CustomerPhone),
				slog.Any("deliveryBetween", filter.DeliveryBetween),
				slog.Uint64("page", page),
				slog.Uint64("size", size),
			)
			writeError(context, err)
			return
		}

		logger.InfoContext(context, "get all orders with success")
		context.JSON(http.StatusOK, res)
	})

}

func writeError(context *gin.Context, err error) {
	problem := mapErrorToProblemDetails(context, err)
	context.JSON(problem.Status, problem)
}

func mapErrorToProblemDetails(context *gin.Context, err error) infra.ProblemDetails {
	if errors.Is(err, ErrOrderNotFound) {
		return NotFound
	}

	if errors.Is(err, ErrAddressIsEmpty) {
		return AddressIsEmpty
	}

	if errors.Is(err, ErrAddressIsTooLarge) {
		return AddressIsTooLarge
	}

	if errors.Is(err, ErrDeliveryAtIsInvalid) {
		return DeliveryAtIsInvalid
	}

	if errors.Is(err, ErrTotalPriceIsInvalid) {
		return TotalPriceIsInvalid
	}

	if errors.Is(err, ErrDiscountIsInvalid) {
		return DiscountIsInvalid
	}

	if errors.Is(err, ErrFinalPriceIsInvalid) {
		return FinalPriceIsInvalid
	}

	if errors.Is(err, ErrCustomerNotFound) {
		return CustomerNotFound
	}

	if errors.Is(err, ErrProductNotFound) {
		return ProductNotFound
	}

	if errors.Is(err, ErrPaymentValueIsInvalid) {
		return PaymentValueIsInvalid
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
		productService:  ProductServiceFactory(ctx),
		customerService: CustomerServiceFactory(ctx),
		repository:      createRepository(ctx),
	}
}

func createRepository(ctx context.Context) OrderRepository {
	return &GormOrderRepository{
		db: infra.GormFactory(ctx),
	}
}

var (
	ProductServiceFactory  func(ctx context.Context) ProductService
	CustomerServiceFactory func(ctx context.Context) CustomerService

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
		Type:   "order-not-found",
		Title:  "Order not found",
	}

	ConcurrencyIssue = infra.ProblemDetails{
		Status: http.StatusConflict,
		Type:   "order-conflict",
		Title:  "Order update conflict",
	}

	AddressIsEmpty = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "order-address-is-empty",
		Title:  "Order address is empty",
	}

	AddressIsTooLarge = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "order-address-is-too-large",
		Title:  "Order address is too large",
	}

	DeliveryAtIsInvalid = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "order-delivery-at-is-invalid",
		Title:  "Order delivery at is invalid",
	}

	TotalPriceIsInvalid = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "order-total-price-at-is-invalid",
		Title:  "Order total price is invalid",
	}

	DiscountIsInvalid = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "order-discount-is-invalid",
		Title:  "Order discount is invalid",
	}

	FinalPriceIsInvalid = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "order-final-price-at-is-invalid",
		Title:  "Order final price is invalid",
	}

	CustomerNotFound = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "order-customer-not-found",
		Title:  "Order customer not found",
	}

	ProductNotFound = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "order-product-not-found",
		Title:  "Order product not found",
	}

	PaymentValueIsInvalid = infra.ProblemDetails{
		Status: http.StatusUnprocessableEntity,
		Type:   "order-payment-value-is-invalid",
		Title:  "Order payment value is invalid",
	}
)
