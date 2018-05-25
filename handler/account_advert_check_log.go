package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/convert"
)

/**
	添加发送的广告检查结果日志
 */
func AddAccountAdvertCheckLog(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	var advertId = c.FormValue("advertId")
	var content = c.FormValue("content")
	if content == "" {
		return utils.ErrorNull(c, "检查内容为空")
	}
	if advertId == "" || !utils.IsValidNumber(advertId) {
		return utils.ErrorNull(c, "未能获得检查的广告")
	}
	m := GetAccountAdvertCheckLog(convert.MustInt64(advertId), acc.ID)
	if m != nil {
		return utils.ErrorNull(c, "已提交记录")
	}
	id := utils.ID()
	nowTime := utils.CurrentTime()
	_, err = global.DB.InsertMap("account_advert_check_log", map[string]interface{}{
		"id":                id,
		"account_advert_id": advertId,
		"check_time":        nowTime,
		"ip":                c.RealIP(),
		"account_id":        acc.ID,
	})
	if err != nil {
		global.Log.Error("AddAccountAdvertCheckLog sql error：%s", err.Error())
		return utils.ErrorNull(c, "保存失败")
	}
	return utils.SuccessNull(c, "保存成功")
}

func GetAccountAdvertCheckLog(advertId, accId int64) map[string]interface{} {
	m, err := global.DB.Query("SELECT * FROM account_advert_check_log WHERE account_advert_id=? AND account_id=? LIMIT 1", advertId, accId)
	if err != nil {
		global.Log.Error("GetAccountAdvertCheckLog sql error：%s", err.Error())
		return nil
	}
	if len(m) != 1 {
		return nil
	}
	return m[0]
}
