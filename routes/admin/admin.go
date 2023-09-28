package ar

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

func GetCartDataFromOrderCode(c echo.Context) error {
	// orderCodeからorder情報を取得する
	var barcode models.Barcode
	if err := database.DB.Where("barcode_data = ?", c.Param("orderCode")).First(&barcode).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("注文情報の取得でエラーが発生しました。"))
	}

	// order情報からorder_item情報を取得する
	var orderItems []models.OrderItem
	if err := database.DB.Where("order_id = ?", barcode.OrderID).Find(&orderItems).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("注文情報の取得でエラーが発生しました。"))
	}

	// レスポンスを作成する
	response := []types.CartItem{}
	for _, orderItem := range orderItems {
		response = append(response, types.CartItem{
			ID:       strconv.FormatUint(uint64(orderItem.MenuID), 10),
			Quantity: orderItem.Quantity,
		})
	}

	// JSONで返す
	return c.JSON(http.StatusOK, response)
}

func sendOrderData(c echo.Context) error {
	// cartデータをJSONから構造体に変換する

	cart := c.QueryParam("cart")
	// カートをJSON形式に変換
	cartItems := []types.CartItem{}

	err := json.Unmarshal([]byte(cart), &cartItems)
	if err != nil {
		return c.JSON(http.StatusBadRequest, epr.APIError("カートのJSON形式が不正です。"))
	}

	// order情報を作成する
	var order models.Order
	order.BarcodeData = c.Param("orderCode")
	order.ReceptionTime = time.Now()
	order.CompletionTime = time.Now()
	order.Qty = 0

	// cart情報をorder情報に追加する
	for _, cartItem := range cart {
		order.Qty += cartItem.Quantity
		order.Items = append(order.Items, models.OrderItem{
			MenuID:   cartItem.ID,
			Quantity: cartItem.Quantity,
		})
	}

	// order情報をDBに保存する
	if err := database.DB.Save(&order).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, epr.APIError("order情報の保存に失敗しました。"))
	}

	// 注文完了
	return c.JSON(http.StatusOK, true)
}
