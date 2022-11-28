package main

import (
	"happy_day/apis"
	"happy_day/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	initConfiguration()
	e := echo.New()

	if isDebug() {
		e.Debug = true
	}

	e.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))
	e.Use(middleware.Recover())

	e.Use(middlewares.ErrorMiddleware)

	e.Use(middleware.CORSWithConfig(getCorsConfig()))

	apis.MapCustomerEndpoints(e)
	apis.MapProductEndpoints(e)
	apis.MapReservationEndpoints(e)

	e.Logger.Info(e.Start(":5100"))
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

func getCorsConfig() middleware.CORSConfig {
	var config middleware.CORSConfig
	config.AllowOrigins = viper.GetStringSlice("cors.allow_origins")
	config.AllowHeaders = viper.GetStringSlice("cors.allow_headers")
	config.AllowMethods = viper.GetStringSlice("cors.allow_methods")
	return config
}

func isDebug() bool {
	return viper.GetBool("is_debug")
}
