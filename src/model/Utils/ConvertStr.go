package Utils

import (
	"strconv"
)

//把字符串数组，转化成int数组
func ConvertStrToInt(strs []string) []int {
	res := make([]int, len(strs))
	for i, v := range strs {
		res[i], _ = strconv.Atoi(v)
	}
	return res
}

//判断字符串是不是在指定的字符串数组中
func InArray(needle string, haystack []string) bool {

	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
