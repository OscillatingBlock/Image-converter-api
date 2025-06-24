package main

import (
	"log/slog"
	"os"
	"strconv"

	"imageConverter/internal/api"
	"imageConverter/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config := config.GetConfig()
	e := echo.New()

	e.Use(middleware.BodyLimit(config.MaxImageSize + "M"))

	portStr := config.Port
	if _, err := strconv.Atoi(portStr); err != nil {
		slog.Error("Invalid port number", "error", err)
		os.Exit(1)
	}

	e.POST("/convert", api.Convert)
	e.POST("/square-crop", api.SquareCropHandler)
	e.POST("/fit-to-square", api.FitToSquareHandler)
	e.POST("/invert", api.InvertHandler)
	e.POST("/apply-filter", api.ApplyFilterHandler)

	slog.Info("Starting server", "port", portStr)
	e.Logger.Fatal(e.Start(":" + portStr))
}
