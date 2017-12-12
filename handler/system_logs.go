package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/pay/global"
	"github.com/beewit/beekit/utils/enum"
)

func AddSystemLog(c echo.Context) error {
	source := c.FormValue("source")
	t := c.FormValue("type")
	content := c.FormValue("content")
	if source==""{
		source=enum.ACTION_PC
	}
	if t == "" {
		t = "ERROR"
	}
	_, err := global.DB.InsertMap("system_logs", map[string]interface{}{
		"source":  source,
		"type":    t,
		"content": content,
		"ct_time": utils.CurrentTime(),
		"ct_ip":   c.RealIP(),
		"id":      utils.GetIp(),
	})
	if err != nil {
		return utils.ErrorNull(c, err.Error())
	}
	return utils.SuccessNull(c, "添加日志成功")
}
