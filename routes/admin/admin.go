package ar

import (
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
	var dbOrders []models.Order
	if err := database.DB.Where("order_status = ?", "ordered").Find(&dbOrders).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("注文情報の取得でエラーが発生しました。"))
	}

	// レスポンスを作成する
	response := []types.OrderedPotato{}
	for _, order := range dbOrders {
		// DBから注文した商品の数を数えて取得する
		var orderedItems []models.OrderItem
		if err := database.DB.Where("order_id = ?", order.ID).Find(&orderedItems).Error; err != nil {
			return c.JSON(http.StatusOK, epr.APIError("注文情報の取得でエラーが発生しました。"))
		}

		var qty int
		qty = 0

		for _, orderedItem := range orderedItems {
			if orderedItem.MenuID >= 16 {
				continue
			}
			qty += orderedItem.Quantity
		}

		response = append(response, types.OrderedPotato{
			ReceptionTime:  order.CreatedAt,
			CompletionTime: order.TimeOfCompletion,
			Qty:            qty,
			Order: struct {
				ID            string `json:"id"`
				IsMobileOrder bool   `json:"isMobileOrder"`
				IsPaid        bool   `json:"isPaid"`
				NumberTag     int    `json:"numberTag"`
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

func GetSeasoningData(c echo.Context) error {
	// 注文情報を取得
	dbOrders := []models.Order{}
	if err := database.DB.Where("order_status = ?", "cooked").Find(&dbOrders).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("注文情報の取得に失敗しました。"))
	}

	// レスポンスを作成
	resOrders := []types.Order{}
	for _, order := range dbOrders {
		// 注文情報を取得
		orderItems := []models.OrderItem{}
		if err := database.DB.Where("order_id = ?", order.ID).Find(&orderItems).Error; err != nil {
			return c.JSON(http.StatusOK, epr.APIError("注文情報の取得に失敗しました。"))
		}

		// cartItemの作成
		var cartItems []types.CartItem
		for _, orderItem := range orderItems {
			cartItems = append(cartItems, types.CartItem{
				ID:       strconv.FormatUint(uint64(orderItem.MenuID), 10),
				Quantity: orderItem.Quantity,
			})
		}

		// Order型に変換
		resOrders = append(resOrders, types.Order{
			ID:            strconv.FormatUint(uint64(order.ID), 10),
			IsMobileOrder: order.IsMobileOrder,
			NumberTag:     order.NumberTag,
			Items:         cartItems,
		})
	}

	// レスポンスを返却
	return c.JSON(http.StatusOK, resOrders)
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

func SendOrderData(c echo.Context) error {
	// cartデータをJSONから構造体に変換する
	orderData := types.OrderData{}
	if err := c.Bind(&orderData); err != nil {
		return c.JSON(http.StatusOK, epr.APIError("bodyが不正です。"))
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if orderData.OrderCode == "" {
		// time_of_completionをほかの注文状況から求める
		// カートの中の商品の数を数える
		var cartItemsCount int
		for _, cartItem := range orderData.Cart {
			cartItemsCount += cartItem.Quantity
		}
		// 注文情報を取得
		var latestCompletionTime time.Time
		latestOrder := models.Order{}
		if err := tx.Where("order_status = ?", "ordered").Order("created_at desc").First(&latestOrder).Error; err != nil {
			latestCompletionTime = time.Now()
		} else {
			latestCompletionTime = latestOrder.TimeOfCompletion
		}

		timeOfCompletion := latestCompletionTime
		if timeOfCompletion.Before(time.Now()) {
			timeOfCompletion = time.Now()
		}
		timeOfCompletion = timeOfCompletion.Add(time.Duration(cartItemsCount*2+2) * time.Minute)

		//新しいorderを作成する
		order := models.Order{
			UserID:           1,
			IsMobileOrder:    false,
			IsPaid:           true,
			NumberTag:        orderData.NumTag,
			OrderStatus:      "ordered",
			TimeOfCompletion: timeOfCompletion,
		}

		// order情報をDBに保存する
		if err := tx.Save(&order).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusOK, epr.APIError("order情報の保存に失敗しました。"))
		}

		// order_item情報をDBに保存する
		for _, cartItem := range orderData.Cart {
			menuId, _ := strconv.ParseUint(cartItem.ID, 10, 32)

			orderItem := models.OrderItem{
				OrderID:  order.ID,
				MenuID:   uint32(menuId),
				Quantity: cartItem.Quantity,
			}

			if err := tx.Save(&orderItem).Error; err != nil {
				tx.Rollback()
				return c.JSON(http.StatusOK, epr.APIError("order_item情報の保存に失敗しました。"))
			}
		}
	} else {
		// orderCodeからorder情報を取得する
		var barcode models.Barcode
		if err := tx.Where("barcode_data = ?", orderData.OrderCode).First(&barcode).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusOK, epr.APIError("注文情報の取得でエラーが発生しました。"))
		}

		// order情報を取得する
		var order models.Order
		if err := tx.Where("id = ?", barcode.OrderID).First(&order).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusOK, epr.APIError("注文情報の取得でエラーが発生しました。"))
		}

		// order情報を更新する
		order.NumberTag = orderData.NumTag
		order.IsPaid = true
		if err := tx.Save(&order).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusOK, epr.APIError("order情報の更新に失敗しました。"))
		}
		//TODO: 予約アイテム変更処理を入れる
	}

	// 注文完了
	tx.Commit()
	return c.JSON(http.StatusOK, true)
}

func FinishedFrying(c echo.Context) error {
	orderId := c.Param("orderId")

	// orderIdが一致しているorderの注文情報尾を更新
	if err := database.DB.Model(&models.Order{}).Where("id = ?", orderId).Update("order_status", "cooked").Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("注文情報の更新に失敗しました。"))
	}

	// レスポンスを返却
	return c.JSON(http.StatusOK, true)
}

func FinishedSeasoning(c echo.Context) error {
	orderId := c.Param("orderId")

	// orderIdが一致しているorderの注文情報尾を更新
	if err := database.DB.Model(&models.Order{}).Where("id = ?", orderId).Update("order_status", "received").Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("注文情報の更新に失敗しました。"))
	}

	// レスポンスを返却
	return c.JSON(http.StatusOK, true)
}
