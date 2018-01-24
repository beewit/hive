/*
领取代金券记录
author：zxb
2018-01-13
*/

package handler

import (
	"fmt"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/hive/global"
	"github.com/labstack/echo"
	"time"
)

func ReceiveCoupon(c echo.Context) error {
	ws, err := GetMiniAppSession(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	id := c.FormValue("id")
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id格式错误")
	}
	coupon := GetCoupon(convert.MustInt64(id))
	if coupon == nil {
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
		if !expirTime.Add(24 * time.Hour).After(time.Now()) {
			return utils.ErrorNull(c, "该代金券已到期无法领取")
		}
	}
	couponId := convert.MustInt64(coupon["id"])
	receiveCoupon := GetReceiveCouponByCouponId(couponId, ws.Unionid)
	if receiveCoupon != nil {
		return utils.ErrorNull(c, "已领取过了")
	}

	acc := GetAccountByUnionId(ws.Unionid, enum.WECHAT)
	var accId interface{}
	if acc != nil {
		accId = acc["id"]
	} else {
		accId = nil
	}
	flog := false
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		_, err = tx.InsertMap("account_receive_coupon", map[string]interface{}{
			"id":                utils.ID(),
			"account_id":        accId,
			"account_coupon_id": couponId,
			"ct_time":           utils.CurrentTime(),
			"ip":                c.RealIP(),
			"status":            "未使用",
			"wx_union_id":       ws.Unionid,
		})
		if err != nil {
			global.Log.Error("insertMap account_receive_coupon sql error:%s", err.Error())
			panic(err)
		}
		x, err := tx.Update("UPDATE account_coupon SET receive_number=receive_number+1 WHERE id=? AND receive_number<number", couponId)
		if err != nil {
			global.Log.Error("update account_receive_coupon sql error:%s", err.Error())
			panic(err)
		}
		if x <= 0 {
			panic(err)
		}
		flog = true
	}, func(err error) {
		if err != nil {
			global.Log.Error("领取代金券失败：%v", err)
			flog = false
		}
	})
	if !flog {
		return utils.ErrorNull(c, "领取代金券失败")
	}
	return utils.SuccessNull(c, "领取代金券成功")
}

func GetReceiveCoupon(id int64) map[string]interface{} {
	m, err := global.DB.Query("SELECT * FROM account_receive_coupon WHERE id=? LIMIT 1", id)
	if err != nil {
		global.Log.Error("getCouponByID sql ERROR：", err.Error())
		return nil
	}
	if len(m) != 1 {
		return nil
	}
	return m[0]
}

func GetReceiveCouponByCouponId(id int64, unionId string) map[string]interface{} {
	m, err := global.DB.Query("SELECT * FROM account_receive_coupon WHERE account_coupon_id=? AND wx_union_id=? LIMIT 1", id, unionId)
	if err != nil {
		global.Log.Error("getCouponByID sql ERROR：", err.Error())
		return nil
	}
	if len(m) != 1 {
		return nil
	}
	return m[0]
}

//领取代金券分页
func GetReceiveCouponList(c echo.Context) error {
	ws, err := GetMiniAppSession(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	t := c.FormValue("t")
	//未使用和未过期
	where := "receiveStatus='未使用' AND (expire_time is NULL OR expire_time>=DATE_FORMAT(now(),'%y-%m-%d'))"
	switch t {
	case "expire":
		//已过期
		where = "expire_time<DATE_FORMAT(now(),'%y-%m-%d')"
		break
	case "finish":
		//已使用
		where = "receiveStatus='已使用'"
		break
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "v_account_receive_coupon",
		Where:     "wx_union_id=? AND status=? AND " + where,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Order:     "receiveTime DESC",
	}, ws.Unionid, enum.NORMAL)
	if err != nil {
		global.Log.Error("QueryPage v_account_receive_coupon sql error:%s", err.Error())
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func CreateCouponQrCode(c echo.Context) error {
	ws, err := GetMiniAppSession(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	id := c.FormValue("id")
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id格式错误")
	}
	rows, err := global.DB.Query("SELECT * FROM account_receive_coupon WHERE id=? AND wx_union_id=? LIMIT 1", id, ws.Unionid)
	if err != nil {
		global.Log.Error("account_receive_coupon sql error:%s", err.Error())
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if len(rows) != 1 {
		return utils.Error(c, "该代金券不存在", nil)
	}
	qrcode := convert.ToString(rows[0]["qrcode"])
	if qrcode != "" {
		return utils.Success(c, "获取使用代金券二维码成功", qrcode)
	}
	fileName := utils.ID()
	//生成二维码
	buf, err := utils.CreateQrCodeBytes(fmt.Sprintf("%s|%s", id, convert.ToString(time.Now().UnixNano())))
	if err != nil {
		return utils.ResultHtml(c, "生成二维码失败")
	}

	filePath := fmt.Sprintf("%sqrcode/%s/%d.jpg", global.FileConf.Path, utils.CurrentDateByPlace("/"), fileName)
	dst, err := utils.CreateFile(filePath)
	if err != nil {
		return utils.ErrorNull(c, "创建文件失败")
	}
	defer dst.Close()

	path := global.GetSavePath(filePath)
	name := fmt.Sprintf("%d.jpg", fileName)
	b := buf.Bytes()
	// Copy
	if _, err = dst.Write(b); err != nil {
		return utils.ErrorNull(c, "保存文件失败")
	}
	go func() {
		//更新二维码
		_, err = global.DB.Update("UPDATE account_receive_coupon SET qrcode=? WHERE id=?", path, id)
		if err != nil {
			global.Log.Error("UPDATE account_receive_coupon，ERROR：%s", err.Error())
		}

		_, err = global.DB.InsertMap("file_log", map[string]interface{}{
			"id":      utils.ID(),
			"name":    name,
			"path":    path,
			"size":    len(b),
			"ext":     ".jpg",
			"ct_time": utils.CurrentTime(),
			"ct_ip":   c.RealIP(),
		})
		if err != nil {
			global.Log.Error("保存上传文件日志失败，ERROR：%s", err.Error())
		}
	}()
	return utils.Success(c, "获取使用代金券二维码成功", path)
}
