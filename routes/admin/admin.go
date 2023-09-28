package ar

import (
	"net/http"
	"strconv"

	epr "github.com/dokimiki/scfes23-foodorder_backend/libs/errorPayloadResponse"
	"github.com/dokimiki/scfes23-foodorder_backend/models"

	"github.com/dokimiki/scfes23-foodorder_backend/database"
	"github.com/dokimiki/scfes23-foodorder_backend/types"
	"github.com/labstack/echo/v4"
)

func GetPotatoData(c echo.Context) error {
	// DBから注文情報を取得する
	var orders []models.Order
	if err := database.DB.Where("order_status = ?", "ordered").Find(&orders).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("注文情報の取得でエラーが発生しました。"))
	}

	// レスポンスを作成する
	response := []types.OrderedPotato{}
	for _, order := range orders {
		// DBから注文した商品の数を数えて取得する
		var orderedItems []models.OrderItem
		if err := database.DB.Where("order_id = ?", order.ID).Find(&orderedItems).Error; err != nil {
			return c.JSON(http.StatusOK, epr.APIError("注文情報の取得でエラーが発生しました。"))
		}

		var qty int
		qty = 0

		for _, orderedItem := range orderedItems {
			qty += orderedItem.Quantity
		}

		response = append(response, types.OrderedPotato{
			ReceptionTime:  order.CreatedAt,
			CompletionTime: order.TimeOfCompletion,
			Qty:            qty,
			Order: struct {
				ID            string
				IsMobileOrder bool
				IsPaid        bool
				NumberTag     int
			}{
				ID:            strconv.FormatUint(uint64(order.ID), 10),
				IsMobileOrder: order.IsMobileOrder,
				IsPaid:        order.IsPaid,
				NumberTag:     order.NumberTag,
			},
		})
	}

	// JSONで返す
	return c.JSON(http.StatusOK, response)
}
