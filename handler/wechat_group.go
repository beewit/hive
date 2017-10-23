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
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	where := "status=?"
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
	}, enum.NORMAL)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
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
 return 	utils.SuccessNullMsg(c, m)
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
	m, err := global.DB.Query("SELECT type ,count(type) as sum FROM wechat_group GROUP BY type")
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
	m, err := global.DB.Query("SELECT area ,count(area) as sum FROM wechat_group GROUP BY area")
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
