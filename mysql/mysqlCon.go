// mysqlCon
package gomysql

import (
	"log"

	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type GoMysql struct {
	db *sql.DB
}

func NewMysql(host string, port string, username string, password string, database string, charset string) GoMysql {
	driverName := "gomysql"
	dataSourceName := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=" + charset
	log.Println(dataSourceName)
	connMaxLifetime := 59 * time.Second
	maxOpenConns := 3
	maxIdleConns := 5

	db, _ = sql.Open(driverName, dataSourceName)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.Ping()

	gomysql := GoMysql{db}

	return gomysql
}

func (gomysql *GoMysql) NewQueryRow(sql string, args ...interface{}) (result map[string]interface{}) {
	rows, err := db.Query(sql, args...)
	if err != nil {
		log.Fatal(err)
	}

	columns, err := rows.Columns()
	count := len(columns)

	defer rows.Close()

	value := make([]string, count)
	val := make([]interface{}, count)
	result = make(map[string]interface{}, 15) // TODO why require
	for rows.Next() {
		for i := 0; i < count; i++ {
			val[i] = &value[i]
		}

		rows.Scan(val...)
		for i, col := range columns {
			result[col] = value[i]
		}
		break
	}

	return result
}

func (gomysql *GoMysql) NewQuery(sql string, args ...interface{}) (results []map[string]interface{}) {
	rows, err := db.Query(sql, args...)
	if err != nil {
		log.Fatal(err)
	}

	columns, err := rows.Columns()
	count := len(columns)

	defer rows.Close()

	value := make([]string, count)
	val := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			val[i] = &value[i]
		}

		rows.Scan(val...)
		result := make(map[string]interface{}, 15)
		for i, col := range columns {
			result[col] = value[i]
		}
		results = append(results, result)
	}

	return results
}
