func getMenuItems(c echo.Context) error {
	var menus []MenuItem

    // FIND all menus
    if err := db.Find(&menus).Error; err != nil {
        return c.JSON(http.StatusOK, epr.APIError("メニュー取得でエラーが発生しました。"))
    }

	// JSONで返す
	return c.JSON(http.StatusOK, menus)
})