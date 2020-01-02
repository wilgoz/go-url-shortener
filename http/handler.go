package http

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"

	"github.com/wilgoz/go-url-shortener/shortener"
)

type RedirectHandler interface {
	Listen() error
}

type handler struct {
	redirectService shortener.RedirectService
	echo            *echo.Echo
}

func NewHandler(service shortener.RedirectService) RedirectHandler {
	h := &handler{
		redirectService: service,
		echo:            echo.New(),
	}
	h.setHandlers()
	return h
}

func (h *handler) setHandlers() {
	h.echo.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{
				echo.GET,
				echo.POST,
			},
		}),
	)
	h.echo.GET("/:code", h.get)
	h.echo.POST("/", h.post)
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

func (h *handler) Listen() error {
	return h.echo.Start(
		func() string {
			port := "8080"
			if os.Getenv("PORT") != "" {
				port = os.Getenv("PORT")
			}
			return fmt.Sprintf(":%s", port)
		}(),
	)
}
