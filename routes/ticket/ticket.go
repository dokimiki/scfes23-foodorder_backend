package ticket_route

import (
	"net/http"

	"github.com/dokimiki/scfes23-foodorder_backend/database"
	"github.com/dokimiki/scfes23-foodorder_backend/models"
	"github.com/labstack/echo/v4"
)

type price struct {
	YenPricePerTicket int `json:"yenPrice"`
}

func GetPrice(c echo.Context) error {
	result := price{}

	database.DB.Model(models.Ticket{}).Limit(1).Scan(&result)

	return c.JSON(http.StatusOK, result)
}