package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var (
	levels = map[string]log.Lvl{
		"DEBUG": log.DEBUG,
		"INFO":  log.INFO,
		"WARN":  log.WARN,
		"ERROR": log.ERROR,
		"OFF":   log.OFF,
	}
)

func getLoggerLevel(level string) (log.Lvl, bool) {
	c, ok := levels[strings.ToUpper(level)]
	return c, ok
}

func main() {

	server := echo.New()

	loggerLevelEnv := os.Getenv("LOGGER_LEVEL")
	loggerLevel, ok := getLoggerLevel(loggerLevelEnv)

	if ok {
		server.Logger.SetLevel(loggerLevel)
	}

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	server.GET("/health", func(context echo.Context) error {
		return context.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	targetLabel := os.Getenv("TARGET_LABEL")
	if targetLabel == "" {
		server.Logger.Fatal("TARGET_LABEL not set")
		os.Exit(1)
	}
	server.Logger.Debug("TARGET_LABEL=" + targetLabel)

	valueLabel := os.Getenv("VALUE_LABEL")
	if valueLabel == "" {
		server.Logger.Fatal("VALUE_LABEL not set")
		os.Exit(1)
	}
	server.Logger.Debug("TARGET_LABEL=" + valueLabel)

	server.GET("/", func(context echo.Context) error {
		target := context.Request().Header.Get(targetLabel)
		if target == "" {
			server.Logger.Debug(targetLabel + " not found")
			return context.NoContent(http.StatusUnauthorized)
		}

		value := context.Request().Header.Get(valueLabel)
		if value == "" {
			server.Logger.Debug(valueLabel + " not found")
			return context.NoContent(http.StatusUnauthorized)
		}

		if target == value {
			return context.NoContent(http.StatusOK)
		} else {
			return context.NoContent(http.StatusUnauthorized)
		}
	})

	server.Logger.Fatal(server.Start(":" + httpPort))
}
