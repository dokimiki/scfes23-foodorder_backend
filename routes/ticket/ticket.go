package ticket_route

import (
	"net/http"

	"github.com/dokimiki/scfes23-foodorder_backend/database"
	"github.com/dokimiki/scfes23-foodorder_backend/models"
	"github.com/labstack/echo/v4"
)

func GetPrice(c echo.Context) error {
	type res struct {
		YenPricePerTicket int `json:"yenPrice"`
	}
	r := res{}

	database.DB.Model(models.Ticket{}).Limit(1).Scan(&r)

	return c.JSON(http.StatusOK, r)
}
