/*
领取代金券记录
author：zxb
2018-01-13
*/

package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/hive/global"
	"github.com/labstack/echo"
	"time"
)

func ReceiveCoupon(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	id := c.FormValue("id")
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id格式错误")
	}
	coupon := GetCoupon(convert.MustInt64(id))
	if coupon != nil {
		return utils.ErrorNull(c, "该代金券不存在")
	}
	sendCouponAcc := GetAccountById(convert.MustInt64(coupon["account_id"]))
	if sendCouponAcc == nil {
		return utils.ErrorNull(c, "发代金券的账号已被锁定无法领取")
	}
	if coupon["status"] != enum.NORMAL {
		return utils.ErrorNull(c, "该代金券已失效")
	}
	if coupon["number"] == coupon["receive_number"] {
		return utils.ErrorNull(c, "该代金券已领完")
	}
	if convert.ToString(coupon["expire_time"]) != "" {
		expirTime, errExpirTime := time.Parse("2006-01-02 15:04:05", convert.ToString(coupon["expire_time"]))
		if errExpirTime != nil {
			return utils.ErrorNull(c, "该代金券到期时间异常")
		}
		if !expirTime.After(time.Now()) {
			return utils.ErrorNull(c, "该代金券已到期无法领取")
		}
	}

	_, err = global.DB.InsertMap("account_receive_coupon", map[string]interface{}{
		"id":                utils.ID,
		"account_id":        acc.ID,
		"account_coupon_id": coupon["id"],
		"ct_time":           utils.CurrentTime(),
		"ip":                c.RealIP(),
		"status":            "未使用",
	})
	if err != nil {
		return utils.ErrorNull(c, "领取代金券失败")
	}
	return utils.SuccessNull(c, "领取代金券成功")
}

func GetReceiveCouponList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	t := c.FormValue("t")
	where := "status='未使用' AND (expire_time is NULL OR expire_time>now())"
	switch t {
	case "expire":
		where = "expire_time<now()"
		break
	case "finish":
		where = "status='未使用'"
		break
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "v_account_receive_coupon",
		Where:     "account_id=? AND " + where,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, acc.ID)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}
