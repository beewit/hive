package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/enum"
	"strings"
)

/**
	账号设置的广告
 */
func GetAccountAdvertList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	rows, err := global.DB.Query("SELECT * FROM account_advert WHERE account_id=? AND status=?", acc.ID, enum.NORMAL)
	if err != nil {
		global.Log.Error("GetAccountAdvertList sql error：%s", err.Error())
		return utils.ErrorNull(c, "获取设置的广告失败")
	}
	if len(rows) == 0 {
		return utils.NullData(c)
	}
	return utils.SuccessNullMsg(c, rows)
}

/**
	添加账号广告
 */
func AddAccountAdvert(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	t := strings.TrimSpace(c.FormValue("type"))
	applyType := strings.TrimSpace(c.FormValue("applyType"))
	title := strings.TrimSpace(c.FormValue("title"))
	content := strings.TrimSpace(c.FormValue("content"))
	img := strings.TrimSpace(c.FormValue("img"))
	status := strings.TrimSpace(c.FormValue("status"))
	ip := c.RealIP()
	switch t {
	case enum.ACCOUNT_ADVERT_IMG:
		if img == "" {
			return utils.ErrorNull(c, "请选择图片")
		}
		break
	case enum.ACCOUNT_ADVERT_TEXT:
		if content == "" {
			return utils.ErrorNull(c, "请输入内容")
		}
		break
	case enum.ACCOUNT_ADVERT_IMG_TEXT:
		if img == "" {
			return utils.ErrorNull(c, "请选择图片")
		}
		if content == "" {
			return utils.ErrorNull(c, "请输入内容")
		}
		break
	default:
		return utils.ErrorNull(c, "类型不存在")
	}
	switch applyType {
	case enum.ACCOUNT_ADVERT_APPLY_ALL:
	case enum.ACCOUNT_ADVERT_APPLY_FISSION:
	case enum.ACCOUNT_ADVERT_APPLY_AUTO:
	case enum.ACCOUNT_ADVERT_APPLY_SMS:
		break
	case enum.ACCOUNT_ADVERT_APPLY_EMAIL:
		if title == "" {
			return utils.ErrorNull(c, "邮件需要添加标题")
		}
		break
	default:
		return utils.ErrorNull(c, "适用类型不存在")
	}

	switch status {
	case enum.NORMAL:
	case enum.DELETE:
		break
	default:
		return utils.ErrorNull(c, "操作的数据状态错误")
		break
	}

	if len(title) > 255 {
		return utils.ErrorNull(c, "标题长度限定255个字符")
	}

	if len(content) > 65535 {
		return utils.ErrorNull(c, "内容过长")
	}
	if len(img) > 1500 {
		return utils.ErrorNull(c, "图片超出限定数量")
	}
	if len(strings.Split(img, ",")) > 9 {
		return utils.ErrorNull(c, "图片超出限定数量")
	}
	nowTime := utils.CurrentTime()
	_, err = global.DB.InsertMap("account_advert", map[string]interface{}{
		"id":         utils.ID(),
		"account_id": acc.ID,
		"type":       t,
		"title":      title,
		"content":    content,
		"img":        img,
		"status":     status,
		"ct_time":    nowTime,
		"ut_time":    nowTime,
		"ip":         ip,
	})
	if err != nil {
		global.Log.Error("AddAccountAdvert sql error：%s", err.Error())
		return utils.ErrorNull(c, "保存失败")
	}
	return utils.SuccessNull(c, "保存成功")
}
