package miniutils

import "reflect"

// GetIndexOf 获取切片元素的位置，不存在返回-1
func GetIndexOf[T any](val T, vals []T) int64 {
	var result int64 = -1
	for i, v := range vals {
		if reflect.DeepEqual(v, val) {
			result = int64(i)
			break
		}
	}
	return result
}
