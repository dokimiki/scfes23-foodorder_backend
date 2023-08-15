package user_route

import (
	"errors"
	"net/http"
	"time"

	"github.com/dokimiki/scfes23-foodorder_backend/database"
	epr "github.com/dokimiki/scfes23-foodorder_backend/libs/errorPayloadResponse"
	"github.com/dokimiki/scfes23-foodorder_backend/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetOrder(c echo.Context) error {
	type res struct {
		OrderStatus      string `json:"orderStatus"`
		IsCanceled       bool   `json:"isCanceled"`
		IsPaid           bool   `json:"isPaid"`
		TimeOfCompletion *time.Time
	}

	jwtToken := c.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(jwt.MapClaims)
	token := claims["sub"].(string)
	storeId := c.Param("store_id")

	order := models.Order{}
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("deleted_at IS NULL").Where("user_token = ? AND store_str_id = ?", token, storeId).Last(&order).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			// create new order
			order = models.Order{
				UserToken:  token,
				StoreStrID: storeId,
			}
			if err := tx.Create(&order).Error; err != nil {
				return errors.New("オーダー作成に失敗しました。(failed to create new order)")
			}
		} else if err != nil {
			return errors.New("オーダーが見つかりません(failed to find order)")
		}

		return nil
	}); err != nil {
		return c.JSON(http.StatusOK, epr.APIError(err.Error()))
	}

	r := res{
		OrderStatus:      order.OrderStatus,
		IsCanceled:       order.IsCanceled,
		IsPaid:           order.IsPaid,
		TimeOfCompletion: order.TimeOfCompletion,
	}

	return c.JSON(http.StatusOK, r)
}

func GetOrderItems(c echo.Context) error {
	type res struct {
		OrderItems []struct {
			MenuID   uint32 `json:"-"`
			StrID    string `json:"id"`
			Quantity int    `json:"quantity"`
		} `json:"orderItems"`
	}

	jwtToken := c.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(jwt.MapClaims)
	token := claims["sub"].(string)
	storeId := c.Param("store_id")
	r := res{}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		order := models.Order{}
		if err := tx.Where("deleted_at IS NULL").Where("user_token = ? AND store_str_id = ?", token, storeId).Last(&order).Error; err != nil {
			return errors.New("オーダーが見つかりません(failed to find order)")
		}

		if err := tx.Model(models.OrderItem{}).
			Select("menus.str_id, order_items.quantity").
			Joins("INNER JOIN menus ON menus.id = order_items.menu_id").
			Where("order_items.order_id = ?", order.ID).
			Scan(&r.OrderItems).Error; err != nil {
			return errors.New("オーダーアイテムが見つかりません(failed to find order items)")
		}

		return nil
	}); err != nil {
		return c.JSON(http.StatusOK, epr.APIError(err.Error()))
	}

	return c.JSON(http.StatusOK, r)
}

func CancelOrder(c echo.Context) error {
	// storeId := c.Param("store_id")
	return nil
}
