// mysqlCon_test
package gomysql

import (
	"fmt"
	"testing"
)

func TestNewQuery(t *testing.T) {
	mysql := NewMysql("127.0.0.1", "3306", "root", "root", "test2", "utf8")
	result := mysql.NewQuery("SELECT * FROM new_news")
	fmt.Println("----------------")
	fmt.Println(result)
	if result == nil {
		t.Error("query row fail", result)
	}
}
