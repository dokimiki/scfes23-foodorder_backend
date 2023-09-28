package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	epr "github.com/dokimiki/scfes23-foodorder_backend/libs/errorPayloadResponse"
	cr "github.com/dokimiki/scfes23-foodorder_backend/routes/common"
	ur "github.com/dokimiki/scfes23-foodorder_backend/routes/user"
)

func Hello(c echo.Context) error {
	return c.JSON(http.StatusOK, epr.APIError("Hello, World!"))
}

func main() {
	godotenv.Load()
	signature := os.Getenv("SCFES23FOODORDER_JWT_SIGNATURE")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 1,
	}))
	e.Use(middleware.Secure())

	v1 := e.Group("/v1")
	e.IPExtractor = echo.ExtractIPFromXFFHeader()
	v1.GET("", Hello)

	/* Common */
	common := v1.Group("/common")
	common.GET("/menus", cr.GetMenuItems)
	common.GET("/allergens/:menuId", cr.GetAllergen)

	/* User */
	user := v1.Group("/user")
	user.POST("/signup", ur.SignUp)
	user.POST("/inviteregistry/:userId", ur.InviteRegistry)

	userWithAuth := user.Group("/me")
	userWithAuth.Use(echojwt.JWT([]byte(signature)))
	userWithAuth.GET("/signin", ur.SignIn)
	userWithAuth.GET("/drawbulklots", ur.DrawBulkLots)
	userWithAuth.GET("/drawinvitelots", ur.DrawInviteLots)
	userWithAuth.GET("/getcouponitemids", ur.GetCouponItemIds)
	userWithAuth.GET("/getcompletestate", ur.GetCompleteState)
	userWithAuth.GET("/sendcartdata", ur.SendCartData)
	userWithAuth.GET("/getcompleteinfo", ur.GetCompleteInfo)

	e.Logger.Fatal(e.StartTLS(":3939", "server.crt", "server.key"))
}
