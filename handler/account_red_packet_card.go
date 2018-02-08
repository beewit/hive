package handler

import (
	"github.com/labstack/echo"
	"strings"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"fmt"
)

func AddRedPacketCard(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	idStr := strings.TrimSpace(c.FormValue("id"))
	name := strings.TrimSpace(c.FormValue("name"))
	contact := strings.TrimSpace(c.FormValue("contact"))
	tel := strings.TrimSpace(c.FormValue("tel"))
	address := strings.TrimSpace(c.FormValue("address"))
	def := strings.TrimSpace(c.FormValue("def"))
	if name == "" {
		return utils.ErrorNull(c, "请填写卡片名称")
	}
	if len(name) > 255 {
		return utils.ErrorNull(c, "卡片名称，最长不超过255字")
	}
	if contact == "" {
		return utils.ErrorNull(c, "请填写联系人")
	}
	if len(contact) > 255 {
		return utils.ErrorNull(c, "联系人，最长不超过255字")
	}
	if tel == "" {
		return utils.ErrorNull(c, "请填写联系方式")
	}
	if len(tel) > 50 {
		return utils.ErrorNull(c, "联系方式，最长不超过50字")
	}
	if address == "" {
		return utils.ErrorNull(c, "请填写联系地址")
	}
	if len(address) > 255 {
		return utils.ErrorNull(c, "联系地址，最长不超过255字")
	}
	if def == "true" {
		def = "1"
	} else {
		def = "0"
	}
	var id int64
	ip := c.RealIP()
	currentTime := utils.CurrentTime()
	var m map[string]interface{}
	if idStr != "" && utils.IsValidNumber(idStr) {
		id = convert.MustInt64(idStr)
		m = GetRedPacketCard(id)
	}
	if m != nil {
		if convert.MustInt64(m["account_id"]) != acc.ID {
			return utils.ErrorNull(c, "无权限修改红包卡片")
		}
		x, err := global.DB.Update("UPDATE account_red_packet_card SET name=?,contact=?,address=?,tel=?,ut_time=?,ip=?,def=?,status=? WHERE id=?",
			name, contact, address, tel, currentTime, ip, def, enum.NORMAL, id)
		if err != nil {
			global.Log.Error("UPDATE account_company_auth sql error:%s", err.Error())
			return utils.ErrorNull(c, "创建红包卡片失败")
		}
		if x > 0 {
			if def == "1" {
				removeRedPacketDefault(id, acc.ID)
			}
			return utils.SuccessNull(c, "创建红包卡片成功")
		} else {
			return utils.ErrorNull(c, "创建红包卡片失败")
		}
	} else {
		id = utils.ID()
		_, err := global.DB.InsertMap("account_red_packet_card", map[string]interface{}{
			"id":         id,
			"name":       name,
			"contact":    contact,
			"address":    address,
			"tel":        tel,
			"ct_time":    currentTime,
			"ut_time":    currentTime,
			"ip":         ip,
			"def":        def,
			"status":     enum.NORMAL,
			"account_id": acc.ID,
		})
		if err != nil {
			global.Log.Error("insert account_company_auth sql error:%s", err.Error())
			return utils.ErrorNull(c, "创建红包卡片失败")
		}
		if def == "1" {
			removeRedPacketDefault(id, acc.ID)
		}
		return utils.SuccessNull(c, "创建红包卡片成功")
	}
}

func removeRedPacketDefault(id, accId int64) bool {
	var where string
	if id > 0 {
		where = fmt.Sprintf(" AND id<>%v", id)
	}
	x, err := global.DB.Update("UPDATE account_red_packet_card SET def=? WHERE account_id=? "+where, 0, accId)
	if err != nil || x <= 0 {
		return false
	}
	return true
}

func GetRedPacketCardByID(c echo.Context) error {
	id := c.FormValue("id")
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id格式错误")
	}
	m := GetRedPacketCard(convert.MustInt64(id))
	if m == nil {
		return utils.NullData(c)
	}
	return utils.SuccessNullMsg(c, m)
}

func GetRedPacketCard(id int64) map[string]interface{} {
	m, err := global.DB.Query("SELECT * FROM account_red_packet_card WHERE id=? LIMIT 1", id)
	if err != nil {
		global.Log.Error("account_red_packet_card sql ERROR：", err.Error())
		return nil
	}
	if len(m) != 1 {
		return nil
	}
	return m[0]
}

func GetRedPacketCardDef(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	funcMap := GetEffectiveFunc(acc.ID, enum.FUNC_RED_PACKET)
	if funcMap == nil {
		return nil
	}
	m := GetRedPacketCardByDef(acc.ID)
	if m == nil {
		return utils.NullData(c)
	}
	return utils.SuccessNullMsg(c, m)
}

//获取默认红包卡片
func GetRedPacketCardByDef(accId int64) map[string]interface{} {
	m, err := global.DB.Query("SELECT * FROM account_red_packet_card WHERE account_id=? AND def=1 LIMIT 1", accId)
	if err != nil {
		global.Log.Error("account_red_packet_card sql ERROR：", err.Error())
		return nil
	}
	if len(m) != 1 {
		return nil
	}
	return m[0]
}

func DeleteRedPacketCard(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	id := c.FormValue("id")
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id格式错误")
	}
	card := GetRedPacketCard(convert.MustInt64(id))
	if card == nil {
		return utils.ErrorNull(c, "该红包卡片不存在")
	}
	if convert.MustInt64(card["account_id"]) != acc.ID {
		return utils.ErrorNull(c, "无权限删除红包卡片")
	}
	if convert.ToString(card["status"]) == enum.DELETE {
		return utils.ErrorNull(c, "已是删除状态")
	}
	x, err := global.DB.Update("UPDATE account_red_packet_card SET status=? WHERE id=?", enum.DELETE, id)
	if err != nil {
		global.Log.Error("account_red_packet_card sql ERROR：", err.Error())
		return nil
	}
	if x <= 0 {
		return utils.ErrorNull(c, "删除红包卡片失败")
	}
	return utils.SuccessNull(c, "删除红包卡片成功")
}

func GetRedPacketCardList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "account_red_packet_card",
		Where:     "account_id=? AND status=?",
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Order:     "def DESC,ct_time DESC",
	}, acc.ID, enum.NORMAL)
	if err != nil {
		global.Log.Error("QueryPage account_red_packet_card sql error:%s", err.Error())
		return utils.ErrorNull(c, "数据异常")
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}
