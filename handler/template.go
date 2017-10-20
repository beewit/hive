package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/labstack/echo"
)

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
