package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
)

/**
  获取账号对功能的所有试用记录
 */
func GetAccountFuncTryout(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	return utils.SuccessNullMsg(c, getListByAccId(acc.ID))
}

/**
获取账号功能的所有试用记录
 */
func getListByAccId(accId int64) []map[string]interface{} {
	rows, err := global.DB.Query("SELECT * FROM account_func_tryout WHERE account_id=?", accId)
	if err != nil {
		global.Log.Error("getListByAccId sql error:%s", err.Error())
		return nil
	}
	return rows
}

/**
获取功能操作试用记录
 */
func get(accId, funcId int64, funcHandleName string) map[string]interface{} {
	rows, err := global.DB.Query("SELECT * FROM account_func_tryout WHERE account_id=? AND func_id=? AND func_handle_name=? LIMIT 1",
		accId, funcId, funcHandleName)
	if err != nil {
		global.Log.Error("get sql error:%s", err.Error())
		return nil
	}
	if len(rows) != 1 {
		return nil
	}
	return rows[0]
}

/**
  获取功能的试用记录
 */
func AddAccountFuncTryout(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	funcId := c.FormValue("funcId")
	funcHandleName := c.FormValue("funcHandleName")
	if funcId == "" || funcHandleName == "" || !utils.IsValidNumber(funcId) {
		return utils.ErrorNull(c, "参数值无效")
	}
	nowTime := utils.CurrentTime()
	m := get(acc.ID, convert.MustInt64(funcId), funcHandleName)
	if m == nil {
		_, err = global.DB.InsertMap("account_func_tryout", map[string]interface{}{
			"id":               utils.ID(),
			"account_id":       acc.ID,
			"func_id":          funcId,
			"func_handle_name": funcHandleName,
			"count":            1,
			"ct_time":          nowTime,
			"ut_time":          nowTime,
			"ip":               c.RealIP(),
		})
	} else {
		_, err = global.DB.Update("UPDATE account_func_tryout SET count=count+1 WHERE account_id=? AND func_id=? AND func_handle_name=? LIMIT 1",
			acc.ID, funcId, funcHandleName)
	}
	if err != nil {
		global.Log.Error("AddAccountFuncTryout sql error:%s", err.Error())
		return utils.ErrorNull(c, "保存失败")
	}
	return utils.SuccessNull(c, "保存成功")
}
