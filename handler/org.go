package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/beekit/utils/convert"
)

/**
   获取子集组织数据
 */
func GetOrgChild(c echo.Context) error {
	var parentId = c.FormValue("parentId")
	if parentId == "" {
		parentId = "0"
	}
	if !utils.IsValidNumber(parentId) {
		return utils.ErrorNull(c, "parentId错误")
	}
	rows, err := global.DB.Query("SELECT * FROM org WHERE parent_id=? AND status=?", parentId, enum.NORMAL)
	if err != nil {
		global.Log.Error("DelOrg GetOrgChild error：%s", err.Error())
		return utils.ErrorNull(c, "获取组织结构失败")
	}
	return utils.SuccessNullMsg(c, rows)
}

/**
   获取下级所有组织数据
 */
func GetOrgChildAll(c echo.Context) error {
	var parentId = c.FormValue("parentId")
	if parentId == "" {
		parentId = "0"
	}
	if !utils.IsValidNumber(parentId) {
		return utils.ErrorNull(c, "parentId错误")
	}
	m := GetOrgById(convert.MustInt64(parentId))
	if m == nil {
		return utils.ErrorNull(c, "查询的组织不存在")
	}
	//根据关系查询
	rows, err := global.DB.Query("SELECT * FROM org WHERE  Concat('-',relation,'-') like '-?-%' AND status=?", m["relation"], enum.NORMAL)
	if err != nil {
		global.Log.Error("DelOrg GetOrgAllChild error：%s", err.Error())
		return utils.ErrorNull(c, "获取组织结构失败")
	}
	return utils.SuccessNullMsg(c, rows)
}

/**
  获取组织结构
 */
func GetOrg(c echo.Context) error {
	var id = c.FormValue("id")
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id错误")
	}
	m := GetOrgById(convert.MustInt64(id))
	return utils.SuccessNullMsg(c, m)
}

func GetOrgById(id int64) map[string]interface{} {
	rows, err := global.DB.Query("SELECT * FROM org WHERE id=? AND status=? LIMIT 1", id, enum.NORMAL)
	if err != nil {
		global.Log.Error("GetOrg sql error：%s", err.Error())
		return nil
	}
	if len(rows) != 1 {
		return nil
	}
	return rows[0]
}

/**
	删除组织结构
 */
func DelOrg(c echo.Context) error {
	var id = c.FormValue("id")
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id错误")
	}
	x, err := global.DB.Update("UPDATE org SET status=? WHERE id=?", enum.DELETE, id)
	if err != nil {
		global.Log.Error("DelOrg sql error：%s", err.Error())
		return utils.ErrorNull(c, "删除失败")
	}
	if x > 0 {
		return utils.SuccessNull(c, "删除成功")
	}
	return utils.ErrorNull(c, "删除失败")
}

func AddOrg(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	var name = c.FormValue("name")
	var parentId = c.FormValue("parentId")
	if name == "" {
		return utils.ErrorNull(c, "请输入组织名称")
	}
	if len(name) > 255 {
		return utils.ErrorNull(c, "组织名称长度过长")
	}
	if parentId == "" {
		parentId = "0"
	}
	if !utils.IsValidNumber(parentId) {
		return utils.ErrorNull(c, "父级组织参数错误")
	}
	relation := "0"
	if parentId != "0" {
		org := GetOrgById(convert.MustInt64(parentId))
		if org == nil {
			return utils.ErrorNull(c, "父级组织不存在")
		}
		relation = convert.ToString(org["relation"])
	}
	id := utils.ID()
	nowTime := utils.CurrentTime()
	_, err = global.DB.InsertMap("org", map[string]interface{}{
		"id":         id,
		"name":       name,
		"parent_id":  parentId,
		"relation":   relation + "-" + convert.ToString(id),
		"ct_time":    nowTime,
		"ut_time":    nowTime,
		"ip":         c.RealIP(),
		"status":     enum.NORMAL,
		"account_id": acc.ID,
	})
	if err != nil {
		return utils.ErrorNull(c, "保存失败")
	}
	return utils.SuccessNull(c, "保存成功")
}
