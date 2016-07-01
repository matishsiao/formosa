package main

import (
	"log"
	"strconv"
	"syscall"
)

func SetUlimit(number uint64) {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Println("[Error]: Getting Rlimit ", err)
	}
	rLimit.Max = number
	rLimit.Cur = number
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Println("[Error]: Setting Rlimit ", err)
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Println("[Error]: Getting Rlimit ", err)
	}
	log.Println("set file limit done:", rLimit)
}

func ToInt64(data string) int64 {
	val, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		log.Println("Ticker ParseInt error", err, data)
		return 0
	}
	return val
}
