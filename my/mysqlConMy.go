// mysqlCon
package gomysql

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//数据库配置
const (
	userName = "root"
	password = "root"
	ip       = "127.0.0.1"
	port     = "3306"
	dbName   = "test2"
)

//Db数据库连接池
var DB *sql.DB

//注意方法名大写，就是public
func init() {
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(100) //设置数据库最大连接数
	DB.SetMaxIdleConns(10)     //设置上数据库最大闲置连接数
	//验证连接
	if err := DB.Ping(); err != nil {
		log.Fatal("open database fail: ", err)
		return
	}
	log.Println("database connnect success")
}

func insert(sql string, params ...interface{}) int64 {
	stmt, err := DB.Prepare(sql)
	checkErr(err)
	res, err := stmt.Exec(params...)
	checkErr(err)
	id, err := res.LastInsertId()
	checkErr(err)
	return id
}

func update(sql string, params ...interface{}) int64 {
	stmt, err := DB.Prepare(sql)
	checkErr(err)
	res, err := stmt.Exec(params...)
	checkErr(err)
	num, err := res.RowsAffected()
	checkErr(err)
	return num
}

func deleter(sql string, params ...interface{}) int64 {
	stmt, err := DB.Prepare(sql)
	checkErr(err)
	res, err := stmt.Exec(params...)
	checkErr(err)
	num, err := res.RowsAffected()
	checkErr(err)
	return num
}

func query(sql string, params ...interface{}) {
	rows, err := DB.Query(sql, params...)
	checkErr(err)
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		record := make(map[string]interface{})
		for j, cols := range values {
			if cols != nil {
				record[columns[j]] = cols
			}
		}
		log.Println(record)
	}
}

type User struct {
	CertNo string
	Name   string
	Age    int
	Id     int
}

func queryObject(sql string, params ...interface{}) {
	rows, err := DB.Query(sql, params...)
	checkErr(err)
	var users User
	rows.Scan(users.CertNo, users.Name, users.Age, users.Id)
	log.Println(users.Name)
	//jason, john := people[0], people[1]
	//log.Printf("%#v\n%#v", jason, john)
}

//func main() {
//	//	insertSql := `INSERT user (certNo,name,age) values (?,?,?)`
//	//	inertResult := insert(insertSql, "111", "222", 2)
//	//	log.Println(inertResult)

//	//	updateSql := `UPDATE user SET age=?,certNo=? WHERE id=?`
//	//	updateResult := update(updateSql, 100, "3333", 7)
//	//	log.Println(updateResult)

//	//	deleteSql := `delete from user where id=?`
//	//	deleteResult := deleter(deleteSql, 7)
//	//	log.Println(deleteResult)

//	//	querySql := `SELECT * FROM user where id > ?`
//	//	query(querySql, 1)

//	querySql2 := `select * from user where id > ?`
//	queryObject(querySql2, 5)
//}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
