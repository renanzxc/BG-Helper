package httphtml

import (
	"bghelper/internal/config"
	"bghelper/pkg/utils/cache"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rotisserie/eris"
)

type HTTPHTML struct {
	echo              *echo.Echo
	basePathTemplates string

	cache    cache.Cache
	validate *validator.Validate
	cfg      config.Config
}

var (
	//go:embed templates/*.html static/*
	content embed.FS
)

func (h *HTTPHTML) Setup() (err error) {
	h.cfg = config.GetConfig()
	h.cache, err = cache.NewJSONCache("./")
	if err != nil {
		return eris.Wrap(err, "")
	}
	h.validate = validator.New(validator.WithRequiredStructEnabled())

	h.echo = echo.New()
	h.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out},"data:"${custom}"}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		CustomTagFunc: func(c echo.Context, buf *bytes.Buffer) (int, error) {
			var errDataStr string

			switch errData := c.Get("err").(type) {
			case string:
				if errData != "" {
					errDataStr = `"` + errData + `"`
				}
			case map[string]interface{}:
				if errData != nil {
					var errDataBytes []byte
					if errDataBytes, err = json.Marshal(errData); err != nil {
						return 0, err
					}
					errDataStr = string(errDataBytes)
				}
			}

			if errDataStr != "" {
				return buf.WriteString(`{"err":` + errDataStr + `}`)
			}

			return 0, nil
		},
	}))
	h.echo.Use(middleware.Recover())
	h.echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("html", h)
			return next(c)
		}
	})
	h.echo.HTTPErrorHandler = func(err error, c echo.Context) {
		errData := eris.ToJSON(err, true)

		c.Set("err", errData)
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	h.echo.GET("/", func(c echo.Context) error {
		return yourUsername(c, h)
	})
	h.echo.GET("/playnext", func(c echo.Context) error {
		return whoPlayNext(c, h)
	})
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
	cfg := config.GetConfig()
	h.echo.Logger.Fatal(h.echo.Start(fmt.Sprintf(":%d", cfg.MyPort)))
}

func (h *HTTPHTML) Shutdown() {
	sig := make(chan os.Signal, 1)

	// Notify the sig channel of the following signals: SIGINT SIGTERM
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	fmt.Println("Waiting for a signal...")

	<-sig

	h.cache.Down()
}
