/*
发送红包
author：zxb
2018-01-13
*/
package handler

import (
	"fmt"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/hive/global"
	"github.com/labstack/echo"
	"strings"
)

func AddCoupon(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	money := strings.TrimSpace(c.FormValue("money"))
	condition := strings.TrimSpace(c.FormValue("condition"))
	expireTime := strings.TrimSpace(c.FormValue("expireTime"))
	number := strings.TrimSpace(c.FormValue("number"))
	desc := strings.TrimSpace(c.FormValue("desc"))
	if money == "" || !utils.IsValidNumber(money) {
		return utils.ErrorNull(c, "请正确填写代金券抵扣金额")
	}
	if condition != "" && len(condition) > 50 {
		return utils.ErrorNull(c, "使用条件字数过长，最长不超过50字")
	}
	if expireTime != "" && !utils.IsValidDate(expireTime) {
		return utils.ErrorNull(c, "到期时间格式错误")
	}
	if number == "" || !utils.IsValidNumber(number) {
		return utils.ErrorNull(c, "代金券数量格式错误")
	}
	if desc != "" && len(desc) > 500 {
		return utils.ErrorNull(c, "使用说明字数过长，最长不超过500字")
	}
	id := utils.ID()
	coupon := map[string]interface{}{}
	coupon["id"] = id
	coupon["account_id"] = acc.ID
	coupon["money"] = money
	if condition != "" {
		coupon["condition"] = condition
	}
	if expireTime != "" {
		coupon["expire_time"] = expireTime
	}
	if desc != "" {
		coupon["desc"] = desc
	}
	coupon["ct_time"] = utils.CurrentTime()
	coupon["ut_time"] = coupon["ct_time"]
	coupon["ip"] = c.RealIP()
	_, err = global.DB.InsertMap("account_coupon", coupon)
	if err != nil {
		return utils.ErrorNull(c, "发代金券失败")
	}
	return utils.Success(c, "发代金券成功", id)
}

func GetCouponList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	t := c.FormValue("t")
	where := ""
	switch t {
	case "expire":
		where = " AND expire_time<now()"
		break
	case "finish":
		where = " AND number=receive_number"
		break
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "account_coupon",
		Where:     "account_id=? AND status=?" + where,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, acc.ID, enum.NORMAL)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func GetCouponByID(c echo.Context) error {
	id := c.FormValue("id")
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id格式错误")
	}
	m := GetCoupon(convert.MustInt64(id))
	if m == nil {
		return utils.NullData(c)
	}
	return utils.SuccessNullMsg(c, m)
}

func GetCoupon(id int64) map[string]interface{} {
	m, err := global.DB.Query("SELECT * FROM account_coupon WHERE status=? AND id=? LIMIT 1", enum.NORMAL, id)
	if err != nil {
		global.Log.Error("getCouponByID sql ERROR：", err.Error())
		return nil
	}
	if len(m) != 1 {
		return nil
	}
	return m[0]
}

func GetCouponByIds(ids string) []map[string]interface{} {
	if ids == "" {
		return nil
	}
	m, err := global.DB.Query(fmt.Sprintf("SELECT * FROM account_coupon WHERE status=? AND id in(%s) LIMIT 1", ids), enum.NORMAL)
	if err != nil {
		global.Log.Error("getCouponByID sql ERROR：", err.Error())
		return nil
	}
	if len(m) < 1 {
		return nil
	}
	return m
}
