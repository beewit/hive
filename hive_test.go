package main

import (
	"testing"
	"encoding/json"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/beekit/utils"
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
