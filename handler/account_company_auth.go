package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"strings"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/enum"
)

func SubmitAccountCompanyAuth(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	name := strings.TrimSpace(c.FormValue("name"))
	legalPerson := strings.TrimSpace(c.FormValue("legal_person"))
	tel := strings.TrimSpace(c.FormValue("tel"))
	address := strings.TrimSpace(c.FormValue("address"))
	businessLicence := strings.TrimSpace(c.FormValue("business_licence"))
	idCardJust := strings.TrimSpace(c.FormValue("id_card_just"))
	idCardBack := strings.TrimSpace(c.FormValue("id_card_back"))
	if name == "" {
		return utils.ErrorNull(c, "请填写企业名称")
	}
	if len(name) > 255 {
		return utils.ErrorNull(c, "企业名称，最长不超过255字")
	}
	if legalPerson == "" {
		return utils.ErrorNull(c, "请填写企业法人")
	}
	if len(legalPerson) > 50 {
		return utils.ErrorNull(c, "企业法人，最长不超过50字")
	}
	if tel == "" {
		return utils.ErrorNull(c, "请填写联系电话")
	}
	if len(tel) > 50 {
		return utils.ErrorNull(c, "联系电话，最长不超过50字")
	}
	if address == "" {
		return utils.ErrorNull(c, "请填写办公地址")
	}
	if len(address) > 255 {
		return utils.ErrorNull(c, "办公地址，最长不超过255字")
	}
	if businessLicence == "" {
		return utils.ErrorNull(c, "请上传营业执照")
	}
	if len(businessLicence) > 255 {
		return utils.ErrorNull(c, "营业执照，最长不超过255字")
	}
	if idCardJust == "" {
		return utils.ErrorNull(c, "请上传身份证(正)")
	}
	if len(idCardJust) > 255 {
		return utils.ErrorNull(c, "身份证(正)，最长不超过255字")
	}
	if idCardBack == "" {
		return utils.ErrorNull(c, "请填写身份证(反)")
	}
	if len(idCardBack) > 255 {
		return utils.ErrorNull(c, "身份证(反)，最长不超过255字")
	}
	ip := c.RealIP()
	currentTime := utils.CurrentTime()
	m := GetAccountCompanyAuth(acc.ID)
	if m != nil {
		if m["status"] == enum.REVIEW_OK {
			return utils.ErrorNull(c, "企业认证审核已通过，请勿继续提交！")
		}
		x, err := global.DB.Update("UPDATE account_company_auth SET name=?,legal_person=?,address=?,business_licence=?,tel=?,id_card_just=?,"+
			"id_card_back=?,status=?,ut_time=?,ip=? WHERE id=?", name, legalPerson, address, businessLicence, tel, idCardJust, idCardBack,
			enum.REVIEW_IN, currentTime, ip, m["id"])
		if err != nil {
			global.Log.Error("UPDATE account_company_auth sql error:%s", err.Error())
			return utils.ErrorNull(c, "提交企业认证失败")
		}
		if x > 0 {
			return utils.SuccessNull(c, "已提交企业认证，请耐心等待审核！")
		} else {
			return utils.ErrorNull(c,"提交企业认证失败")
		}

	} else {
		//新增
		m = map[string]interface{}{}
		m["id"] = utils.ID()
		m["account_id"] = acc.ID
		m["ct_time"] = currentTime
		m["ut_time"] = currentTime
		m["ip"] = ip
		m["status"] = enum.REVIEW_IN
		m["name"] = name
		m["legal_person"] = legalPerson
		m["address"] = address
		m["business_licence"] = businessLicence
		m["tel"] = tel
		m["id_card_just"] = idCardJust
		m["id_card_back"] = idCardBack
		_,err:=global.DB.InsertMap("account_company_auth",m)
		if err != nil {
			global.Log.Error("insert account_company_auth sql error:%s", err.Error())
			return utils.ErrorNull(c, "提交企业认证失败")
		}
		return utils.SuccessNull(c, "已提交企业认证，请耐心等待审核！")
	}
}

func GetAccountCompanyAuthByID(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	m := GetAccountCompanyAuth(acc.ID)
	if m == nil {
		return utils.NullData(c)
	}
	return utils.SuccessNullMsg(c, m)
}

func GetAccountCompanyAuth(accId int64) map[string]interface{} {
	m, err := global.DB.Query("SELECT * FROM account_company_auth WHERE account_id=? LIMIT 1", accId)
	if err != nil {
		global.Log.Error("account_company_auth sql ERROR：", err.Error())
		return nil
	}
	if len(m) != 1 {
		return nil
	}
	return m[0]
}
