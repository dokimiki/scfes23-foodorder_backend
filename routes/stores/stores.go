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
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Features string `json:"features"`
}

type menuListInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Features string `json:"features"`
	ImageURL string `json:"imageUrl"`
	Price    int    `json:"price"`
	Discount int    `json:"discount"`
}

func checkStoreIsExist(storeID string) (uint32, error) {
	var count int
	database.DB.Model(models.Store{}).
		Select("COUNT(*)").
		Where("str_id = ?", storeID).
		Limit(1).
		Scan(&count)

	if count > 0 {
		result := models.Store{}
		database.DB.Model(models.Store{}).
			Where("str_id = ?", storeID).
			Limit(1).
			Scan(&result)
		return result.ID, nil
	} else {
		return 0, errors.New("store is not exist")
	}
}

func GetStoreInfo(c echo.Context) error {
	strStoreId := c.Param("store_id")
	intStoreId, err := checkStoreIsExist(strStoreId)
	if err != nil {
		return c.JSON(http.StatusOK, epr.APIError(err.Error()))
	}

	db_result := models.Store{}
	database.DB.Model(models.Store{}).Where("id = ?", intStoreId).Scan(&db_result)

	result := storeInfo{
		ID:       db_result.StrID,
		Name:     db_result.Name,
		Location: db_result.Location,
		Features: db_result.Features,
	}

	return c.JSON(http.StatusOK, result)
}

func GetMenuList(c echo.Context) error {
	strStoreId := c.Param("store_id")
	intStoreId, err := checkStoreIsExist(strStoreId)
	if err != nil {
		return c.JSON(http.StatusOK, epr.APIError(err.Error()))
	}

	// TODO: JOIN句を使って一回で取得するようにする
	menuListResult := []models.Menu{}
	menuDetailListResult := []models.MenuDetails{}
	database.DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(models.Menu{}).Where("store_id = ?", intStoreId).Scan(&menuListResult)

		var menuIdList []uint32

		for _, menu := range menuListResult {
			menuIdList = append(menuIdList, menu.ID)
		}

		tx.Model(models.MenuDetails{}).Where("menu_id IN ?", menuIdList).Scan(&menuDetailListResult)

		return nil
	})

	result := []menuListInfo{}

	for _, menu := range menuListResult {
		for _, menuDetail := range menuDetailListResult {
			if menu.ID == menuDetail.MenuID {
				result = append(result, menuListInfo{
					ID:       menu.StrID,
					Name:     menu.Name,
					Features: menu.Features,
					ImageURL: menu.ImgURL,
					Price:    menuDetail.TicketPrice,
					Discount: menuDetail.Discount,
				})
			}
		}
	}

	return c.JSON(http.StatusOK, result)
}
