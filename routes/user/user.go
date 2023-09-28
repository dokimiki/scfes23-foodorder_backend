package ur

import (
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
	rand.Seed(time.Now().UnixNano())
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
		fmt.println(n)

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
