package model

import (
	"regexp"
)

var Form map[string][]string

func Isint(key string, value string) bool {
	if res, err := regexp.MatchString("^[0-9]+$", value); err != nil || res == false {
		//if res, err := regexp.MatchString("^[0-9]+$", r.Form["aid"][0]); err != nil || res == false {
		//Println("the num is not number")
		return false
	}
	return true
}

//检查是否存在对应的数组
func IsExistOne(key string) bool {
	if cap(Form[key]) != 1 {
		//Println("aid is NULL")
		return false
	}
	return true
}
