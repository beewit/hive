package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/beekit/utils/convert"
	"fmt"
)

/*
*  系统帮助
 */
func GetHelpList(c echo.Context) error {
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	t := c.FormValue("t")
	o := c.FormValue("order")
	where := ""
	if t != "" {
		where += fmt.Sprintf(" AND type='%s'", t)
	}
	order := "ct_time DESC"
	if o == "hot" {
		order = "page_view DESC"
	}
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "id,type,img,title,page_view,ct_time,status,href",
		Table:     "help",
		Where:     "status=?" + where,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Order:     order,
	},   enum.NORMAL)
	if err != nil {
		global.Log.Error("GetHelpList sql error:%s", err.Error())
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func GetHelp(c echo.Context) error {
	id := c.FormValue("id")
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "数据参数错误")
	}
	m := GetHelpById(convert.MustInt64(id))
	if m == nil {
		return utils.NullData(c)
	}
	return utils.SuccessNullMsg(c, m)
}

func GetHelpById(id int64) map[string]interface{} {
	m, err := global.DB.Query("SELECT * FROM help WHERE status=? AND id=? LIMIT 1", enum.NORMAL, id)
	if err != nil {
		global.Log.Error("GetHelpById sql ERROR：", err.Error())
		return nil
	}
	if len(m) != 1 {
		return nil
	}
	return m[0]
}
