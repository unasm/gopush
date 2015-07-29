//进行内存的管理，比如数据的扩张

package Rom

//对字符串数组进行扩张，达到上限翻倍,动态管理内存
func ExtendStringArr(old []string) []string {
	//容量和长度相同的时候，切片达到上限，要进行扩容
	if len(old) == cap(old) {
		newSlice := make([]string, len(old), 2*(len(old)))
		copy(newSlice, old)
		old = newSlice[0:cap(newSlice)]
	}
	return old
}

//对二维字符串map数组进行扩张，达到上限翻倍,动态管理内存
/*
func ExtendMapStringArr(old map[int][]string) map[int][]string {
	//容量和长度相同的时候，切片达到上限，要进行扩容
	length := len(old)
	if length == cap(old) {
		newSlice := make(map[int][]string, length, 2*(length))
		copy(newSlice, old)
		old = newSlice[0:cap(newSlice)]
	}
	return old
}
*/
