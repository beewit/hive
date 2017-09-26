package handler

import (
	"encoding/json"
	"io/ioutil"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/hive/global"
	"github.com/labstack/echo"
)

func readBody(c echo.Context) (map[string]string, error) {
	body, bErr := ioutil.ReadAll(c.Request().Body)
	if bErr != nil {
		global.Log.Error("读取http body失败，原因：", bErr.Error())
		return nil, bErr
	}
	defer c.Request().Body.Close()

	var bm map[string]string
	bErr = json.Unmarshal(body, &bm)
	if bErr != nil {
		global.Log.Error("解析http body失败，原因：", bErr.Error())
		return nil, bErr
	}
	return bm, bErr
}

func Filter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token string
		bm, _ := readBody(c)
		if bm != nil {
			token = bm["token"]
		}
		if token == "" {
			return utils.AuthFail(c, "登陆信息token无效，请重新登陆")
		}

		accMapStr, err := global.RD.GetString(token)
		if err != nil {
			global.Log.Error(err.Error())
			return utils.AuthFail(c, "登陆信息已失效，请重新登陆")
		}
		if accMapStr == "" {
			global.Log.Error(token + "已失效")
			return utils.AuthFail(c, "登陆信息已失效，请重新登陆")
		}
		accMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(accMapStr), &accMap)
		if err != nil {
			global.Log.Error(accMapStr + "，error：" + err.Error())
			return utils.AuthFail(c, "登陆信息已失效，请重新登陆")
		}
		m, err := global.DB.Query("SELECT id,nickname,photo,mobile,status FROM account WHERE id=? LIMIT 1", accMap["id"])
		if err != nil {
			return utils.AuthFail(c, "获取用户信息失败")
		}
		if convert.ToString(m[0]["status"]) != enum.NORMAL {
			return utils.AuthFail(c, "用户已被冻结")
		}

		c.Set("account", global.ToMapAccount(m[0]))
		return next(c)
	}
}

func GetTemplateByListPage(c echo.Context) error {
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "article_template",
		Where:     "status = 1 ORDER BY `order` DESC,ct_time DESC",
		PageIndex: pageIndex,
		PageSize:  pageSize,
	})
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func UpdateTemplateReferById(c echo.Context) error {
	id := c.Param("id")
	if !utils.IsValidNumber(id) {
		return utils.Error(c, "id非法", nil)
	}
	x, err := global.DB.Update("UPDATE article_template SET refer_num=refer_num+1 WHERE id=?", id)
	if err != nil {
		return utils.Error(c, "更新引用数失败，"+err.Error(), nil)
	}
	if x > 0 {
		return utils.Success(c, "更新成功", nil)
	} else {
		return utils.Error(c, "更新引用数失败", nil)
	}
}

func GetTemplateByList(c echo.Context) error {
	sql := "SELECT * FROM article_template WHERE status = 1 ORDER BY `order` DESC,ct_time DESC"
	rows, err := global.DB.Query(sql)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if len(rows) <= 0 {
		return utils.Success(c, "无数据", nil)
	}
	return utils.Success(c, "获取数据成功", rows)
}

func GetTemplateById(c echo.Context) error {
	id := c.Param("id")
	if !utils.IsValidNumber(id) {
		return utils.Error(c, "id非法", nil)
	}
	sql := `SELECT * FROM article_template WHERE id=? AND status = 1 LIMIT 1`
	rows, err := global.DB.Query(sql, id)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if len(rows) != 1 {
		return utils.Success(c, "无数据", nil)
	}
	return utils.Success(c, "获取数据成功", rows[0])
}
