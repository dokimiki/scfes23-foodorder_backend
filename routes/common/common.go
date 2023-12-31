package cr

import (
	"net/http"
	"strconv"

	"github.com/dokimiki/scfes23-foodorder_backend/database"
	epr "github.com/dokimiki/scfes23-foodorder_backend/libs/errorPayloadResponse"
	"github.com/dokimiki/scfes23-foodorder_backend/models"
	"github.com/dokimiki/scfes23-foodorder_backend/types"
	"github.com/labstack/echo/v4"
)

func GetMenuItems(c echo.Context) error {
	dbResult := []models.Menu{}

	// FIND all menus
	if err := database.DB.Find(&dbResult).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("メニュー取得でエラーが発生しました。"))
	}

	// Convert to response
	response := []types.MenuItem{}
	for _, menu := range dbResult {
		response = append(response, types.MenuItem{
			ID:     strconv.FormatUint(uint64(menu.ID), 10),
			Name:   menu.Name,
			Price:  menu.Price,
			Image:  menu.ImgUrl,
			IsShow: menu.IsShow,
		})
	}

	// JSONで返す
	return c.JSON(http.StatusOK, response)
}

func GetAllergen(c echo.Context) error {
	id := c.Param("menuId")

	// アレルギー情報を取得する
	allergens := models.MenuAllergen{}
	if err := database.DB.Where("menu_id = ?", id).Take(&allergens).Error; err != nil {
		return c.JSON(http.StatusOK, epr.APIError("アレルギー情報取得でエラーが発生しました。"))
	}

	// アレルギー情報を返す
	response := types.AllergensList{
		Ebi:    allergens.Ebi,
		Kani:   allergens.Kani,
		Komugi: allergens.Komugi,
		Kurumi: allergens.Kurumi,
		Milk:   allergens.Milk,
		Peanut: allergens.Peanut,
		Soba:   allergens.Soba,
		Tamago: allergens.Tamago,
	}

	return c.JSON(http.StatusOK, response)
}
