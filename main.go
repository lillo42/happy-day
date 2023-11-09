package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lmittmann/tint"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"happyday/customers"
	"happyday/discounts"
	"happyday/infra"
	"happyday/orders"
	"happyday/products"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func main() {
	loadConfig()
	slog.SetDefault(createLogger())
	runDatabaseMigration()
	runHttpServer()
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func createLogger() *slog.Logger {
	mode := viper.GetString("logger.mode")
	if mode == "text" {
		return slog.New(tint.NewHandler(os.Stdout, nil))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func runDatabaseMigration() {
	infra.GormFactory = func(ctx context.Context) *gorm.DB {
		logger := infra.ResolverLogger(ctx)
		connectionString := viper.GetString("connectionString")
		db, err := gorm.Open(postgres.New(postgres.Config{DSN: connectionString}), &gorm.Config{
			Logger: &infra.SlogGorm{
				Logger:                    logger,
				LogLevel:                  gormlogger.Info,
				IgnoreRecordNotFoundError: true,
			},
		})

		if err != nil {
			logger.ErrorContext(ctx, "fatal error open database connection", slog.Any("err", err))
			panic(fmt.Errorf("fatal error open database connection: %w", err))
		}

		return db
	}

	db := infra.GormFactory(context.Background())

	err := db.AutoMigrate(&infra.Customer{},
		&infra.Product{},
		&infra.Discount{}, &infra.DiscountProducts{},
		&infra.Order{}, &infra.OrderPayment{}, &infra.OrderProduct{})

	if err != nil {
		slog.Error("fatal error run database migration", slog.Any("err", err))
		panic(fmt.Errorf("fatal error run database migration: %w", err))
	}
}

func runHttpServer() {
	gin.SetMode(viper.GetString("mode"))

	engine := gin.New()
	engine.Use(slogLogger)
	engine.Use(slogRecovery)
	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	apiRouter := engine.Group("/api")
	customers.Map(apiRouter)
	products.Map(apiRouter)
	discounts.Map(apiRouter)
	orders.Map(apiRouter)

	discounts.ProductServiceFactory = func(ctx context.Context) discounts.ProductService {
		return &GlobalProductService{
			repository: products.CreateRepository(ctx),
		}
	}

	orders.ProductServiceFactory = func(ctx context.Context) orders.ProductService {
		return &GlobalProductService{
			repository: products.CreateRepository(ctx),
		}
	}

	orders.CustomerServiceFactory = func(ctx context.Context) orders.CustomerService {
		return &GlobalCustomerService{
			repository: customers.CreateRepository(ctx),
		}
	}

	if err := engine.Run(); err != nil {
		slog.Error("fatal error to run HTTP Server", slog.Any("err", err))
		panic(fmt.Errorf("fatal error to run HTTP Server: %w", err))
	}
}

func slogLogger(context *gin.Context) {
	requestID := uuid.NewString()
	logger := createLogger()

	child := logger.With(
		slog.Group("request",
			slog.String("request_id", requestID),
			slog.String("method", context.Request.Method),
			slog.String("user_agent", context.Request.UserAgent()),
			slog.String("path", context.Request.URL.Path),
			slog.String("query_parameter", context.Request.URL.RawQuery),
			slog.String("client_ip", context.ClientIP()),
			slog.Any("header", context.Request.Header),
		),
	)

	child.Info("processing request")
	context.Set("logger", child)

	start := time.Now()
	context.Next()
	stop := time.Now()

	child.Info("request processed",
		slog.Int("status", context.Writer.Status()),
		slog.Int("body_size", context.Writer.Size()),
		slog.Any("response_headers", context.Writer.Header()),
		slog.Duration("time_taken_ms", stop.Sub(start)))
}

func slogRecovery(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// Check for a broken connection, as it is not really a
			// condition that warrants a panic stack trace.
			var brokenPipe bool
			if ne, ok := err.(*net.OpError); ok {
				if se, ok := ne.Err.(*os.SyscallError); ok {
					if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						brokenPipe = true
					}
				}
			}

			val, _ := context.Get("logger")
			logger := val.(*slog.Logger)
			httpRequest, _ := httputil.DumpRequest(context.Request, false)
			if brokenPipe {
				logger.ErrorContext(context, context.Request.URL.Path,
					slog.Any("error", err),
					slog.String("request_info", string(httpRequest)),
				)

				// If the connection is dead, we can't write a status to it.
				_ = context.Error(err.(error)) // nolint: errcheck
				context.Abort()
				return
			}

			logger.Error("Recovery from panic",
				slog.Any("error", err),
				slog.Group("recovery",
					slog.Time("time", time.Now()),
					slog.String("request_info", string(httpRequest)),
					slog.String("stack", string(debug.Stack())),
				),
			)

			defaultHandleRecovery(context)
		}
	}()

	context.Next()
}

func defaultHandleRecovery(c *gin.Context) {
	c.AbortWithStatus(http.StatusInternalServerError)
}
