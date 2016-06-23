package main

import (
	"fmt"
)

func EncodeKV(key string) string {
	return fmt.Sprintf("%s%s", DATATYPE_KV, key)
}

func DecodeKV(key string) string {
	if len(key) > 0 {
		return key[1:]
	}
	return ""
}

func EncodeBinlog(key string) string {
	return fmt.Sprintf("%s%s", DATATYPE_BINLOG, key)
}

func EncodeHash(hash string, key string) string {
	return fmt.Sprintf("%s%s%d=%s%d", DATATYPE_HASH, hash, len(hash), key, len(key))
}
