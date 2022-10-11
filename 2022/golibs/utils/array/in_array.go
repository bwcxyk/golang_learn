package array

// InArray 判断是否是字符串数组中的一员
func InArray(val string, array []string) (ok bool, i int) {
	for i = range array {
		if ok = array[i] == val; ok {
			return
		}
	}
	return
}
