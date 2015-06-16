package model

import (
	"config"
	"database/sql"
	"fmt"
	. "fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
	//"reflect"
)

var db *sql.DB

type Field struct {
	Id     uint32
	Item   string
	Value  string
	Times  uint32
	Symbol string
}
type Result struct {
	Status uint16
	Fields map[int]Field
	//key的列表
	//Keys []string
}

//连接数据库
func Connect() {
	var err error

	Println(config.DB_USER + ":" + config.DB_PASSWORD + "@tcp(" + config.DB_HOST + ":" + config.DB_PORT + ")/" + config.DB_NAME)
	//
	db, err = sql.Open("mysql", config.DB_USER+":"+config.DB_PASSWORD+"@tcp("+config.DB_HOST+":"+config.DB_PORT+")/"+config.DB_NAME)
	//强行连接，判断是不是连接成功
	CheckErr(err, "open数据库失败")
	CheckErr(db.Ping(), "ping数据库失败")
}

//查询数据库
//func Query(query string) (r Result) {
//select的时候，必须全部取出来
func Query(query string) (r *Result) {
	var res Result
	var err error
	rows, err := db.Query(query)
	CheckErr(err, "查询失败")
	//Keys, _ := rows.Columns()
	/*
		for k, col := range res.Keys {
			fmt.Printf("%d \t %10s \n", k, col)
		}
	*/
	cnt := 0
	tmap := make(map[int]Field)
	for rows.Next() {
		var tmp Field
		rows.Scan(&tmp.Id, &tmp.Symbol, &tmp.Item, &tmp.Value, &tmp.Times)
		tmap[cnt] = tmp
		cnt++
		fmt.Printf("%d\n", cnt)
	}
	cnt = 0
	for rows.Next() {
		fmt.Printf("%d\n", cnt)
	}
	res.Fields = tmap
	defer rows.Close()
	return &res
}

func ErrorInsert(data map[string]string) {
	tmpStamp := Sprintf("%d", time.Now().Unix())
	tmp := make(map[string]string)
	for k, v := range data {
		tmp["item"] = k
		tmp["value"] = v
		tmp["link"] = tmpStamp
		Insert(tmp)
	}
}

func Insert(data map[string]string) int64 {
	if db == nil {
		//Println("connecting")
		Connect()
	}
	//sql := "INSERT error"
	keys := ""
	values := ""
	//values := make([]string, len(data))
	//cnt := 0
	for k, v := range data {
		keys += "`" + k + "`,"
		values += "'" + v + "',"
	}
	keys = strings.Trim(keys, ",")
	values = strings.Trim(values, ",")
	sql := "INSERT INTO error (" + keys + ") VALUES (" + values + ")"
	Println(sql)
	res, err := db.Exec(sql)
	CheckErr(err, "插入数据库失败")
	if err == nil {
		return -1
	}
	id, err := res.LastInsertId()
	CheckErr(err, "获取插入id失败")
	return id

}

//查询的样例
/*
func example() {
	model.Connect()
	var ans *model.Result
	ans = model.Query("SELECT * FROM  list where id < 10")
	for k, v := range ans.Fields {
		fmt.Println("idx is : ", k)
		fmt.Println("idx is : ", v.Symbol, "time is ", v.Symbol, v.Id)
	}
}
*/
