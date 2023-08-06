package stores_route

import (
	"errors"
	"net/http"

	"github.com/dokimiki/scfes23-foodorder_backend/database"
	epr "github.com/dokimiki/scfes23-foodorder_backend/libs/errorPayloadResponse"
	"github.com/dokimiki/scfes23-foodorder_backend/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type storeInfo struct {
	StrID    string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Features string `json:"features"`
}

type menuListInfo struct {
	StrID       string `json:"id"`
	Name        string `json:"name"`
	ImgURL      string `json:"imageUrl"`
	TicketPrice int    `json:"price"`
	Discount    int    `json:"discount"`
}

type menuAllergenInfo struct {
	Ebi    string `json:"ebi"`
	Kani   string `json:"kani"`
	Komugi string `json:"komugi"`
	Kurumi string `json:"kurumi"`
	Milk   string `json:"milk"`
	Peanut string `json:"peanut"`
	Soba   string `json:"soba"`
	Tamago string `json:"tamago"`
}

type menuDetailInfo struct {
	StrID       string           `json:"id"`
	Name        string           `json:"name"`
	Features    string           `json:"features"`
	ImgURL      string           `json:"imageUrl"`
	Remaining   int              `json:"remaining"`
	TicketPrice int              `json:"price"`
	Discount    int              `json:"discount"`
	Allergen    menuAllergenInfo `json:"allergen"`
}

func checkStoreIsExist(storeID string) (uint32, error) {
	result := models.Store{}

	if err := database.DB.Where("str_id = ?", storeID).Take(&result).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, errors.New("指定された店舗がありません(store is not exist)")
	}

	return result.ID, nil
}

func GetStoreInfo(c echo.Context) error {
	strStoreId := c.Param("store_id")
	intStoreId, err := checkStoreIsExist(strStoreId)
	if err != nil {
		return c.JSON(http.StatusOK, epr.APIError(err.Error()))
	}

	result := storeInfo{}
	database.DB.Model(models.Store{}).Where("id = ?", intStoreId).Scan(&result)

	return c.JSON(http.StatusOK, result)
}

func GetMenuList(c echo.Context) error {
	strStoreId := c.Param("store_id")
	intStoreId, err := checkStoreIsExist(strStoreId)
	if err != nil {
		return c.JSON(http.StatusOK, epr.APIError(err.Error()))
	}

	result := []menuListInfo{}

	database.DB.Model(models.Menu{}).
		Select("menus.str_id, menus.name, menus.img_url, menu_details.ticket_price, menu_details.discount").
		Where("store_id = ?", intStoreId).
		Joins("INNER JOIN menu_details ON menus.id = menu_details.menu_id").
		Scan(&result)

	return c.JSON(http.StatusOK, result)
}

func GetMenuDetail(c echo.Context) error {
	strStoreId := c.Param("store_id")
	intStoreId, err := checkStoreIsExist(strStoreId)
	if err != nil {
		return c.JSON(http.StatusOK, epr.APIError(err.Error()))
	}

	strMenuId := c.Param("menu_id")

	// begin a transaction
	tx := database.DB.Begin()

	var menu models.Menu = models.Menu{}

	result := menuDetailInfo{}

	// TODO: 出てるエラーを直す
	if err := tx.Where("str_id = ? AND store_id = ?", strMenuId, intStoreId).Take(&menu).Scan(&result).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return c.JSON(http.StatusOK, epr.APIError("指定されたメニューがありません(menu is not exist)"))
	}
	menuId := menu.ID

	if err := tx.Model(models.MenuDetail{}).Where("menu_id = ?", menuId).Scan(&result).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return c.JSON(http.StatusOK, epr.APIError("指定されたメニューがありません(menu detail is not exist)"))
	}

	if err := tx.Model(models.MenuAllergen{}).Where("menu_id = ?", menuId).Scan(&result.Allergen).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return c.JSON(http.StatusOK, epr.APIError("指定されたメニューがありません(menu allergen is not exist)"))
	}

	// commit the transaction
	tx.Commit()

	return c.JSON(http.StatusOK, result)
}
