package handler

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/hive/global"
	"github.com/beewit/wechat/mp/user/oauth2"
	"github.com/labstack/echo"
)

var MPSessionId string = "mpSessionId"

func readBody(c echo.Context) (map[string]string, error) {
	body, bErr := ioutil.ReadAll(c.Request().Body)
	if bErr != nil {
		global.Log.Error("读取http body失败，原因：", bErr.Error())
		return nil, bErr
	}
	defer c.Request().Body.Close()

	var bm map[string]string
	bErr = json.Unmarshal(body, &bm)
	if bErr != nil {
		global.Log.Error("解析http body失败，原因：", bErr.Error())
		return nil, bErr
	}
	return bm, bErr
}

func Filter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var token string
		token = c.FormValue("token")
		if token == "" {
			bm, _ := readBody(c)
			if bm != nil {
				token = bm["token"]
			}
		}
		if token == "" {
			return utils.AuthFail(c, "登陆信息token无效，请重新登陆")
		}

		accMapStr, err := global.RD.GetString(token)
		if err != nil {
			global.Log.Error(err.Error())
			return utils.AuthFail(c, "登陆信息已失效，请重新登陆")
		}
		if accMapStr == "" {
			global.Log.Error(token + "已失效")
			return utils.AuthFail(c, "登陆信息已失效，请重新登陆")
		}
		accMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(accMapStr), &accMap)
		if err != nil {
			global.Log.Error(accMapStr + "，error：" + err.Error())
			return utils.AuthFail(c, "登陆信息已失效，请重新登陆")
		}
		m, err := global.DB.Query("SELECT id,nickname,photo,mobile,status FROM account WHERE id=? LIMIT 1", accMap["id"])
		if err != nil {
			return utils.AuthFail(c, "获取用户信息失败")
		}
		if convert.ToString(m[0]["status"]) != enum.NORMAL {
			return utils.AuthFail(c, "用户已被冻结")
		}

		c.Set("account", global.ToMapAccount(m[0]))
		return next(c)
	}
}

func GetAccount(c echo.Context) (acc *global.Account, err error) {
	itf := c.Get("account")
	if itf == nil {
		err = utils.AuthFailNull(c)
		return
	}
	acc = global.ToInterfaceAccount(itf)
	if acc == nil {
		err = utils.AuthFailNull(c)
		return
	}
	return
}

type WxSesstion struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
}

//func GetMiniAppSession(c echo.Context) (*WxSesstion, error) {
//	miniAppSessionId := strings.TrimSpace(c.FormValue("miniAppSessionId"))
//	if miniAppSessionId == "" {
//		return nil, errors.New("未识别到用户标识")
//	}
//	wsStr, err := global.RD.GetString(miniAppSessionId)
//	if err != nil {
//		return nil, errors.New("未识别到用户标识")
//	}
//	var ws *WxSesstion
//	err = json.Unmarshal([]byte(wsStr), &ws)
//	if err != nil {
//		return nil, errors.New("获取用户登录标识失败")
//	}
//	return ws, nil
//}

func GetOauthUser(c echo.Context) *oauth2.UserInfo {
	mpSessionId := strings.TrimSpace(c.FormValue(MPSessionId))
	if mpSessionId == "" {
		return nil
	}
	us, err := global.RD.GetString(mpSessionId)
	if err != nil {
		return nil
	}
	var u *oauth2.UserInfo
	err = json.Unmarshal([]byte(us), &u)
	if err != nil {
		return nil
	}
	return u
}
