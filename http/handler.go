package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/wilgoz/go-url-shortener/shortener"
)

type RedirectHandler interface {
	Get(c echo.Context) error
	Post(c echo.Context) error
}

type handler struct {
	redirectService shortener.RedirectService
}

func NewHandler(service shortener.RedirectService) RedirectHandler {
	return &handler{service}
}

func (h *handler) Get(c echo.Context) error {
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
	return c.Redirect(http.StatusMovedPermanently, redirect.Original)
}

func (h *handler) Post(c echo.Context) error {
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
