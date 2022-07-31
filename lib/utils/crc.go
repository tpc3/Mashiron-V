package utils

import (
	"hash/crc64"
	"strconv"
)

var table *crc64.Table

func init() {
	table = crc64.MakeTable(crc64.ECMA)
}

func Crc(str string) string {
	data := []byte(str)
	return strconv.FormatUint(crc64.Checksum(data, table), 10)
}
