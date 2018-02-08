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
	"time"
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
	name := strings.TrimSpace(c.FormValue("name"))
	photo := strings.TrimSpace(c.FormValue("photo"))
	if name == "" {
		return utils.ErrorNull(c, "请填写商家名称")
	}
	if len(name) > 50 {
		return utils.ErrorNull(c, "商家名称，最长不超过500字")
	}
	if photo == "" {
		return utils.ErrorNull(c, "请上传发现金券的商家标志")
	}
	if photo != "" && len(photo) > 255 {
		return utils.ErrorNull(c, "商家标志地址错误")
	}

	if money == "" || !utils.IsValidNumber(money) {
		return utils.ErrorNull(c, "请正确填写现金券抵扣金额")
	}

	if convert.MustFloat64(money) <= 0 {
		return utils.ErrorNull(c, "现金券抵扣金额必须大于0")
	}

	if condition == "" || !utils.IsValidNumber(money) {
		return utils.ErrorNull(c, "请正确填写现金券使用条件")
	}
	if convert.MustFloat64(condition) <= 0 {
		return utils.ErrorNull(c, "使用条件满可用金额必须大于等于0")
	}

	if expireTime != "" && !utils.IsValidDate(expireTime) {
		return utils.ErrorNull(c, "到期时间格式错误")
	}
	if number == "" || !utils.IsValidNumber(number) {
		return utils.ErrorNull(c, "现金券数量格式错误")
	}
	if desc != "" && len(desc) > 500 {
		return utils.ErrorNull(c, "使用说明字数过长，最长不超过500字")
	}
	id := utils.ID()
	coupon := map[string]interface{}{}
	coupon["id"] = id
	coupon["account_id"] = acc.ID
	coupon["money"] = money
	coupon["number"] = number
	coupon["photo"] = photo
	coupon["name"] = name
	if condition != "" {
		coupon["`condition`"] = condition
	}
	if expireTime != "" {
		coupon["expire_time"] = expireTime
	}
	if desc != "" {
		coupon["`desc`"] = desc
	}
	coupon["ct_time"] = utils.CurrentTime()
	coupon["ut_time"] = coupon["ct_time"]
	coupon["ip"] = c.RealIP()
	_, err = global.DB.InsertMap("account_coupon", coupon)
	if err != nil {
		global.Log.Error("AddCoupon sql error:%s", err.Error())
		return utils.ErrorNull(c, "发现金券失败")
	}
	return utils.Success(c, "发现金券成功", id)
}

func GetCouponList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	t := c.FormValue("t")
	//未过期并未领完的
	where := " AND (expire_time is NULL OR expire_time>=DATE_FORMAT(now(),'%y-%m-%d')) AND number>receive_number"
	switch t {
	case "expire":
		//已过期
		where = " AND expire_time<DATE_FORMAT(now(),'%y-%m-%d')"
		break
	case "finish":
		//已领完
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
		Order:     "ct_time DESC",
	}, acc.ID, enum.NORMAL)
	if err != nil {
		global.Log.Error("QueryPage account_coupon sql error:%s", err.Error())
		return utils.ErrorNull(c, "数据异常")
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
	m, err := global.DB.Query(fmt.Sprintf("SELECT * FROM account_coupon WHERE status=? AND id in(%s)", ids), enum.NORMAL)
	if err != nil {
		global.Log.Error("getCouponByID sql ERROR：", err.Error())
		return nil
	}
	if len(m) < 1 {
		return nil
	}
	return m
}

//使用现金券
func UseCoupon(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	qrCodeKey := strings.TrimSpace(c.FormValue("qrCodeKey"))
	if qrCodeKey == "" {
		return utils.ErrorNull(c, "未识别有效现金券二维码")
	}
	keys := strings.Split(qrCodeKey, "|")
	if len(keys) < 2 {
		return utils.ErrorNull(c, "未识别有效现金券二维码")
	}
	id := keys[0]
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "未识别有效现金券二维码")
	}
	receiveCoupon := GetReceiveCoupon(convert.MustInt64(id))
	if receiveCoupon == nil {
		return utils.ErrorNull(c, "现金券不存在")
	}
	if convert.ToString(receiveCoupon["status"]) != enum.COUPON_USE_NOT {
		return utils.ErrorNull(c, "现金券已使用或已过期")
	}
	coupon := GetCoupon(convert.MustInt64(receiveCoupon["account_coupon_id"]))
	if coupon == nil {
		return utils.ErrorNull(c, "现金券已过期")
	}
	if convert.MustInt64(coupon["account_id"]) != acc.ID {
		return utils.ErrorNull(c, "使用失败，不属于我们的现金券")
	}
	if convert.ToString(coupon["expire_time"]) != "" {
		expirTime, errExpirTime := time.Parse("2006-01-02 15:04:05", convert.ToString(coupon["expire_time"]))
		if errExpirTime != nil {
			return utils.ErrorNull(c, "该现金券到期时间异常")
		}
		if !expirTime.Add(24 * time.Hour).After(time.Now()) {
			return utils.ErrorNull(c, "该现金券已到期")
		}
	}
	//开始使用
	x, err := global.DB.Update("UPDATE account_receive_coupon SET status=? WHERE id=?", enum.COUPON_USE, id)
	if err != nil {
		global.Log.Error("account_receive_coupon sql ERROR：", err.Error())
		return nil
	}
	if x <= 0 {
		return utils.ErrorNull(c, "使用现金券失败")
	}
	return utils.ErrorNull(c, "使用现金券成功")
}

func DeleteCoupon(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	id := c.FormValue("id")
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id格式错误")
	}
	coupon := GetCoupon(convert.MustInt64(id))
	if coupon == nil {
		return utils.ErrorNull(c, "该现金券不存在")
	}
	if convert.MustInt64(coupon["account_id"]) != acc.ID {
		return utils.ErrorNull(c, "删除失败，不是你发的现金券")
	}
	if convert.ToString(coupon["status"]) == enum.DELETE {
		return utils.ErrorNull(c, "已是删除状态")
	}
	x, err := global.DB.Update("UPDATE account_coupon SET status=? WHERE id=?", enum.DELETE, id)
	if err != nil {
		global.Log.Error("account_coupon sql ERROR：", err.Error())
		return nil
	}
	if x <= 0 {
		return utils.ErrorNull(c, "删除现金券失败")
	}
	return utils.SuccessNull(c, "删除现金券成功")
}

//现金券使用次数
func GetUseCouponNum(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	rows, err := global.DB.Query("SELECT count(1) as num FROM v_account_receive_coupon WHERE account_id=? AND receiveStatus=?", acc.ID, enum.COUPON_USE)
	if err != nil {
		global.Log.Error("v_account_receive_coupon sql error:%s", err.Error())
		return utils.ErrorNull(c, "数据异常，"+err.Error())
	}
	if len(rows) != 1 {
		return utils.SuccessNullMsg(c, 0)
	}
	num := convert.MustInt64(rows[0]["num"])
	if convert.MustInt64(rows[0]["num"]) <= 0 {
		return utils.SuccessNullMsg(c, 0)
	}
	return utils.SuccessNullMsg(c, num)
}
