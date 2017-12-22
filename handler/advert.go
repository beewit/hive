package handler

import (
	"fmt"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/enum"
)

func GetAdvertList(t string) ([]map[string]interface{}, error) {
	where := ""
	if t != "" {
		where += fmt.Sprintf(" AND type='%s'", t)
	}
	sql := "SELECT * FROM advert_text WHERE status=? " + where
	rows, err := global.DB.Query(sql, enum.NORMAL)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
