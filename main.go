package main

import (
	"net/http"
	"os"

	epr "github.com/dokimiki/scfes23-foodorder_backend/libs/errorPayloadResponse"
	stores_route "github.com/dokimiki/scfes23-foodorder_backend/routes/stores"
	ticket_route "github.com/dokimiki/scfes23-foodorder_backend/routes/ticket"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 1,
	}))
	e.Use(middleware.Secure())
	// e.Pre(middleware.HTTPSRedirect()) // https化した際に有効化する

	v1 := e.Group("/v1")
	v1.GET("", Hello)

	ticket := v1.Group("/ticket")
	ticket.GET("/price", ticket_route.GetPrice) // TODO: Replace HELLO to get price

	stores := v1.Group("/stores/:store_id")
	stores.GET("", stores_route.GetStoreInfo)      // TODO: Replace HELLO to get store info
	stores.GET("/menus", stores_route.GetMenuList) // TODO: Replace HELLO to get menus
	stores.GET("/menus/:menu_id", Hello)           // TODO: Replace HELLO to get menu

	user := v1.Group("/user")
	user.POST("/signup", Hello) // TODO: Replace HELLO to signup

	user.Use(echojwt.JWT([]byte(signature)))
	user.GET("/me", Hello)                          // TODO: Replace HELLO to get user info
	user.GET("/me/orders", Hello)                   // TODO: Replace HELLO to get user orders
	user.GET("/me/orders/:order_id", Hello)         // TODO: Replace HELLO to get user order
	user.POST("/me/orders/:order_id/cancel", Hello) // TODO: Replace HELLO to cancel user order

	storekeeper := v1.Group("/storekeeper")
	storekeeper.POST("/request", Hello) // TODO: Replace HELLO to request storekeeper

	storekeeper.Use(echojwt.JWT([]byte(signature)))
	storekeeper.POST("/approve", Hello) // TODO: Replace HELLO to approve storekeeper
	storekeeper.GET("/me", Hello)       // TODO: Replace HELLO to get storekeeper info

	e.Logger.Fatal(e.Start(":3030"))
}
