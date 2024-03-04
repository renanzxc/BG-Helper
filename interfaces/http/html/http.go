package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/renanzxc/BG-Helper/utils/cache"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
)

type HTTPHTML struct {
	echo              *echo.Echo
	basePathTemplates string

	cache cache.Cache
}

func (h *HTTPHTML) Setup() (err error) {
	h.echo = echo.New()
	h.cache, err = cache.NewJSONCache("./")
	if err != nil {
		return
	}

	h.echo.Use(middleware.Logger())
	h.echo.Use(middleware.Recover())
	h.echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("html", h)
			return next(c)
		}
	})
	// TODO: Refactor
	h.echo.HTTPErrorHandler = func(err error, c echo.Context) {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	h.echo.GET("/playnext/:username/with/:another_player", func(c echo.Context) error {
		return playnext(c, h)
	})

	actualPath, err := os.Getwd()
	if err != nil {
		return
	}

	h.basePathTemplates = path.Join(actualPath, "/template")

	return
}

func (h *HTTPHTML) Run() {
	// Start server
	h.echo.Logger.Fatal(h.echo.Start(":1323"))
}

func (h *HTTPHTML) Shutdown() {
	sig := make(chan os.Signal, 1)

	// Notify the sig channel of the following signals: SIGINT, SIGTERM, and SIGKILL
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// Block until a signal is received
	fmt.Println("Waiting for a signal...")

	_ = <-sig

	h.cache.Down()

	return
}
