package main

import (
	"fmt"
	"model"
	//	"reflect"
	//"log"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	model.Connect()
	//var ans *model.Result
	var ans *model.Result
	ans = model.Query("SELECT * FROM  list where id < 10")
	//fmt.Println("dfadsf", ans)
	for k, v := range ans.Fields {
		fmt.Println("idx is : ", k)
		fmt.Println("idx is : ", v.Symbol, "time is ", v.Symbol, v.Id)
	}
}
