package handler

import (
	"github.com/labstack/echo"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils"
)

/**
	添加安装日志记录
 */
func AddInstallLog(c echo.Context) error {
	source := c.FormValue("source")
	deviceCode := c.FormValue("deviceCode")
	device := c.FormValue("device")
	if source == "" {
		source = "beewit"
	}
	_, err := global.DB.InsertMap("install_log", map[string]interface{}{
		"source":      source,
		"device_code": deviceCode,
		"device":      device,
		"ct_time":     utils.CurrentTime(),
		"ct_ip":       c.RealIP(),
		"id":          utils.ID(),
	})
	if err != nil {
		return utils.ErrorNull(c, err.Error())
	}
	return utils.SuccessNull(c, "添加安装日志成功")
}
