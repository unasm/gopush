package StaffModel

import (
	"config"
	"database/sql"
	"fmt"
	"model/Rom"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	//每40s发送一次心跳信号，如果确认，延长关闭一份中
	TABLE = "company_staff_801"
)

type Field struct {
	mobilephone uint32
	id_no       string
	name        string
	id          uint32
	company_id  uint64
}
type Result struct {
	Status uint16
	Fields map[int]Field
	//key的列表
	//Keys []string
}
type Model struct {
	db      *sql.DB
	TabName string
}

/*
func (db *Model) Init() {

}
*/
func NewObj() *Model {
	/*
		var model Model
		model.db = db
		model.tabName = "company_staff_801"
	*/
	//return &model
	var db *sql.DB
	return &Model{db, "company_staff_801"}
	//return &Model{db, "company_staff_801"}
}

//连接数据库
func (table *Model) Connect() {
	//func Connect(db *sql.DB, config map[string]string) db *sql.DB {
	//var err error
	//Println(config["user"] + ":" + config["passwd"] + "@tcp(" + config["host"] + ":" + config["port"] + ")/" + config["name"])
	//
	//table.db, err = sql.Open("mysql", config["user"]+":"+config["passwd"]+"@tcp("+config["host"]+":"+config["port"]+")/"+config["dbName"])
	dsn := config.DB_USER + ":" + config.DB_PASSWORD + "@tcp(" + config.DB_HOST + ":" + config.DB_PORT + ")/" + config.DB_NAME
	fmt.Println(dsn)
	db, err := sql.Open("mysql", dsn)
	table.db = db
	//强行连接，判断是不是连接成功
	//model.CheckErr(err, "open数据库失败")
	//model.CheckErr(db.Ping(), "ping数据库失败")
	if err != nil {
		fmt.Println(err)
		panic("connection error")
	}
}

//查询数据库
//func Query(query string) (r Result) {
//select的时候，必须全部取出来
/*
func (table *Model) Query(query string) (r *Result) {
	var res Result
	var err error
	rows, err := table.db.Query(query)
	//model.CheckErr(err, "查询失败")
	//Keys, _ := rows.Columns()
	cnt := 0
	tmap := make(map[int]Field)
	for rows.Next() {
		var tmp Field
		rows.Scan(&tmp.Id, &tmp.Symbol, &tmp.Item, &tmp.Value, &tmp.Times)
		tmap[cnt] = tmp
		cnt++
		Printf("%d\n", cnt)
	}
	cnt = 0
	for rows.Next() {
		Printf("%d\n", cnt)
	}
	res.Fields = tmap
	defer rows.Close()
	return &res
}
*/

//单纯获取value，不要k，如show tables这类命令,结果只有一列多行
//select的时候，必须全部取出来
func (table *Model) getBySql(query string) ([]string, int) {
	buff := make([]string, 16, 32)
	//var res []string
	rows, _ := table.db.Query(query)
	defer rows.Close()
	//model.CheckErr(err, "查询失败")

	Columns, _ := rows.Columns()
	cnum := len(Columns)
	scanArgs := make([]interface{}, cnum)
	values := make([]interface{}, cnum)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	idx := 0
	cnt := 0
	line := make([]string, cnum)
	for rows.Next() {
		//model.CheckErr(rows.Scan(scanArgs...), "获取行失败")
		cnt = 0
		if len(values) == 1 && values[0] == nil {
			break
		}
		for _, col := range values {
			//Println(col)
			line[cnt] = string(col.([]byte))
			cnt++
		}
		//每行必须只能有一个，结果才能只是一个Arr,不是getOneRow的情况
		buff = Rom.ExtendStringArr(buff)
		if cnt == 1 {
			buff = buff[0 : idx+1]
			buff[idx] = line[0]
			idx++
		}
	}
	return buff, idx
}

//单纯获取value，不要k，如show tables这类命令
//select的时候，必须全部取出来
func (table *Model) QueryLines(query string) ([]string, map[int][]string) {
	//data := make(map[int][]string)
	data := make(map[int][]string, 64)
	rows, _ := table.db.Query(query)
	defer rows.Close()
	//model.CheckErr(err, "查询失败")
	Columns, _ := rows.Columns()
	//Println(Columns)
	scanArgs := make([]interface{}, len(Columns))
	values := make([]interface{}, len(Columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	cnt := 0
	for rows.Next() {
		//model.CheckErr(rows.Scan(scanArgs...), "获取行失败")

		lineBuf := make([]string, len(scanArgs))
		for i, col := range values {
			if col != nil {
				lineBuf[i] = string(col.([]byte))
			}
		}
		//data = Rom.ExtendMapStringArr(data)
		data[cnt] = lineBuf
		cnt++
	}
	return Columns, data
	//return &res
}

//只获取一行数据
// @return []string  字段名
// @return []string  值
func (table *Model) QueryOne(query string) ([]string, []string) {
	//data := make(map[int][]string)
	var data []string
	rows, _ := table.db.Query(query)
	defer rows.Close()
	//model.CheckErr(err, "查询失败, query is "+query)
	Columns, _ := rows.Columns()
	//Println(Columns)
	scanArgs := make([]interface{}, len(Columns))
	values := make([]interface{}, len(Columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	cnt := 0
	for rows.Next() {
		//model.CheckErr(rows.Scan(scanArgs...), "获取行失败")
		if len(values) == 1 && values[0] == nil {
			break
		}
		data = make([]string, len(scanArgs))
		for i, col := range values {
			//Println(col)
			data[i] = string(col.([]byte))
			cnt++
		}
		return Columns, data
		//每行必须只能有一个，结果才能只是一个Arr,不是getOneRow的情况
	}
	return Columns, data
}

//把对应的map插入到数据库
func (table *Model) Insert(data map[string]string) int64 {
	//fmt.Println("asdfa")
	//return 1
	//if table.db == nil {
	table.Connect()
	//}
	keys := ""
	values := ""
	for k, v := range data {
		keys += "`" + k + "`,"
		values += "'" + v + "',"
	}
	keys = strings.Trim(keys, ",")
	values = strings.Trim(values, ",")
	sql := "INSERT INTO " + table.TabName + " (" + keys + ") VALUES (" + values + ")"
	res, err := table.db.Exec(sql)
	//model.CheckErr(err, "插入数据库失败")
	if err == nil {
		return -1
	}
	id, err := res.LastInsertId()
	//model.CheckErr(err, "获取插入id失败")
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
