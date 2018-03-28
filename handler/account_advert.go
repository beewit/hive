package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/enum"
	"strings"
	"github.com/beewit/beekit/utils/convert"
)

/**
	账号设置的广告
 */
func GetAccountAdvertList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "account_advert",
		Where:     "status=? AND account_id=?",
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Order:     "ct_time DESC",
	}, enum.NORMAL, acc.ID)
	if err != nil {
		global.Log.Error("GetAccountAdvertList sql error：%s", err.Error())
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func getAccountAdvert(id int64) map[string]interface{} {
	if id < 1 {
		return nil
	}
	rows, err := global.DB.Query("SELECT * FROM account_advert WHERE id=? LIMIT 1", id)
	if err != nil {
		global.Log.Error("getAccountAdvert sql error：%s", err.Error())
		return nil
	}
	if len(rows) != 1 {
		return nil
	}
	return rows[0]
}

/**
	添加账号广告
 */
func AddAccountAdvert(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	id := strings.TrimSpace(c.FormValue("id"))
	t := strings.TrimSpace(c.FormValue("type"))
	applyType := strings.TrimSpace(c.FormValue("applyType"))
	title := strings.TrimSpace(c.FormValue("title"))
	content := strings.TrimSpace(c.FormValue("content"))
	img := strings.TrimSpace(c.FormValue("img"))
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
	if id != "" && utils.IsValidNumber(id) {
		m := getAccountAdvert(convert.MustInt64(id))
		if m == nil {
			return utils.ErrorNull(c, "修改失败")
		}
		if convert.ToString(m["status"]) != enum.NORMAL {
			return utils.ErrorNull(c, "已被删除")
		}
		if convert.MustInt64(m["account_id"]) != acc.ID {
			return utils.ErrorNull(c, "无权限修改")
		}
		_, err = global.DB.Update("UPDATE account_advert SET type=?,title=?,content=?,img=?,ut_time=?,ip=? WHERE id=?",
			t, title, content, img, nowTime, ip, id)
	} else {
		_, err = global.DB.InsertMap("account_advert", map[string]interface{}{
			"id":         utils.ID(),
			"account_id": acc.ID,
			"type":       t,
			"title":      title,
			"content":    content,
			"img":        img,
			"status":     enum.NORMAL,
			"ct_time":    nowTime,
			"ut_time":    nowTime,
			"ip":         ip,
		})
	}
	if err != nil {
		global.Log.Error("AddAccountAdvert sql error：%s", err.Error())
		return utils.ErrorNull(c, "保存失败")
	}
	return utils.SuccessNull(c, "保存成功")
}

/**
	删除个人广告
 */
func DelAccountAdvert(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	id := strings.TrimSpace(c.FormValue("id"))
	if id != "" && utils.IsValidNumber(id) {
		m := getAccountAdvert(convert.MustInt64(id))
		if m == nil {
			return utils.ErrorNull(c, "修改失败")
		}
		if convert.ToString(m["status"]) != enum.NORMAL {
			return utils.ErrorNull(c, "已被删除")
		}
		if convert.MustInt64(m["account_id"]) != acc.ID {
			return utils.ErrorNull(c, "无权限修改")
		}
		_, err = global.DB.Update("UPDATE account_advert SET status=?,ut_time=?,ip=? WHERE id=?",
			enum.DELETE, utils.CurrentTime(), c.RealIP(), id)
		if err != nil {
			global.Log.Error("DelAccountAdvert sql error：%s", err.Error())
			return utils.ErrorNull(c, "删除失败")
		}
		return utils.SuccessNull(c, "删除成功")
	} else {
		return utils.ErrorNull(c, "删除失败")
	}
}


/**
	删除个人广告
 */
func GetAccountAdvert(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	id := strings.TrimSpace(c.FormValue("id"))
	if id != "" && utils.IsValidNumber(id) {
		m := getAccountAdvert(convert.MustInt64(id))
		if m == nil {
			return utils.ErrorNull(c, "未获取到数据")
		}
		if convert.ToString(m["status"]) != enum.NORMAL {
			return utils.ErrorNull(c, "已被删除")
		}
		if convert.MustInt64(m["account_id"]) != acc.ID {
			return utils.ErrorNull(c, "无权限修改")
		}
		return utils.SuccessNullMsg(c, m)
	} else {
		return utils.ErrorNull(c, "未获取到数据")
	}
}
