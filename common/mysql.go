package common

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewMysqlConn() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "pimg:Wuzeping#55@tcp(42.192.43.242:3306)/testP")
	return
}

func GetResultRow(rows *sql.Rows) map[string]string {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([][]byte, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[string]string)
	for rows.Next() {
		//将行数据保存到record字典
		rows.Scan(scanArgs...)
		for i, v := range values {
			if v != nil {
				record[columns[i]] = string(v)
			}
		}
	}
	return record
}

func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	//返回所有列
	columns, _ := rows.Columns()
	vals := make([][]byte, len(columns))
	scans := make([]interface{}, len(columns))

	for k := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		rows.Scan(scans...)
		row := make(map[string]string)
		for k, v := range vals {
			key := columns[k]
			row[key] = string(v)
		}
		result[i] = row
		i++
	}
	return result
}
