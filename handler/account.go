package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/encrypt"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"fmt"
)

func UpdatePassword(c echo.Context) error {
	pwd := c.FormValue("pwd")
	pwdNew := c.FormValue("pwdNew")
	itf := c.Get("account")
	if itf == nil {
		return utils.AuthFailNull(c)
	}
	acc := global.ToInterfaceAccount(itf)
	if acc == nil {
		return utils.AuthFailNull(c)
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
	identityAuth, err := getIdentityAuth(acc)
	if err != nil {
		global.Log.Error(fmt.Sprintf("获取认证信息错误，ERROR：%s", err.Error()))
		return utils.ErrorNull(c, "获取认证信息错误")
	}
	if identityAuth != nil {
		status := convert.ToString(identityAuth["status"])
		if status == enum.REVIEW_OK {
			return utils.ErrorNull(c, "已通过审核请勿重复提交")
		} else if status == enum.REVIEW_IN {
			return utils.ErrorNull(c, "正在审核中请耐心等待")
		}
	}
	sql := "REPLACE INTO account_identity_auth(id,account_id,type,name,card_number,photo_front,photo_back,status,ct_time,ct_ip)VALUES(?,?,?,?,?,?,?,?,?,?)"
	_, err = global.DB.Insert(sql, utils.ID(), acc.ID, enum.IDENTITY_ID_CARD, name, cardNumber, photoFront, photoBack, enum.REVIEW_NO, utils.CurrentTime(), c.RealIP())
	if err != nil {
		global.Log.Error(fmt.Sprintf("申请认证信息保存失败，错误：%s", err.Error()))
		return utils.ErrorNull(c, "申请认证信息保存失败")
	}
	return utils.SuccessNull(c, "提交实名认证成功，请耐心等待审核！")
}

func GetIndetityAuth(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return err
	}
	identityAuth, err := getIdentityAuth(acc)
	if err != nil {
		global.Log.Error(fmt.Sprintf("获取认证信息错误，ERROR：%s", err.Error()))
		return utils.ErrorNull(c, "获取认证信息错误")
	}
	if identityAuth == nil {
		return utils.NullData(c)
	}
	return utils.SuccessNullMsg(c, identityAuth)
}

func getIdentityAuth(acc *global.Account) (map[string]interface{}, error) {
	sql := `SELECT * FROM account_identity_auth WHERE account_id = ? LIMIT 1`
	rows, err := global.DB.Query(sql, acc.ID)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, nil
	}
	return rows[0], nil
}
