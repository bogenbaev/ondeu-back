package utils

import "strconv"

func ParseUint(s string) (uint, error) {
	parsed, err := strconv.ParseUint(s, 10, 64)
	return uint(parsed), err
}
