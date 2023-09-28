package ur

import (
	"encoding/json"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dokimiki/scfes23-foodorder_backend/database"
	epr "github.com/dokimiki/scfes23-foodorder_backend/libs/errorPayloadResponse"
	gt "github.com/dokimiki/scfes23-foodorder_backend/libs/generateToken"
	"github.com/dokimiki/scfes23-foodorder_backend/models"
	"github.com/dokimiki/scfes23-foodorder_backend/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func issueUserJWT(token string) string {
	byteJWT := jwt.New(jwt.SigningMethodHS256)
	claims := byteJWT.Claims.(jwt.MapClaims)
	claims["sub"] = token
	claims["exp"] = time.Now().Add(time.Hour * 24 * 15).Unix()

	secret := os.Getenv("SCFES23FOODORDER_JWT_SIGNATURE")
	strJWT, err := byteJWT.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}

	return strJWT
}

func SignUp(c echo.Context) error {
	// ユーザーIDを生成
	token := gt.GenUserToken()
	userID := issueUserJWT(token)

	// ユーザー情報を保存
	user := models.User{
		Token: token,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("ユーザー登録に失敗しました。"))
	}

	// ユーザーIDを返す
	response := types.User{
		ID: userID,
	}
	return c.JSON(http.StatusOK, response)
}

func SignIn(c echo.Context) error {
	// ユーザーIDを取得
	jwtToken := c.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(jwt.MapClaims)
	token := claims["sub"].(string)

	// ユーザー情報を取得
	user := models.User{}
	if err := database.DB.Where("token = ?", token).First(&user).Error; err != nil {
		// ユーザーIDを生成
		token := gt.GenUserToken()
		userID := issueUserJWT(token)

		// ユーザー情報を保存
		user := models.User{
			Token: token,
		}
		if err := database.DB.Create(&user).Error; err != nil {
			return c.JSON(http.StatusOK, epr.APIError("ユーザー登録に失敗しました。"))
		}

		// ユーザーIDを返す
		response := types.User{
			ID: userID,
		}
		return c.JSON(http.StatusOK, response)
	}

	// ユーザー情報を返す
	response := types.User{
		ID:        strconv.FormatUint(uint64(user.ID), 10),
		IsOrdered: user.IsOrdered,
	}
	return c.JSON(http.StatusOK, response)
}

func InviteRegistry(c echo.Context) error {
	// ユーザーIDを取得
	userId := c.Param("userId")

	// ユーザー情報を取得
	user := models.User{}
	if err := database.DB.Where("token = ?", userId).First(&user).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("ユーザーIDが見つかりません。"))
	}

	// ユーザーのisInvitationをtrueに更新
	user.IsInvitation = true
	if err := database.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("ユーザーの招待状態を更新できませんでした。"))
	}

	// ユーザーのisInvitationを返す
	response := true
	return c.JSON(http.StatusOK, response)
}

func DrawBulkLots(c echo.Context) error {
	// ユーザーIDを取得
	jwtToken := c.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(jwt.MapClaims)
	token := claims["sub"].(string)

	// ユーザー情報を取得
	user := models.User{}
	if err := database.DB.Where("token = ?", token).First(&user).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("ユーザーIDが見つかりません。"))
	}

	// ユーザーのbulk_couponを取得
	bulkCoupon := user.BulkCoupon

	// bulk_couponがnoneの場合
	if bulkCoupon == "none" {
		// ランダムにkindを生成
		n := rand.Intn(100)
		var kind string

		if n < 20 { // 20%
			kind = "100"
		} else if n < 26 { // 6%
			kind = "200"
		} else if n < 30 { // 4%
			kind = "300"
		} else {
			kind = "0"
		}

		// bulk_couponを更新
		user.BulkCoupon = kind
		if err := database.DB.Save(&user).Error; err != nil {
			return c.JSON(http.StatusOK, epr.APIError("bulk_couponの更新に失敗しました。"))
		}

		// 生成したkindを返す
		response := types.Coupon{
			Kind: kind,
		}
		return c.JSON(http.StatusOK, response)
	}

	// bulk_couponがnoneでない場合
	// そのまま返す
	response := types.Coupon{
		Kind: bulkCoupon,
	}
	return c.JSON(http.StatusOK, response)
}

func DrawInviteLots(c echo.Context) error {
	// ユーザーIDを取得
	jwtToken := c.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(jwt.MapClaims)
	token := claims["sub"].(string)

	// ユーザー情報を取得
	user := models.User{}
	if err := database.DB.Where("token = ?", token).First(&user).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("ユーザーIDが見つかりません。"))
	}

	// ユーザーが誰かを招待したかを調べる
	if !user.IsInvitation {
		return c.JSON(http.StatusOK, epr.APIError("QRコードをほかの人に読み込んでもらってください。"))
	}

	// ユーザーのinviteCouponを取得
	inviteCoupon := user.InviteCoupon

	// inviteCouponがnoneの場合
	if inviteCoupon == "none" {
		// ランダムにkindを生成
		n := rand.Intn(100)
		var kind string

		if n < 20 { // 20%
			kind = "100"
		} else if n < 26 { // 6%
			kind = "200"
		} else if n < 30 { // 4%
			kind = "300"
		} else {
			kind = "0"
		}

		// inviteCouponを更新
		user.InviteCoupon = kind
		if err := database.DB.Save(&user).Error; err != nil {
			return c.JSON(http.StatusOK, epr.APIError("inviteCouponの更新に失敗しました。"))
		}

		// 生成したkindを返す
		response := types.Coupon{
			Kind: kind,
		}
		return c.JSON(http.StatusOK, response)
	}

	// inviteCouponがnoneでない場合
	// そのまま返す
	response := types.Coupon{
		Kind: inviteCoupon,
	}
	return c.JSON(http.StatusOK, response)
}

func GetCouponItemIds(c echo.Context) error {
	// クーポン種別とクーポンIDの対応表を取得
	oneHundredCouponItemId := "16"
	twoHundredCouponItemId := "17"
	threeHundredCouponItemId := "18"
	couponItemIds := types.CouponItemIds{
		None:         nil,
		Zero:         nil,
		OneHundred:   &oneHundredCouponItemId,
		TwoHundred:   &twoHundredCouponItemId,
		ThreeHundred: &threeHundredCouponItemId,
	}

	return c.JSON(http.StatusOK, couponItemIds)
}

func GetCompleteState(c echo.Context) error {
	// ユーザーIDを取得
	jwtToken := c.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(jwt.MapClaims)
	token := claims["sub"].(string)

	// ユーザー情報を取得
	user := models.User{}
	if err := database.DB.Where("token = ?", token).First(&user).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("ユーザー情報が見つかりません。"))
	}

	// ユーザーIDを取得
	userId := user.ID

	// ユーザーの注文情報を取得
	order := models.Order{}
	if err := database.DB.Where("user_id = ?", userId).First(&order).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("注文情報が見つかりません。"))
	}

	// 注文状態を取得
	orderStatus := order.OrderStatus

	// 注文状態に応じて完了状態を返す
	switch orderStatus {
	case "ordered":
		return c.JSON(http.StatusOK, types.CompleteState{
			State: "Cooking",
		})
	case "cooked":
		return c.JSON(http.StatusOK, types.CompleteState{
			State: "Cooked",
		})
	case "received":
		return c.JSON(http.StatusOK, types.CompleteState{
			State: "Delivered",
		})
	default:
		return c.JSON(http.StatusOK, types.CompleteState{
			State: "Unknown",
		})
	}
}

func SendCartData(c echo.Context) error {
	// カートを取得
	cart := c.QueryParam("cart")

	// ユーザーIDを取得
	jwtToken := c.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(jwt.MapClaims)
	token := claims["sub"].(string)

	// カートをJSON形式に変換
	cartItems := []types.CartItem{}
	err := json.Unmarshal([]byte(cart), &cartItems)
	if err != nil {
		return c.JSON(http.StatusBadRequest, epr.APIError("カートのJSON形式が不正です。"))
	}

	// カートの中の商品の数を数える
	var cartItemsCount int
	for _, cartItem := range cartItems {
		cartItemsCount += cartItem.Quantity
	}

	// ユーザー情報を取得
	user := models.User{}
	if err := database.DB.Where("token = ?", token).First(&user).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("ユーザー情報が見つかりません。"))
	}

	// ユーザーの注文状況を取得
	userLatestOrder := models.Order{}
	if err := database.DB.Where("user_id = ?", user.ID).First(&userLatestOrder).Error; err == nil {
		return c.JSON(http.StatusOK, epr.APIError("すでに注文が完了しています。"))
	} else if err.Error() != "record not found" {
		return c.JSON(http.StatusOK, epr.APIError("注文情報の取得に失敗しました。"))
	}

	// time_of_completionをほかの注文状況から求める
	// 注文情報を取得
	var latestCompletionTime time.Time
	latestOrder := models.Order{}
	if err := database.DB.Where("order_status = ?", "ordered").Order("created_at desc").First(&latestOrder).Error; err != nil {
		latestCompletionTime = time.Now()
	} else {
		latestCompletionTime = latestOrder.TimeOfCompletion
	}

	timeOfCompletion := latestCompletionTime
	timeOfCompletion = timeOfCompletion.Add(time.Duration(math.Ceil(float64(cartItemsCount)/3)) * 4 * time.Minute)

	// 注文を作成
	order := models.Order{
		UserID:           user.ID,
		OrderStatus:      "ordered",
		IsMobileOrder:    true,
		TimeOfCompletion: timeOfCompletion,
	}
	if err := database.DB.Create(&order).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("注文の作成に失敗しました。"))
	}

	// 注文明細を作成
	for _, cartItem := range cartItems {
		// 注文明細を作成
		menuId, _ := strconv.Atoi(cartItem.ID)
		orderItem := models.OrderItem{
			OrderID:  order.ID,
			MenuID:   uint32(menuId),
			Quantity: cartItem.Quantity,
		}
		if err := database.DB.Create(&orderItem).Error; err != nil {
			return c.JSON(http.StatusOK, epr.APIError("注文明細の作成に失敗しました。"))
		}
	}

	// ユーザーのisOrderをtrueにする
	if err := database.DB.Model(&user).Update("is_ordered", true).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("ユーザー情報の更新に失敗しました。"))
	}

	// 注文情報を返却
	return c.JSON(http.StatusOK, true)
}

func genBarcode() string {
	const length = 20

	var barcode string

	for i := 0; i < length; i++ {
		n := rand.Intn(10)

		if i == 0 {
			n = 3
		}
		if i == 1 {
			n = 9
		}

		barcode += strconv.Itoa(n)
	}

	return barcode
}

func GetCompleteInfo(c echo.Context) error {
	// リクエストヘッダーからJWTトークンを取得
	jwtToken := c.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(jwt.MapClaims)
	token := claims["sub"].(string)

	// ユーザー情報を取得
	user := models.User{}
	if err := database.DB.Where("token = ?", token).First(&user).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("ユーザー情報が見つかりません。"))
	}

	// ユーザーの注文情報を取得
	order := models.Order{}
	err := database.DB.Where("user_id = ?", user.ID).First(&order).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, epr.APIError("注文情報が見つかりません。"))
	}

	// 注文明細情報を取得
	items := []models.OrderItem{}
	err = database.DB.Where("order_id = ?", order.ID).Find(&items).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, epr.APIError("注文明細の取得に失敗しました。"))
	}

	// 注文明細情報を整形
	var cartItems []types.CartItem
	for _, item := range items {
		cartItems = append(cartItems, types.CartItem{
			ID:       strconv.FormatUint(uint64(item.MenuID), 10),
			Quantity: item.Quantity,
		})
	}

	// バーコードを生成
	barcode := genBarcode()

	// 完了情報を返却
	info := types.CompleteInfo{
		Barcode:      barcode,
		CompleteTime: order.TimeOfCompletion,
		Items:        cartItems,
	}
	return c.JSON(http.StatusOK, info)
}
