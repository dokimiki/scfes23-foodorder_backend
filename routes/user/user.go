package user_route

import (
	"crypto/sha512"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dokimiki/scfes23-foodorder_backend/database"
	epr "github.com/dokimiki/scfes23-foodorder_backend/libs/errorPayloadResponse"
	gt "github.com/dokimiki/scfes23-foodorder_backend/libs/generateToken"
	"github.com/dokimiki/scfes23-foodorder_backend/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type signUpBody struct {
	ScreenWidth  int `json:"screenWidth"`
	ScreenHeight int `json:"screenHeight"`
}

type JWTResponse struct {
	JWT string `json:"JWT"`
}

type UserInfoResponse struct {
	Token       string `json:"token"`
	CancelCount int    `json:"cancelCount"`
}

func issueUserJWT(id string) string {
	idToken := jwt.New(jwt.SigningMethodHS256)
	claims := idToken.Claims.(jwt.MapClaims)
	claims["sub"] = id
	claims["exp"] = time.Now().Add(time.Hour * 24 * 15).Unix()

	secret := os.Getenv("SCFES23FOODORDER_JWT_SIGNATURE")
	strIdToken, err := idToken.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}

	return strIdToken
}

func SignUp(c echo.Context) error {
	// デバイスの情報を取得
	body := signUpBody{}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, epr.APIError(err.Error()))
	}

	ip := c.RealIP()
	rh, _ := net.LookupAddr(ip)
	srh := rh[len(rh)-1]
	pt := fmt.Sprintf("%d,%d,%s,%s", body.ScreenWidth, body.ScreenHeight, ip, srh)
	hs := sha512.Sum512([]byte(pt))
	deviceInfo := models.Device{
		HashedData:   hs[:],
		ScreenWidth:  body.ScreenWidth,
		ScreenHeight: body.ScreenHeight,
		IPAddress:    ip,
		RemoteHost:   srh,
	}

	// トークン生成
	token := gt.GenUserToken()

	// アカウント作成
	user := models.User{
		Token: token,
	}
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("token").Create(&user).Error; err != nil {
			return err
		}

		deviceInfo.UserID = user.ID

		if err := tx.Create(&deviceInfo).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return c.JSON(http.StatusOK, epr.APIError(err.Error()))
	}

	JWT := issueUserJWT(fmt.Sprintf("%d", user.ID))

	return c.JSON(http.StatusOK, JWTResponse{JWT: JWT})
}

func UserInfo(c echo.Context) error {
	result := UserInfoResponse{}

	jwtToken := c.Get("user").(*jwt.Token)
	claims := jwtToken.Claims.(jwt.MapClaims)
	userId := claims["sub"].(string)

	userInfo := models.User{}
	if err := database.DB.Where("token = ?", userId).Take(&userInfo).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError(err.Error()))
	}

	result.Token = userInfo.Token
	result.CancelCount = userInfo.CancelCount

	return c.JSON(http.StatusOK, result)
}
