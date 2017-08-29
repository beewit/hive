package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/convert"
)

func GetTemplateByList(c echo.Context) error {
	sql := `SELECT * FROM article_template WHERE status = 1`
	rows, err := global.DB.Query(sql)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if len(rows) <= 0 {
		return utils.Success(c, "无数据", nil)
	}
	for i := range rows {
		v := rows[i]
		m := convert.ToMapString(v)
		for k, v := range m {
			println(k, "   -    ", v)
		}
	}
	return utils.Success(c, "有数据", convert.ToArrayMapString(rows))
}

func GetTemplateById(c echo.Context) error {
	id := c.Param("id")
	if !utils.IsValidNumber(id) {
		return utils.Error(c, "id非法", nil)
	}
	sql := `SELECT * FROM article_template WHERE id=? AND status = 1`
	rows, err := global.DB.Query(sql, id)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if len(rows) != 1 {
		return utils.Success(c, "无数据", nil)
	}
	return utils.Success(c, "有数据", convert.ToMapString(rows[0]))
}
