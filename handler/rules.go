package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils"
	"time"
	"github.com/beewit/beekit/utils/enum"
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
	if acc.MemberTypeId <= 0 || acc.MemberExpirTime == "" {
		//非会员
		r["member"] = enum.NOT_MEMBER
	} else {
		t := time.Now()
		t1, errTime := time.Parse("2006-01-02 15:04:05", acc.MemberExpirTime)
		if errTime == nil && t.After(t1) {
			//已过期
			r["member"] = enum.MEMBER_EXPIRED
		} else {
			r["member"] = enum.MEMBER_NOT_EXPIRED
		}
	}
	sql := "SELECT r.*,group_concat(rmt.member_type_id),group_concat(rmt.member_type_name) FROM rules r LEFT JOIN " +
		"rules_member_type rmt ON r.id=rmt.rules_id GROUP BY r.id"
	m, err := global.DB.Query(sql)
	if err != nil {
		global.Log.Error(err.Error())
		return utils.ErrorNull(c, "查询会员规则失败")
	}
	r["rules"] = m
	return utils.SuccessNullMsg(c, r)
}
