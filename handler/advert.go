package handler

import (
	"fmt"
	"github.com/beewit/hive/global"
	"github.com/beewit/beekit/utils/enum"
)

func GetAdvertTextList(cls, t string) ([]map[string]interface{}, error) {
	where := ""
	if t != "" {
		where += fmt.Sprintf(" AND type='%s'", t)
	}
	sql := "SELECT * FROM advert_text WHERE status=? AND class=? " + where
	rows, err := global.DB.Query(sql, enum.NORMAL, cls)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func GetAdvertList(cls, t string) ([]map[string]interface{}, error) {
	where := ""
	if t != "" {
		where += fmt.Sprintf(" AND type='%s'", t)
	}
	sql := "SELECT * FROM advert WHERE status=? AND class=? " + where+" ORDER BY sort ASC"
	rows, err := global.DB.Query(sql, enum.NORMAL, cls)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
