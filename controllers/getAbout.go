package controllers

import (
	"net/http"
	"urlier/templates"

	"github.com/labstack/echo/v4"
)

func GetAbout(c echo.Context) error {
	return renderPage(c, http.StatusOK, templates.About(), "About")
}
