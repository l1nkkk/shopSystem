package common

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// NewMysqlConn 创建mysql 连接
func NewMysqlConn() (db *sql.DB, err error) {
	// 账号:密码@tcp(地址:端口)/数据库?charset=utf8
	db, err = sql.Open("mysql", "root:linuxlin000@tcp(127.0.0.1:3306)/imooc?charset=utf8")
	return
}

// GetResultRow 获取返回值，获取一条
func GetResultRow(rows *sql.Rows) map[string]string{
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}	
	record := make(map[string]string)

	// l1nkkk: 这种写法，看起来是以最后一行的数据作为返回
	for rows.Next() {
		//将行数据保存到record字典
		rows.Scan(scanArgs...)
		for i, v := range values {
			if v != nil {
				//fmt.Println(reflect.TypeOf(col))
				record[columns[i]] = string(v.([]byte))
			}
		}
	}
	return record
}

// GetResultRows 获取所有
func GetResultRows(rows *sql.Rows) map[int]map[string]string {
	// 返回一行中所有列的列名
	columns, _ := rows.Columns()

	// 这里表示一行中所有列的值，用[]byte表示
	vals := make([][]byte, len(columns))
	// 这里表示一行填充数据的地址
	scans := make([]interface{}, len(columns))


	// 这里scans引用vals，为了后面调用rows.Scan，把数据填充到[]byte里
	for k, _ := range vals {
		// scans 存的是&[]byte，即对应每个 vals 中元素的地址
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		// 第一层循环，遍历行

		// 填充数据，传入 列数个地址，用于接收行数据
		rows.Scan(scans...)
		row := make(map[string]string)

		// 把vals中的数据复制到row中
		for k, v := range vals {
			// 第二层循环，遍历列数据

			// 获得列名
			key := columns[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result[i] = row
		i++
	}
	return result
}

