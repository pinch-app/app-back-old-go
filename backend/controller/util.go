package controller

import "strconv"

func ParseUint64(str string) (uint64, error) {
	id, err := strconv.ParseUint(str, 10, 64)
	return id, err
}
