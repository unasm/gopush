package model

import (
	"config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//连接数据库
func Connect() {
	var err error
	db, err = sql.Open("mysql", config.DB_USER+":"+config.DB_PASSWORD+"@tcp("+config.DB_HOST+":"+config.DB_PORT+")/"+config.DB_NAME)
	checkErr(err)
}

//查询数据库
//func Query(query string) (r Result) {
//select的时候，必须全部取出来
func Query(query string) (r *Result) {
	var res Result
	var err error
	rows, err := db.Query(query)
	checkErr(err)
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
