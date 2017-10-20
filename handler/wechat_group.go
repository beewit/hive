package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/enum"
	"fmt"
)

func GetWechatGroupList(c echo.Context) error {
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	where := "status=?"
	area := c.FormValue("area")
	types := c.FormValue("type")
	if area != "" {
		where += fmt.Sprintf(" AND area='%s'", area)
	}
	if types != "" {
		where += fmt.Sprintf(" AND type='%s'", types)
	}

	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "wechat_group",
		Where:     where,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Order:     "ut_time DESC",
	}, enum.NORMAL)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}
