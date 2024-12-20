package httphtml

import (
	"bghelper/pkg/utils/cache"
	"embed"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rotisserie/eris"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type HTTPHTML struct {
	echo              *echo.Echo
	basePathTemplates string

	cache    cache.Cache
	validate *validator.Validate
}

var (
	//go:embed templates/*.html static/*
	content embed.FS
)

func (h *HTTPHTML) Setup() (err error) {
	h.cache, err = cache.NewJSONCache("./")
	if err != nil {
		return eris.Wrap(err, "")
	}
	h.validate = validator.New(validator.WithRequiredStructEnabled())

	h.echo = echo.New()
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
		c.Logger().Error(eris.ToString(err, true))

		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	h.echo.GET("/playnext/:username/with/:another_player", func(c echo.Context) error {
		return playnext(c, h)
	})
	h.echo.GET("/playnext/:username/with/:another_player/pdf", func(c echo.Context) error {
		return playnextPDF(c, h)
	})

	staticFS, err := fs.Sub(content, "static")
	if err != nil {
		return eris.Wrap(err, "")
	}

	staticHandler := http.FileServer(http.FS(staticFS))
	h.echo.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", staticHandler)))

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
