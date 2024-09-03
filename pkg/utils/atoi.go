package utils

import "strconv"

func Atoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}