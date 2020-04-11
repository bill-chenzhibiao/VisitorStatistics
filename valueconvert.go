package main

import (
	"strconv"
)

func obj2int(value interface{},df int) (result int) {
	switch value.(type) {
	case string:
		if op, ok := value.(string);ok{
			if r,err:=strconv.Atoi(op);err == nil{
				result = r
			}else {
				result = df
			}
		}else {
			result = df
		}
	case float64:
		if op, ok := value.(float64);ok{
			result = int(op)
		}else {
			result = df
		}
	default:
		result = df
	}
	return result
}
