package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/hive/global"
	"github.com/labstack/echo"
	"fmt"
)

func GetFuncList(c echo.Context) error {
	accID := c.FormValue("accId")
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "f.*,af.expiration_time",
		Table:     "func f LEFT JOIN account_func af ON af.func_id=f.id AND af.account_id=?",
		Where:     "f.status=?",
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, accID, enum.NORMAL)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func GetAccountFuncList(c echo.Context) error {
	accID := c.FormValue("accId")
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "f.*,af.expiration_time",
		Table:     "account_func af LEFT JOIN func f ON af.func_id=f.id",
		Where:     "f.status=? AND af.account_id=?",
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, enum.NORMAL, accID)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func GetFuncAllByIdsAndAccId(c echo.Context) error {
	accID := c.FormValue("accId")
	funcIds := c.FormValue("funcIds")
	platformIds := c.FormValue("platformIds")
	where := "af.expiration_time>now() AND f.status=? AND af.account_id=?"
	if funcIds != "" {
		where += fmt.Sprintf(" AND f.id in(%s)", funcIds)
	}
	if platformIds != "" {
		where += fmt.Sprintf(" AND f.platform_id in(%s)", platformIds)
	}
	sql := fmt.Sprintf(
		"SELECT f.* FROM account_func af LEFT JOIN func f ON af.func_id=f.id LEFT JOIN platform p ON p.id=f.platform_id WHERE %s",
		where)
	m, err := global.DB.Query(sql, enum.NORMAL, accID)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	return utils.Success(c, "获取数据成功", m)
}

func GetEffectiveFuncById(c echo.Context) error {
	itf := c.Get("account")
	if itf == nil {
		return utils.AuthFailNull(c)
	}
	acc := global.ToInterfaceAccount(itf)
	if acc == nil {
		return utils.AuthFailNull(c)
	}
	funcId := c.FormValue("funcId")
	if funcId == "" {
		return utils.ErrorNull(c, "funcId不能为空")
	}
	if !utils.IsValidNumber(funcId) {
		return utils.ErrorNull(c, "funcId参数错误")
	}
	sql := "SELECT f.* FROM account_func af LEFT JOIN func f ON af.func_id=f.id " +
		"LEFT JOIN platform p ON p.id=f.platform_id " +
		"WHERE af.expiration_time>now() AND f.status=? AND af.account_id=? AND f.id=? LIMIT 1"
	m, err := global.DB.Query(sql, enum.NORMAL, acc.ID, funcId)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if len(m) != 1 {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", m[0])
}

func GetFuncGroupByAccountId(c echo.Context) error {
	itf := c.Get("account")
	if itf == nil {
		return utils.AuthFailNull(c)
	}
	acc := global.ToInterfaceAccount(itf)
	if acc == nil {
		return utils.AuthFailNull(c)
	}
	openSql := "SELECT * FROM account_func WHERE expiration_time>NOW() AND account_id=?"
	expireSql := "SELECT * FROM account_func WHERE expiration_time<NOW() AND account_id=?"
	notExpireSql := "SELECT * FROM account_func WHERE expiration_time>date_add(NOW(),interval 15 DAY) AND account_id=?"
	openMap, _ := global.DB.Query(openSql, acc.ID)
	expireMap, _ := global.DB.Query(expireSql, acc.ID)
	notExpireMap, _ := global.DB.Query(notExpireSql, acc.ID)
	data := map[string]interface{}{
		"openMap":      openMap,
		"expireMap":    expireMap,
		"notExpireMap": notExpireMap,
	}
	return utils.SuccessNullMsg(c, data)
}
