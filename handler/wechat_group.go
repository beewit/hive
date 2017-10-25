package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/enum"
	"fmt"
	"github.com/beewit/beekit/utils/convert"
)

func GetWechatGroupList(c echo.Context) error {
	itf := c.Get("account")
	if itf == nil {
		return utils.AuthFailNull(c)
	}
	acc := global.ToInterfaceAccount(itf)
	if acc == nil {
		return utils.AuthFailNull(c)
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	where := "status=? AND id NOT IN(SELECT wechat_group_id FROM account_wechat_group WHERE account_id=?)"
	area := c.FormValue("area")
	types := c.FormValue("type")
	if area != "" {
		where += fmt.Sprintf(" AND area='%s'", area)
	}
	if types != "" {
		where += fmt.Sprintf(" AND type='%s'", types)
	}

	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "wechat_group",
		Where:     where,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Order:     "ut_time DESC",
	}, enum.NORMAL, acc.ID)
	if err != nil {
		return utils.ErrorNull(c, "数据异常，"+err.Error())
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func GetWechatGroupClass(c echo.Context) error {
	typeMap := GetWechatGroupType()
	areaMap := GetWechatGroupArea()
	m := map[string]interface{}{}
	m["area"] = areaMap
	m["type"] = typeMap
	return utils.SuccessNullMsg(c, m)
}

func GetWechatGroupType() []map[string]interface{} {
	//缓存
	reidsWechatGroupType := "WechatGroupType"
	wgt, err := global.RD.GetString(reidsWechatGroupType)
	if err != nil {
		global.Log.Error(err.Error())
	} else {
		if wgt != "" {
			redisMap, err := convert.String2MapList(wgt)
			if err != nil {
				global.Log.Error(err.Error())
			} else {
				if redisMap != nil && len(redisMap) > 0 {
					return redisMap
				}
			}
		}
	}
	m, err := global.DB.Query("SELECT type ,count(type) as sum FROM wechat_group WHERE type!='' GROUP BY type")
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if m != nil && len(m) > 0 {
		val, err := convert.ToArrayMapStr(m)
		if err != nil {
			global.Log.Error(err.Error())
		} else {
			//十分钟缓存
			global.RD.SetAndExpire(reidsWechatGroupType, val, 60*10)
		}
	}
	return m
}

func GetWechatGroupArea() []map[string]interface{} {
	//缓存
	reidsWechatGroupArea := "WechatGroupArea"
	wgt, err := global.RD.GetString(reidsWechatGroupArea)
	if err != nil {
		global.Log.Error(err.Error())
	} else {
		if wgt != "" {
			redisMap, err := convert.String2MapList(wgt)
			if err != nil {
				global.Log.Error(err.Error())
			} else {
				if redisMap != nil && len(redisMap) > 0 {
					return redisMap
				}
			}
		}
	}
	m, err := global.DB.Query("SELECT area ,count(area) as sum FROM wechat_group WHERE area!='' GROUP BY area")
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if m != nil && len(m) > 0 {
		val, err := convert.ToArrayMapStr(m)
		if err != nil {
			global.Log.Error(err.Error())
		} else {
			//十分钟缓存
			global.RD.SetAndExpire(reidsWechatGroupArea, val, 60*10)
		}
	}
	return m
}

func AddAccountWechatGroup(c echo.Context) error {
	itf := c.Get("account")
	if itf == nil {
		return utils.AuthFailNull(c)
	}
	acc := global.ToInterfaceAccount(itf)
	if acc == nil {
		return utils.AuthFailNull(c)
	}
	wgId := c.FormValue("wgId")
	if !utils.IsValidNumber(wgId) {
		utils.ErrorNull(c, "wgId格式错误")
	}
	if getWechatGroupByid(convert.MustInt64(wgId)) == nil {
		utils.ErrorNull(c, "无此wgId数据")
	}
	iw, _ := utils.NewIdWorker(1)
	id, _ := iw.NextId()
	m := map[string]interface{}{}
	m["id"] = id
	m["account_id"] = acc.ID
	m["wechat_group_id"] = wgId
	m["ct_time"] = utils.CurrentTime()
	_, err := global.DB.InsertMap("account_wechat_group", m)
	if err != nil {
		return utils.ErrorNull(c, "数据异常，"+err.Error())
	}
	return utils.SuccessNull(c, "保存数据成功")
}

func getWechatGroupByid(id int64) map[string]interface{} {
	m, err := global.DB.Query("SELECT * FROM account_wechat_group WHERE id=? LIMIT 1", id)
	if err != nil || len(m) != 1 {
		return nil
	}
	return m[0]
}
