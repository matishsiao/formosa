package main

import (
	"fmt"
)

func EncodeBinlog(key string) string {
	return fmt.Sprintf("%s!%s", DATATYPE_BINLOG, key)
}

func EncodeBinlogNode(node string) string {
	return fmt.Sprintf("%s%s!", DATATYPE_BINLOG, node)
}

func EncodeHash(hash string, key string) string {
	return fmt.Sprintf("%s%s!%s", DATATYPE_HASH, hash, key)
}

func EncodeHashEnd(hash string) string {
	return fmt.Sprintf("%s%s#", DATATYPE_HASH, hash)
}

func DecodeHash(enkey string, hash string) string {
	if len(enkey) > 0 {
		if string(enkey[0]) != DATATYPE_HASH {
			return ""
		}
		enkey = enkey[2+len(hash):]
		return enkey
	}
	return ""
}

func EncodeQueue(queue string, key string) string {
	return fmt.Sprintf("%s%s!%s", DATATYPE_QUEUE, queue, key)
}

func EncodeQueueEnd(queue string) string {
	return fmt.Sprintf("%s%s#", DATATYPE_QUEUE, queue)
}

func EncodeQueueFront(queue string) string {
	return fmt.Sprintf("%s%s!%s", DATATYPE_QUEUE, queue, DATATYPE_QUEUE_FRONT)
}
func EncodeQueueRear(queue string) string {
	return fmt.Sprintf("%s%s!%s", DATATYPE_QUEUE, queue, DATATYPE_QUEUE_REAR)
}

func EncodeKV(key string) string {
	return fmt.Sprintf("%s!%s", DATATYPE_KV, key)
}

func DecodeKV(enkey string) string {
	if len(enkey) > 0 {
		if string(enkey[0]) != DATATYPE_KV {
			return ""
		}
		return enkey[1:]
	}
	return ""
}
