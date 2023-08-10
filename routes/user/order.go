package user_route

import (
	"net/http"

	"github.com/dokimiki/scfes23-foodorder_backend/database"
	epr "github.com/dokimiki/scfes23-foodorder_backend/libs/errorPayloadResponse"
	"github.com/dokimiki/scfes23-foodorder_backend/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetOrder(c echo.Context) error {
	jwtToken := c.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(jwt.MapClaims)
	token := claims["sub"].(string)

	user := models.User{}
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		tx.Where("token = ?", token).Take(&user)
		return nil
	}); err != nil {
		return c.JSON(http.StatusOK, epr.APIError(err.Error()))
	}

	return nil
}
