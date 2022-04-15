package create

import (
	"bytes"
	_ "embed"
	"io"
	"strings"
)

//go:embed mysql.sql
var mysql []byte

func GetSql() ([]string, error) {
	r := bytes.NewReader(mysql)
	b, err := io.ReadAll(r)
	if err == nil {
		sqlArr := strings.Split(string(b), "--")
		var sql = []string{}
		for _, v := range sqlArr {
			if strings.Trim(v, "--") != "" {
				sql = append(sql, "--"+v)
			}
		}

		return sql, nil
	}

	return nil, err
}
