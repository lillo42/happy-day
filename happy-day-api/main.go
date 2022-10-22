package main

import (
	"happy_day/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	initConfiguration()

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))
	e.Use(middlewares.ErrorMiddleware)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		// AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))

	initializeReservationController().Routes(e)
	initializeProductController().Routes(e)
	initializeCustomerController().Routes(e)

	e.Logger.Fatal(e.Start(":5100"))
}

func initConfiguration() {
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