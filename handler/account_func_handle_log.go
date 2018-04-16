package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
)

/**
	添加功能操作统计
 */
func AddAccountFuncHandleLog(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	funcId := c.FormValue("funcId")
	funcHandleFlag := c.FormValue("funcHandleFlag")
	funcHandleName := c.FormValue("funcHandleName")
	objName := c.FormValue("objName")
	content := c.FormValue("content")
	remark := c.FormValue("remark")
	if funcId == "" || funcHandleName == "" || (objName == "" && remark == "") || content == "" ||
		!utils.IsValidNumber(funcId) || !utils.IsValidNumber(funcHandleFlag) {
		return utils.ErrorNull(c, "参数值无效")
	}
	_, err = global.DB.InsertMap("account_func_handle_log", map[string]interface{}{
		"id":               utils.ID(),
		"account_id":       acc.ID,
		"account_func_id":  funcId,
		"func_handle_flag": funcHandleFlag,
		"func_handle_name": funcHandleName,
		"obj_name":         objName,
		"content":          content,
		"remark":          remark,
		"ct_time":          utils.CurrentTime(),
		"ip":               c.RealIP(),
	})
	if err != nil {
		global.Log.Error("AddAccountFuncHandleLog sql error:%s", err.Error())
		return utils.ErrorNull(c, "保存失败")
	}
	return utils.SuccessNull(c, "保存成功")
}

/**
	获取功能操作日志统计
 */
func GetAccountFuncHandleGroup(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	rows := getAccountFuncHandleRankingGroup()
	if rows != nil {
		var mobile string
		for i := 0; i < len(rows); i++ {
			mobile = convert.ToString(rows[i]["mobile"])
			rows[i]["mobile"] = utils.MobileReplaceRepl(mobile)
		}
	}
	return utils.SuccessNullMsg(c, map[string]interface{}{
		"group":        getAccountFuncHandleGroup(acc.ID),
		"rankingGroup": rows,
	})
}

/**
	近一个月内的账号操作发送记录统计
 */
func getAccountFuncHandleGroup(accId int64) []map[string]interface{} {
	rows, err := global.DB.Query("SELECT func_handle_name,count(func_handle_name) AS sum FROM account_func_handle_log"+
		" WHERE account_id=? AND DATE_SUB(CURDATE(), INTERVAL 30 DAY) <=date(ct_time) GROUP BY func_handle_name", accId)
	if err != nil {
		global.Log.Error("getAccountFuncHandleGroup sql error:%s", err.Error())
		return nil
	}
	return rows
}

/**
	近一个月内的账号操作发送记录统计
 */
func getAccountFuncHandleRankingGroup() []map[string]interface{} {
	rows, err := global.DB.Query("SELECT log.account_id,mobile,count(log.account_id) AS sum FROM account_func_handle_log log LEFT JOIN account ON account.id=log.account_id" +
		" WHERE DATE_SUB(CURDATE(), INTERVAL 30 DAY) <=date(log.ct_time) GROUP BY log.account_id ORDER BY sum DESC LIMIT 10")
	if err != nil {
		global.Log.Error("getAccountFuncHandleGroup sql error:%s", err.Error())
		return nil
	}
	return rows
}

/**
   根据标记获取功能操作的详细记录
 */
func GetAccountFuncHandleLogList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	funcHandleFlag := c.FormValue("funcHandleFlag")
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "account_func_handle_log",
		Where:     "account_id=? AND func_handle_flag=? AND DATE_SUB(CURDATE(), INTERVAL 30 DAY) <=date(log.ct_time)",
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Order:     "ct_time DESC",
	}, acc.ID, funcHandleFlag)
	if err != nil {
		global.Log.Error("QueryPage account_coupon sql error:%s", err.Error())
		return utils.ErrorNull(c, "数据异常")
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}
