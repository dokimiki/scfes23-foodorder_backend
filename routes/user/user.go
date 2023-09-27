package ur

import (
	"net/http"
	"os"
	"time"

	"github.com/dokimiki/scfes23-foodorder_backend/database"
	epr "github.com/dokimiki/scfes23-foodorder_backend/libs/errorPayloadResponse"
	gt "github.com/dokimiki/scfes23-foodorder_backend/libs/generateToken"
	"github.com/dokimiki/scfes23-foodorder_backend/models"
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
	type res struct {
		JWT string `json:"JWT"`
	}

	token := gt.GenUserToken()
	r := res{JWT: issueUserJWT(token)}

	// アカウント作成
	user := models.User{
		Token: token,
	}
	if err := database.DB.Select("token").Create(&user).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("アカウント作成に失敗しました(error on create account)"))
	}

	return c.JSON(http.StatusOK, r)
}

func UserInfo(c echo.Context) error {
	type res struct {
		ID          uint32 `json:"id"`
		CancelCount int    `json:"cancelCount"`
	}

	jwtToken := c.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(jwt.MapClaims)
	token := claims["sub"].(string)
	r := res{}

	if err := database.DB.Model(models.User{}).Where("token = ?", token).Scan(&r).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("ユーザー情報の取得に失敗しました(error on get user info)"))
	}

	return c.JSON(http.StatusOK, r)
}
