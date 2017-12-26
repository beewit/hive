package main

import (
	"testing"
	"encoding/json"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/wechat-ai/smartWechat"
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"math"
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
}

func TestFloat(t *testing.T) {
	println(fmt.Sprintf("%v",utils.IsValidNumber("12.6")))
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
