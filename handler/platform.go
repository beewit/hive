package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/enum"
)

func GetPlatformList(c echo.Context) error {
	sql := `SELECT * FROM platform WHERE status = ?`
	rows, err := global.DB.Query(sql, enum.NORMAL)
	if err != nil {
		global.Log.Error("GetPlatformList：" + err.Error())
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if len(rows) <= 0 {
		return utils.Success(c, "无数据", nil)
	}
	return utils.Success(c, "获取数据成功", convert.ToArrayMapString(rows))
}

func GetPlatformId(c echo.Context) error {
	t := c.FormValue("type")
	sql := `SELECT * FROM platform WHERE type=? AND status = ? LIMIT 1`
	rows, err := global.DB.Query(sql, t, enum.NORMAL)
	if err != nil {
		global.Log.Error("GetPlatformId：" + err.Error())
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if len(rows) != 1 {
		return utils.Success(c, "无数据", nil)
	}
	return utils.Success(c, "获取数据成功", convert.ToMapString(rows[0]))
}
