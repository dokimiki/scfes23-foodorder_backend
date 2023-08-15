package stores_route

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/dokimiki/scfes23-foodorder_backend/database"
	epr "github.com/dokimiki/scfes23-foodorder_backend/libs/errorPayloadResponse"
	"github.com/dokimiki/scfes23-foodorder_backend/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetStoreInfo(c echo.Context) error {
	type res struct {
		StrID    string `json:"id"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Features string `json:"features"`
	}

	StoreId := c.Param("store_id")
	r := res{}

	if err := database.DB.Model(models.Store{}).Where("str_id = ?", StoreId).Scan(&r).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusOK, epr.APIError("指定された店舗がありません(store is not exist)"))
	} else if err != nil {
		return c.JSON(http.StatusOK, epr.APIError("エラーが発生しました(error occurred)"))
	}

	return c.JSON(http.StatusOK, r)
}

func GetMenuList(c echo.Context) error {
	type res struct {
		StrID       string `json:"id"`
		Name        string `json:"name"`
		ImgURL      string `json:"imageUrl"`
		TicketPrice int    `json:"price"`
		Discount    int    `json:"discount"`
	}

	StoreID := c.Param("store_id")
	r := []res{}

	if err := database.DB.Model(models.Menu{}).
		Select("menus.str_id, menus.name, menus.img_url, menu_details.ticket_price, menu_details.discount").
		Where("store_str_id = ?", StoreID).
		Joins("INNER JOIN menu_details ON menus.id = menu_details.menu_id").
		Scan(&r).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("エラーが発生しました(error occurred)"))
	}

	return c.JSON(http.StatusOK, r)
}

func GetMenuDetail(c echo.Context) error {
	type res struct {
		StrID       string `json:"id"`
		Name        string `json:"name"`
		Features    string `json:"features"`
		ImgURL      string `json:"imageUrl"`
		Remaining   int    `json:"remaining"`
		TicketPrice int    `json:"price"`
		Discount    int    `json:"discount"`
		Allergen    struct {
			Ebi    string `json:"ebi"`
			Kani   string `json:"kani"`
			Komugi string `json:"komugi"`
			Kurumi string `json:"kurumi"`
			Milk   string `json:"milk"`
			Peanut string `json:"peanut"`
			Soba   string `json:"soba"`
			Tamago string `json:"tamago"`
		} `json:"allergen" gorm:"-"`
	}

	StoreId := c.Param("store_id")
	MenuId := c.Param("menu_id")
	txOptions := &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}
	r := res{}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		// menuの取得
		menu := models.Menu{}
		if err := tx.Where("str_id = ? AND store_str_id = ?", MenuId, StoreId).Take(&menu).Scan(&r).Error; err != nil {
			return err
		}

		// menu detailの取得
		if err := tx.Model(models.MenuDetail{}).Where("menu_id = ?", menu.ID).Scan(&r).Error; err != nil {
			return err
		}

		// menu allergenの取得
		if err := tx.Model(models.MenuAllergen{}).Where("menu_id = ?", menu.ID).Scan(&r.Allergen).Error; err != nil {
			return err
		}

		return nil
	}, txOptions); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusOK, epr.APIError("指定されたメニューがありません(menu is not exist)"))
		} else {
			return c.JSON(http.StatusOK, epr.APIError("エラーが発生しました(error occurred)"))
		}
	}

	return c.JSON(http.StatusOK, r)
}
