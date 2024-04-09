package controllers

import (
	"fmt"
	"net/http"
	"urlier/templates"

	"github.com/labstack/echo/v4"
)

func GetEmpty(c echo.Context) error {
	fmt.Println("GET /empty")
	return renderComponent(c, http.StatusOK, templates.Empty())
}
