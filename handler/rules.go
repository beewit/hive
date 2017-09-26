package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils"
)

func GetRules(c echo.Context) error {
	itf := c.Get("account")
	if itf == nil {
		return utils.AuthFailNull(c)
	}
	acc := global.ToInterfaceAccount(itf)
	if acc == nil {
		return utils.AuthFailNull(c)
	}
	r := make(map[string]interface{})
	sql := "SELECT r.*,group_concat(rmt.func_id),group_concat(rmt.func_name) FROM rules r LEFT JOIN " +
		"rules_func rmt ON r.id=rmt.rules_id GROUP BY r.id"
	m, err := global.DB.Query(sql)
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "查询会员规则失败")
	}
	r["rules"] = m
	return utils.SuccessNullMsg(c, r)
}
