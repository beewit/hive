package handler

import (
	"fmt"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/encrypt"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/hive/global"
	"github.com/labstack/echo"
)

func GetWallet(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	m := getWallet(acc)
	if m == nil {
		//初始化钱包
		m = map[string]interface{}{}
		m["id"] = utils.ID()
		m["account_id"] = acc.ID
		m["money"] = 0
		m["ct_time"] = utils.CurrentTime()
		m["ct_ip"] = c.RealIP()
		_, err := global.DB.InsertMap("account_wallet", m)
		if err != nil {
			global.Log.Error("初始化钱包失败，ERROR：", err.Error())
		}
	}
	return utils.SuccessNullMsg(c, map[string]interface{}{
		"wallet":                 m,
		"unWithdrawCashMoney":    getUnWithdrawCashMoney(acc),
		"applyWithdrawCashMoney": getApplyWithdrawCashMoney(acc),
	})
}

func getWallet(acc *global.Account) map[string]interface{} {
	sql := "SELECT * FROM account_wallet WHERE account_id=? LIMIT 1"
	rows, _ := global.DB.Query(sql, acc.ID)
	if rows == nil || len(rows) != 1 {
		return nil
	}
	return rows[0]
}

/**
不能提现的金额，邀请返利获得的奖励金额需要一个月后才可申请提现
*/
func getUnWithdrawCashMoney(acc *global.Account) float64 {
	sql := "SELECT sum(change_money) as money FROM account_wallet_log WHERE DATE_SUB(CURDATE(), INTERVAL 1 MONTH) <= date(ct_time) AND account_id=? AND type=?"
	rows, _ := global.DB.Query(sql, acc.ID, enum.WALLET_REBATE)
	if rows == nil || len(rows) != 1 {
		return 0
	}
	return convert.MustFloat64(rows[0]["money"])
}

func AddWithdrawCashCard(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	cp := c.FormValue("collectPlatform")
	cn := c.FormValue("cardNumber")
	cname := c.FormValue("cardName")
	cq := c.FormValue("cardQrcode")
	if cp == "" {
		return utils.ErrorNull(c, "请选择提现的收款第三方平台")
	}
	if cname == "" {
		return utils.ErrorNull(c, "请输入收款人真实姓名")
	}
	if cn == "" {
		return utils.ErrorNull(c, "请输入第三方平台的账号")
	}
	if cq == "" {
		return utils.ErrorNull(c, "请上传收款二维码")
	}
	wallet := getWallet(acc)
	if wallet == nil {
		sql := "INSERT INTO account_wallet(id,account_id,wc_collect_platform,wc_card_number,wc_card_qrcode,wc_card_name,last_time,last_ip)VALUES(?,?,?,?,?,?,?,?)"
		_, err = global.DB.Insert(sql, utils.ID(), acc.ID, cp, cn, cq, cname, utils.CurrentTime(), c.RealIP())
		if err != nil {
			tip := "保存提现的收款信息失败"
			global.Log.Error(fmt.Sprintf("%s，错误：%s", tip, err.Error()))
			return utils.ErrorNull(c, tip)
		}
	} else {
		sql := "UPDATE account_wallet SET wc_collect_platform=?,wc_card_number=?,wc_card_qrcode=?,wc_card_name=?,last_time=?,last_ip=? WHERE account_id=?"
		_, err = global.DB.Insert(sql, cp, cn, cq, cname, utils.CurrentTime(), c.RealIP(), acc.ID)
		if err != nil {
			tip := "保存提现的收款信息失败"
			global.Log.Error(fmt.Sprintf("%s，错误：%s", tip, err.Error()))
			return utils.ErrorNull(c, tip)
		}
	}
	return utils.SuccessNull(c, "保存成功")
}

func GetWalletLogList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "account_wallet_log",
		Where:     "account_id=?",
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Order:     "ct_time DESC",
	}, acc.ID)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func ApplyWithdrawCash(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	applyMoney := c.FormValue("applyMoney")
	if applyMoney == "" {
		return utils.ErrorNull(c, "请输入提现金额")
	}
	am := convert.MustFloat64(applyMoney)
	if am <= 0 {
		return utils.ErrorNull(c, "请输入提现金额")
	}
	m := getWallet(acc)
	if m == nil {
		return utils.ErrorNull(c, "金额已超过可提现金额")
	}
	walletMoney := convert.MustFloat64(m["money"])
	//邀请返利，需要一个月后才可提现
	unWithdrawCash := getUnWithdrawCashMoney(acc)
	applyWithdrawCash := getApplyWithdrawCashMoney(acc)
	doMoney := walletMoney - unWithdrawCash
	global.Log.Info("提现金额：%v元,余额：%v元，不可提现金额：%v元，正在审核中的提现金额：%v元,可提现金额：%v元",
		am, walletMoney, unWithdrawCash, applyWithdrawCash, doMoney)
	if am > doMoney {
		return utils.ErrorNull(c, fmt.Sprintf("金额已超过可提现金额"))
	}
	flog := false
	global.DB.Tx(func(tx *mysql.SqlConnTransaction) {
		applyMap := map[string]interface{}{}
		applyMap["id"] = utils.ID()
		applyMap["account_id"] = acc.ID
		applyMap["apply_money"] = am
		applyMap["collect_platform"] = m["wc_collect_platform"]
		applyMap["card_number"] = m["wc_card_number"]
		applyMap["card_name"] = m["wc_card_name"]
		applyMap["card_qrcode"] = m["wc_card_qrcode"]
		applyMap["ct_time"] = utils.CurrentTime()
		applyMap["last_time"] = applyMap["ct_time"]
		applyMap["ct_ip"] = c.RealIP()
		applyMap["status"] = enum.REVIEW_NO
		_, err = tx.InsertMap("account_apply_withdraw_cash", applyMap)
		if err != nil {
			global.Log.Error(err.Error())
			panic(err)
		}
		//扣除余额
		_, err = tx.Update("UPDATE account_wallet SET money=money-? WHERE account_id=?", am, acc.ID)
		if err != nil {
			global.Log.Error(err.Error())
			panic(err)
		}
		flog = true
	}, func(err error) {
		if err != nil {
			global.Log.Error("申请提现失败，%v", err)
			flog = false
		}
	})

	if !flog {
		return utils.ErrorNull(c, "申请提现失败")
	}
	return utils.SuccessNull(c, "申请提现成功")
}

/**
申请提现正在审核的金额
*/
func getApplyWithdrawCashMoney(acc *global.Account) float64 {
	sql := "SELECT sum(apply_money) as applyMoney FROM account_apply_withdraw_cash WHERE `status`=? AND account_id=?"
	rows, _ := global.DB.Query(sql, enum.REVIEW_NO, acc.ID)
	if rows == nil || len(rows) != 1 {
		return 0
	}
	return convert.MustFloat64(rows[0]["applyMoney"])
}

func GetApplyWithdrawCashList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "account_apply_withdraw_cash",
		Where:     "account_id=?",
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Order:     "last_time DESC",
	}, acc.ID)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func GetShareAccountList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "v_account",
		Where:     "status=? AND shareAccountId=?",
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Order:     "ct_time DESC",
	}, enum.NORMAL, acc.ID)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

func GetAppSetting(c echo.Context) error {
	advertTextList, _ := GetAdvertTextList(enum.APP, "")
	advertImgList, _ := GetAdvertList(enum.APP, "")
	return utils.SuccessNullMsg(c, map[string]interface{}{
		"advert":    advertTextList,
		"advertImg": advertImgList,
	})
}

func UpdatePassword(c echo.Context) error {
	pwd := c.FormValue("pwd")
	pwdNew := c.FormValue("pwdNew")
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	sql := `SELECT id, password,salt FROM account WHERE id = ? AND status = ?`
	rows, _ := global.DB.Query(sql, acc.ID, enum.NORMAL)
	if len(rows) != 1 {
		return utils.Error(c, "帐号不存在或已禁用", nil)
	}
	userInfo := rows[0]
	pwdOld := convert.ToString(userInfo["password"])
	salt := convert.ToString(userInfo["salt"])
	if encrypt.Sha1Encode(pwd+salt) != pwdOld {
		return utils.Error(c, "原密码错误", nil)
	}

	sql = `UPDATE account SET password=? WHERE id = ? AND status = ?`
	x, err := global.DB.Update(sql, encrypt.Sha1Encode(pwdNew+salt), acc.ID, enum.NORMAL)
	if err != nil {
		return utils.ErrorNull(c, err.Error())
	}
	if x > 0 {
		return utils.Success(c, "修改密码成功", nil)
	} else {
		return utils.Error(c, "修改密码失败", nil)
	}
}

func GetAccountAuths(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	sql := `SELECT * FROM account_auths WHERE account_id = ?`
	rows, _ := global.DB.Query(sql, acc.ID)
	return utils.SuccessNullMsg(c, rows)
}

func AddActionLogs(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	t := c.FormValue("action")
	if t == "" {
		return utils.ErrorNull(c, "无有效的功能行为记录分类")
	}
	_, err = global.DB.InsertMap("account_action_logs", utils.ActionLogs(c, t, acc.ID))
	if err != nil {
		return utils.ErrorNull(c, err.Error())
	}
	return utils.SuccessNull(c, "success")
}

//申请身份认证
func AddIdentityAuth(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	name := c.FormValue("name")
	cardNumber := c.FormValue("cardNumber")
	photoFront := c.FormValue("photoFront")
	photoBack := c.FormValue("photoBack")
	if name == "" {
		return utils.ErrorNull(c, "请填写身份名称")
	}
	if cardNumber == "" {
		return utils.ErrorNull(c, "请填写身份号")
	}
	if photoFront == "" {
		return utils.ErrorNull(c, "请填写身份正面照片")
	}
	if photoBack == "" {
		return utils.ErrorNull(c, "请填写身份反面照片")
	}
	identityAuth := getIdentityAuth(acc)
	if identityAuth != nil {
		status := convert.ToString(identityAuth["status"])
		if status == enum.REVIEW_OK {
			return utils.ErrorNull(c, "已通过审核请勿重复提交")
		} else if status == enum.REVIEW_IN {
			return utils.ErrorNull(c, "正在审核中请耐心等待")
		}
		sql := "UPDATE account_identity_auth SET type=?,name=?,card_number=?,photo_front=?,photo_back=?,status=?,ct_time=?,ct_ip=? WHERE account_id=?"
		_, err = global.DB.Insert(sql, enum.IDENTITY_ID_CARD, name, cardNumber, photoFront, photoBack, enum.REVIEW_NO, utils.CurrentTime(), c.RealIP(), acc.ID)
		if err != nil {
			global.Log.Error(fmt.Sprintf("申请认证信息保存失败，错误：%s", err.Error()))
			return utils.ErrorNull(c, "申请认证信息保存失败")
		}
	} else {
		sql := "INSERT INTO account_identity_auth(id,account_id,type,name,card_number,photo_front,photo_back,status,ct_time,ct_ip)VALUES(?,?,?,?,?,?,?,?,?,?)"
		_, err = global.DB.Insert(sql, utils.ID(), acc.ID, enum.IDENTITY_ID_CARD, name, cardNumber, photoFront, photoBack, enum.REVIEW_NO, utils.CurrentTime(), c.RealIP())
		if err != nil {
			global.Log.Error(fmt.Sprintf("申请认证信息保存失败，错误：%s", err.Error()))
			return utils.ErrorNull(c, "申请认证信息保存失败")
		}
	}
	return utils.SuccessNull(c, "提交实名认证成功，请耐心等待审核！")
}

func GetIdentityAuth(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	identityAuth := getIdentityAuth(acc)
	if identityAuth == nil {
		return utils.NullData(c)
	}
	return utils.SuccessNullMsg(c, identityAuth)
}

func getIdentityAuth(acc *global.Account) map[string]interface{} {
	sql := `SELECT * FROM account_identity_auth WHERE account_id = ? LIMIT 1`
	rows, err := global.DB.Query(sql, acc.ID)
	if err != nil {
		global.Log.Error(fmt.Sprintf("获取认证信息错误，ERROR：%s", err.Error()))
		return nil
	}
	if len(rows) != 1 {
		return nil
	}
	return rows[0]
}

func GetAccountById(id int64) map[string]interface{} {
	sql := `SELECT * FROM account WHERE id = ? AND status = ? LIMIT 1`
	rows, err := global.DB.Query(sql, id, enum.NORMAL)
	if err != nil {
		global.Log.Error(fmt.Sprintf("GetAccountById sql ERROR：%s", err.Error()))
		return nil
	}
	if len(rows) != 1 {
		return nil
	}
	return rows[0]

}

func GetAccountByUnionId(unionId, t string) map[string]interface{} {
	sql := `SELECT a.* FROM account_auths aa LEFT JOIN account a ON aa.account_id=a.id WHERE aa.unionID = ? AND aa.type=? LIMIT 1`
	rows, err := global.DB.Query(sql, unionId, t)
	if err != nil {
		global.Log.Error(fmt.Sprintf("GetAccountByWechatUnionId sql ERROR：%s", err.Error()))
		return nil
	}
	if len(rows) != 1 {
		return nil
	}
	return rows[0]
}
