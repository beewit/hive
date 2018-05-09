package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/enum"
)

func GetIndustryChild(c echo.Context) error {
	var parentId = c.FormValue("parentId")
	if parentId == "" {
		parentId = "0"
	}
	if !utils.IsValidNumber(parentId) {
		return utils.ErrorNull(c, "parentId错误")
	}
	rows, err := global.DB.Query("SELECT * FROM industry WHERE parent_id=? AND status=?", parentId, enum.NORMAL)
	if err != nil {
		return utils.ErrorNull(c, "获取城市地区失败")
	}
	return utils.SuccessNullMsg(c, rows)
}

func GetIndustry(c echo.Context) error {
	var id = c.FormValue("id")
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id错误")
	}
	rows, err := global.DB.Query("SELECT * FROM industry WHERE id=? AND status=? LIMIT 1", id, enum.NORMAL)
	if err != nil {
		return utils.ErrorNull(c, "获取城市地区失败")
	}
	if len(rows) != 1 {
		return utils.ErrorNull(c, "无数据")
	}
	return utils.SuccessNullMsg(c, rows[0])
}
