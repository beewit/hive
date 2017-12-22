package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/hive/global"
	"github.com/labstack/echo"
	"fmt"
	"github.com/beewit/beekit/utils/convert"
	"time"
)

func GetFuncList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "f.*,af.expiration_time",
		Table:     "func f LEFT JOIN account_func af ON af.func_id=f.id AND af.account_id=?",
		Where:     "f.status=?",
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, acc.ID, enum.NORMAL)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func GetAccountFuncList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "f.*,af.expiration_time",
		Table:     "account_func af LEFT JOIN func f ON af.func_id=f.id",
		Where:     "f.status=? AND af.account_id=?",
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, enum.NORMAL, acc.ID)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func GetFuncGiveLog(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "log.*,f.name",
		Table:     "account_func_give_log log LEFT JOIN func f ON log.func_id=f.id",
		Where:     "log.account_id=?",
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Order:     "log.ct_time DESC",
	},  acc.ID)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func GetFuncAllByIdsAndAccId(c echo.Context) error {
	accID := c.FormValue("accId")
	funcIds := c.FormValue("funcIds")
	platformIds := c.FormValue("platformIds")
	where := "af.expiration_time>now() AND f.status=? AND af.account_id=?"
	if funcIds != "" {
		where += fmt.Sprintf(" AND f.id in(%s)", funcIds)
	}
	if platformIds != "" {
		where += fmt.Sprintf(" AND f.platform_id in(%s)", platformIds)
	}
	sql := fmt.Sprintf(
		"SELECT f.*,af.expiration_time FROM account_func af LEFT JOIN func f ON af.func_id=f.id LEFT JOIN platform p ON p.id=f.platform_id WHERE %s",
		where)
	m, err := global.DB.Query(sql, enum.NORMAL, accID)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	return utils.Success(c, "获取数据成功", m)
}

func GetEffectiveFuncById(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	funcId := c.FormValue("funcId")
	if funcId == "" {
		return utils.ErrorNull(c, "funcId不能为空")
	}
	if !utils.IsValidNumber(funcId) {
		return utils.ErrorNull(c, "funcId参数错误")
	}
	sql := "SELECT f.*,af.expiration_time FROM account_func af LEFT JOIN func f ON af.func_id=f.id " +
		"LEFT JOIN platform p ON p.id=f.platform_id " +
		"WHERE af.expiration_time>now() AND f.status=? AND af.account_id=? AND f.id=? LIMIT 1"
	m, err := global.DB.Query(sql, enum.NORMAL, acc.ID, funcId)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if len(m) != 1 {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", m[0])
}

func GetEffectiveFuncList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	sql := "SELECT f.*,af.expiration_time FROM account_func af LEFT JOIN func f ON af.func_id=f.id " +
		"LEFT JOIN platform p ON p.id=f.platform_id " +
		"WHERE af.expiration_time>now() AND f.status=? AND af.account_id=? "
	m, err := global.DB.Query(sql, enum.NORMAL, acc.ID)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if len(m) <= 0 {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", m)
}

func GetFuncGroupByAccountId(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	openSql := "SELECT * FROM account_func WHERE expiration_time>NOW() AND account_id=?"
	expireSql := "SELECT * FROM account_func WHERE expiration_time<NOW() AND account_id=?"
	notExpireSql := "SELECT * FROM account_func WHERE expiration_time>date_add(NOW(),interval 15 DAY) AND account_id=?"
	openMap, _ := global.DB.Query(openSql, acc.ID)
	expireMap, _ := global.DB.Query(expireSql, acc.ID)
	notExpireMap, _ := global.DB.Query(notExpireSql, acc.ID)
	data := map[string]interface{}{
		"openMap":      openMap,
		"expireMap":    expireMap,
		"notExpireMap": notExpireMap,
	}
	return utils.SuccessNullMsg(c, data)
}

//该帐号开通的功能项目
func getAccountFuncByAccId(accId int64) []map[string]interface{} {
	sql := "SELECT * FROM account_func WHERE account_id=?"
	rows, err := global.DB.Query(sql, accId)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if len(rows) < 1 {
		return nil
	}
	return rows
}

//该帐号开通的功能项目
func getAccountFuncByAccIdOrFuncId(accId, funcId int64) map[string]interface{} {
	sql := "SELECT * FROM account_func WHERE account_id=? AND func_id=? LIMIT 1"
	rows, err := global.DB.Query(sql, accId, funcId)
	if err != nil {
		global.Log.Error(err.Error())
		return nil
	}
	if len(rows) < 1 {
		return nil
	}
	return rows[0]
}

//根据分享链接获得添加微信助手的使用时间
func AddShareWechatAppTime(c echo.Context) error {
	//查询当天分享次数，超过5次即可获得一天的使用权
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	remarks := c.FormValue("remarks")
	sql := "SELECT COUNT(1) as sum FROM account_action_logs WHERE type='func' AND account_id=? AND remarks='分享APP' AND to_days(ct_time) = to_days(now());"
	m, _ := global.DB.Query(sql, acc.ID)
	if m == nil {
		return utils.ErrorNull(c, "查询使用次数失败")
	}
	hour := 8
	sum := convert.MustInt(m[0]["sum"])
	if sum >= 5 {
		//判断是否今日赠送过了时间，如果已经赠送不再赠送了
		sql = "SELECT COUNT(1) as sum FROM account_func_give_log WHERE type='SHARE_APP' AND account_id=? AND func_id=? AND to_days(ct_time) = to_days(now());"
		m, _ = global.DB.Query(sql, acc.ID,enum.FUNC_WECHAT_APP)
		if m == nil {
			return utils.ErrorNull(c, "查询赠送功能使用时间记录失败")
		}
		if convert.MustInt(m[0]["sum"]) > 0 {
			return utils.SuccessNull(c, "今日已经获得赠送功能使用时间")
		}
		err = addFuncGiveLog(acc, "SHARE_APP", remarks, enum.FUNC_WECHAT_APP, hour)
		if err != nil {
			return utils.ErrorNull(c, err.Error())
		}
		return utils.SuccessNull(c, fmt.Sprintf("恭喜您获得%d小时微信营销助手使用时间", hour))
	}
	return utils.SuccessNull(c, fmt.Sprintf("每天分享满5次可获得%d小时微信营销助手使用时间", hour))
}

func addFuncGiveLog(acc *global.Account, t, remarks string, funcId int64, hour int) (err error) {
	err = AddFuncTime(acc, funcId, hour)
	if err != nil {
		return
	}
	giveMap := map[string]interface{}{}
	giveMap["id"] = utils.ID()
	giveMap["account_id"] = acc.ID
	giveMap["type"] = t
	giveMap["func_id"] = funcId
	giveMap["hour"] = hour
	giveMap["remarks"] = remarks
	giveMap["ct_time"] = utils.CurrentTime()
	_, err = global.DB.InsertMap("account_func_give_log", giveMap)
	return
}

func AddFuncTime(acc *global.Account, funcId int64, hour int) (err error) {
	var flog bool
	var daysTime time.Time
	h := time.Hour * time.Duration(hour)
	accFunc := getAccountFuncByAccIdOrFuncId(acc.ID, funcId)
	if accFunc != nil && accFunc["expiration_time"] != nil {
		flog = true
		expirTimeStr, errExpirTime := time.Parse("2006-01-02 15:04:05", convert.ToString(accFunc["expiration_time"]))
		if errExpirTime != nil {
			global.Log.Error(convert.ToString(acc.ID) + "会员的过期时间错误：" + errExpirTime.Error())
			return
		}
		if expirTimeStr.After(time.Now()) {
			//未到期的续费
			daysTime = expirTimeStr.Add(h)
		} else {
			//已到期的续费
			daysTime = time.Now().Add(h)
		}
	} else {
		flog = false
		//无到期时间
		daysTime = time.Now().Add(h)
	}
	if flog {
		//修改
		sql := "UPDATE account_func SET expiration_time=?,ut_time=? WHERE account_id=? AND func_id=?"
		_, err = global.DB.Update(sql, utils.FormatTime(daysTime), utils.CurrentTime(), acc.ID, funcId)
		if err != nil {
			global.Log.Error(err.Error())
			return
		}
	} else {
		//添加
		daysTime = time.Now().Add(h)
		m := make(map[string]interface{})
		m["id"] = utils.ID()
		m["account_id"] = acc.ID
		m["func_id"] = funcId
		m["expiration_time"] = utils.FormatTime(daysTime)
		m["ct_time"] = utils.CurrentTime()
		m["ut_time"] = utils.CurrentTime()
		_, err = global.DB.InsertMap("account_func", m)
		if err != nil {
			global.Log.Error(err.Error())
			return
		}
	}
	return
}
