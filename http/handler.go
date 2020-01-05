package http

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"

	"github.com/wilgoz/go-url-shortener/shortener"
)

// RedirectHandler provides a utility http interface to listen and serve requests
type RedirectHandler interface {
	Listen(port string) error
	Shutdown()
}

type handler struct {
	redirectService shortener.RedirectService
	echo            *echo.Echo
}

func (h *handler) setHandlers() {
	h.echo.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowMethods: []string{
				echo.GET,
				echo.POST,
			},
		}),
	)
	h.echo.GET("/:code", h.get)
	h.echo.POST("/", h.post)
}

// NewHandler initializes server configs and routes and returns the handler
func NewHandler(service shortener.RedirectService) RedirectHandler {
	h := &handler{
		redirectService: service,
		echo:            echo.New(),
	}
	h.setHandlers()
	return h
}

func (h *handler) get(c echo.Context) error {
	redirect, err := h.redirectService.Find(c.Param("code"))
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectNotFound {
			return &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: http.StatusText(http.StatusNotFound),
			}
		}
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}
	return c.Redirect(
		http.StatusMovedPermanently, redirect.Original,
	)
}

func (h *handler) post(c echo.Context) error {
	redirect := &shortener.Redirect{}
	err := c.Bind(redirect)
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusUnprocessableEntity,
			Message: http.StatusText(http.StatusUnprocessableEntity),
		}
	}
	err = h.redirectService.Store(redirect)
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectInvalid {
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: http.StatusText(http.StatusBadRequest),
			}
		}
	}
	return c.JSON(http.StatusCreated, redirect)
}

func (h *handler) Listen(port string) error {
	return h.echo.Start(port)
}

func (h *handler) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.echo.Shutdown(ctx); err != nil {
		h.echo.Logger.Fatal(err)
	}
}
