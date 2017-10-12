package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/enum"
)

func GetPayOrderList(c echo.Context) error {
	accID := c.FormValue("accId")
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "o.*",
		Table:     "order_payment o",
		Where:     "o.status=? AND o.pay_status=? AND o.account_id=?",
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, enum.NORMAL, enum.PAY_STATUS_END, accID)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}
