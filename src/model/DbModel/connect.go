package DbModel

import (
	"crypto/md5"
	"database/sql"
	. "fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"model"
	"model/Rom"
	"strings"
	//"time"
	//"reflect"
)

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
type DB sql.DB

//连接数据库
func Connect(config map[string]string) (db *sql.DB, err error) {
	//func Connect(db *sql.DB, config map[string]string) db *sql.DB {
	//var err error
	//Println(config["user"] + ":" + config["passwd"] + "@tcp(" + config["host"] + ":" + config["port"] + ")/" + config["name"])
	//
	db, err = sql.Open("mysql", config["user"]+":"+config["passwd"]+"@tcp("+config["host"]+":"+config["port"]+")/"+config["dbName"])
	//db, err = sql.Open("mysql", config.DB_USER+":"+config.DB_PASSWORD+"@tcp("+config.DB_HOST+":"+config.DB_PORT+")/"+config.DB_NAME)
	//强行连接，判断是不是连接成功
	model.CheckErr(err, "open数据库失败")
	model.CheckErr(db.Ping(), "ping数据库失败")
	return db, nil
}

//查询数据库
//func Query(query string) (r Result) {
//select的时候，必须全部取出来
func Query(db *sql.DB, query string) (r *Result) {
	var res Result
	var err error
	rows, err := db.Query(query)
	model.CheckErr(err, "查询失败")
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

//单纯获取value，不要k，如show tables这类命令,结果只有一列多行
//select的时候，必须全部取出来
func GetArr(db *sql.DB, query string) ([]string, int) {
	buff := make([]string, 16, 32)
	//var res []string
	rows, err := db.Query(query)
	defer rows.Close()
	model.CheckErr(err, "查询失败")

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
		model.CheckErr(rows.Scan(scanArgs...), "获取行失败")
		cnt = 0
		//Println(len(values))
		//Println(values)
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
func QueryLines(db *sql.DB, query string) ([]string, map[int][]string) {
	//data := make(map[int][]string)
	data := make(map[int][]string, 64)
	rows, err := db.Query(query)
	defer rows.Close()
	model.CheckErr(err, "查询失败")
	Columns, _ := rows.Columns()
	//Println(Columns)
	scanArgs := make([]interface{}, len(Columns))
	values := make([]interface{}, len(Columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	cnt := 0
	for rows.Next() {
		model.CheckErr(rows.Scan(scanArgs...), "获取行失败")

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
func QueryOneRow(db *sql.DB, query string) ([]string, []string) {
	//data := make(map[int][]string)
	var data []string
	rows, err := db.Query(query)
	defer rows.Close()
	model.CheckErr(err, "查询失败, query is "+query)
	Columns, _ := rows.Columns()
	//Println(Columns)
	scanArgs := make([]interface{}, len(Columns))
	values := make([]interface{}, len(Columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	cnt := 0
	for rows.Next() {
		model.CheckErr(rows.Scan(scanArgs...), "获取行失败")
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

/*
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
*/

func GetPri(db *sql.DB, tab string) string {
	//Println("desc " + tab)
	tabHead, tabLine := QueryLines(db, "desc "+tab)
	key := -1
	pri := ""
	for i, head := range tabHead {
		if head == "Key" {
			key = i
			break
		}
	}

	if key != -1 {
		for _, data := range tabLine {
			if data[key] == "PRI" {
				pri = data[0]
				break
			}
			//Println(data)
		}
	} else {
		Println(tab, " has no key###################################")
	}
	//Println(tab, "has pri_key ", pri)
	return pri
}

//把对应的map插入到数据库
func Insert(db *sql.DB, data map[string]string) int64 {
	/*
		if db1 == nil {
			db1.Connect(config.DB1)
		}
		if db2 == nil {
			db2.Connect(config.DB2)
		}
	*/
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
	model.CheckErr(err, "插入数据库失败")
	if err == nil {
		return -1
	}
	id, err := res.LastInsertId()
	model.CheckErr(err, "获取插入id失败")
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

//查询两个表的相同表，相同id是否相同
func CheckRow(db1 *sql.DB, db2 *sql.DB, sql string) bool {
	_, Row1 := QueryOneRow(db1, sql)
	_, Row2 := QueryOneRow(db2, sql)
	b1 := GetMd5(Row1)
	b2 := GetMd5(Row2)
	if Sprintf("%x", b1) == Sprintf("%x", b2) {
		return true
	}
	return false
	//io.WriteString(h, "And Leon's getting laaarger!")
	//Printf("%x", b1)
	//Printf("%x", b2)
}

func GetMd5(row []string) [16]byte {
	h := md5.New()
	for _, v := range row {
		io.WriteString(h, v)
	}
	return md5.Sum(nil)
}
