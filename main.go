package main

import (
	"time"

	"happyday/middlewares"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	initializeConfiguration()
	logger := initializeLogger()
	middlewares.Logger = logger

	engine := gin.New()

	if isDebug() {
		gin.ForceConsoleColor()
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	engine.Use(ginzap.RecoveryWithZap(logger, true))

	customerController := initializeCustomerController()
	customerController.MapEndpoints(engine)
	middlewares.AddErrors(customerController.ErrorMapping())

	productController := initializeProductController()
	productController.MapEndpoint(engine)
	middlewares.AddErrors(productController.ErrorMapping())

	engine.Run()
}

func isDebug() bool {
	return viper.GetBool("isDebug")
}

func initializeConfiguration() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("HAPPY_DAY_")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initializeLogger() *zap.Logger {
	fields := zap.Fields(
		zap.String("Application", "Happy Day"),
		zap.String("Version", Version),
	)

	var logger *zap.Logger
	var err error

	if isDebug() {
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger = zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(colorable.NewColorableStdout()),
			zapcore.DebugLevel)).
			WithOptions(fields)
	} else {
		logger, err = zap.NewProduction(fields)
	}

	if err != nil {
		panic(err)
	}

	return logger
}

var (
	Version = "0.0.0"
)
