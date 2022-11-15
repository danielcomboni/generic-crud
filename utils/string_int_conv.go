package utils

import (
	"log"
	"strconv"
)

func ConvertStrToInt64(str string) int64 {
	intValue, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Println("failed to convert string to int64: " + err.Error())
		return 0
	}
	return intValue
}

func ConvertInt64ToStr(intValue int64) string {
	s := strconv.FormatInt(intValue, 10)
	return s
}
