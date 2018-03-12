package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/hive/global"
	"github.com/beewit/hive/handler"
	"github.com/beewit/wechat-ai/smartWechat"
)

func TestRedis(t *testing.T) {
	accMapStr, err := global.RD.GetString("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.YCypNF-PULJ2zgYCumBYtQg4YmiE41O6HgE1hocZd5Q")
	if err != nil {
		global.Log.Error(err.Error())

	}
	if accMapStr == "" {
		global.Log.Error("已失效")

	}
	accMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(accMapStr), &accMap)
	if err != nil {
		global.Log.Error(accMapStr + "，error：" + err.Error())
	}
}

func TestRules(t *testing.T) {
	rp, err := ApiPost("http://127.0.0.1:8085/api/rules/list", map[string]string{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.1SK0Uj1P_uu1YH-I3_p2JwSNlnb9jGjIjaYmWWLuGzA"})
	if err != nil {
		t.Error(err.Error())
	}
	str, err2 := json.Marshal(rp)
	if err2 != nil {
		t.Error(err2.Error())
	}
	println(string(str))
}

func ApiPost(url string, m map[string]string) (utils.ResultParam, error) {
	b, _ := json.Marshal(m)
	body, err := uhttp.Cmd(uhttp.Request{
		Method: "POST",
		URL:    url,
		Body:   b,
	})
	if err != nil {
		return utils.ResultParam{}, err
	}
	return utils.ToResultParam(body), nil
}

func TestUpdatePwd(t *testing.T) {
	rp, err := ApiPost("http://127.0.0.1:8085/api/account/updatePwd?pwd=123456&pwdNew=1234567", map[string]string{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.LEHRhrcsscya5MbirqEmsqwX8SPqFzIAqm8MU-lGJvQ"})
	if err != nil {
		t.Error(err.Error())
	}
	str, err2 := json.Marshal(rp)
	if err2 != nil {
		t.Error(err2.Error())
	}
	println(string(str))
}

func TestGetUrlPars(t *testing.T) {
	s := smartWechat.GetURLParams(map[string]string{
		"234": "iii",
		"456": "32r"})
	println(s)
}

func TestTime(t *testing.T) {
	println(fmt.Sprintf("%d", time.Now().Unix()))
	_, err := time.Parse("01-02-2006", "02-08-2015")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestFloat(t *testing.T) {
	println(fmt.Sprintf("%v", utils.IsValidNumber("12.6")))
}
func TestPageNumber(t *testing.T) {
	println(fmt.Sprintf("%v", int(math.Ceil(float64(10)/float64(4)))))
}

func TestImg(t *testing.T) {
	rep, err := http.Get("http://sso.9ee3.com/img/code")
	if err != nil {
		println(err.Error())
		return
	}
	_, err = ioutil.ReadAll(rep.Body)
	if err != nil {
		println(err.Error())
		return
	}
	println("可怜")
}

func TestGetReceiveRedPacketAndCouponList(t *testing.T) {
	var id int64
	id = 6436906930684928
	redPacket := handler.GetRedPacket(id)
	if redPacket == nil {
		t.Error("红包不存在或已过期")
	}
	sql := "SELECT ar.money,ar.ct_time as receiveTime,wa.* FROM account_receive_red_packet ar LEFT JOIN wx_account wa ON ar.wx_union_id =wa.union_id WHERE account_send_red_packet_id=?"
	redPacketList, err := global.DB.Query(sql, id)
	if err != nil {
		global.Log.Error("GetReceiveRedPacketAndCouponList account_receive_red_packet sql error:%s", err.Error())
		t.Error("获取领取红包数据失败", nil)
	}
	couponList := []map[string]interface{}{}
	joinCouponIds := convert.ToString(redPacket["join_coupon_ids"])
	if joinCouponIds != "" {
		sql = fmt.Sprintf("SELECT ac.money,ar.ct_time as receiveTime,wa.* FROM account_receive_coupon ar LEFT JOIN wx_account wa ON ar.wx_union_id =wa.union_id "+
			"LEFT JOIN account_coupon ac ON ac.id=ar.account_coupon_id WHERE account_coupon_id in(%s)", joinCouponIds)
		couponList, err = global.DB.Query(sql)
		if err != nil {
			global.Log.Error("GetReceiveRedPacketAndCouponList account_receive_coupon sql error:%s", err.Error())
			t.Error("获取领取现金券数据失败", nil)
		}
	}

	t.Error(map[string]interface{}{
		"redPacket":     redPacket,
		"redPacketList": redPacketList,
		"couponList":    couponList})
}
func TestGetWechatMiniUnionID(t *testing.T) {
	where := " AND (expire_time is NULL OR expire_time>now()) AND number>receive_number"
	pageIndex := utils.GetPageIndex("1")
	pageSize := utils.GetPageSize("10")
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "account_coupon",
		Where:     "account_id=? AND status=?" + where,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, 6416401854972928, enum.NORMAL)
	if err != nil {
		t.Error("数据异常，" + err.Error())
		return
	}
	if page == nil {
		t.Error("无数据")
		return
	}
	t.Error("无数据")
	return
}

func TestRandInt(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	money := 1 + r.Intn(2-1)
	println(money)
	f := float32(money) / float32(100)
	println(fmt.Sprintf("%.2f", f))
}

func TestRandIntByUtils(t *testing.T) {
	for i := 0; i < 10000; i++ {
		ft := utils.NewRandom().NumberByFloat(5.0, 50.0)
		println(fmt.Sprintf("%.2f", ft))
		in := utils.NewRandom().NumberByInt(1, 10000)
		println(in)
		if ft > 50 || ft <= 0 {
			println("------------------------------------------------------------- error")
		}
		if in > 10000 || in <= 0 {
			println("------------------------------------------------------------- error")
		}
		if in > 9000 {
			println("------------------------------------------------------------- yes ###################################3")
		}
	}
}

func TestFloatStr(t *testing.T) {
	money := convert.MustFloat64("10.00")
	feeMoney := convert.MustFloat64("0.2")
	totalPrice := money + feeMoney
	println(convert.MustInt(fmt.Sprintf("%.2f", totalPrice*100)))
}

func TestFloat2(t *testing.T) {
	var a float64 = 8.23
	var b float64 = 0.2
	println(fmt.Sprintf("%.2f", a*b))
	println(fmt.Sprintf("%.2f",a - b))
	println(fmt.Sprintf("%.2f",a + b))
}
