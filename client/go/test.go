package main

import (
	"fmt"
	"os"

	"./formosa"
)

func main() {
	test(0)
}

func test(id int) int64 {
	//ip := "104.155.206.199";
	ip := "127.0.0.1"
	port := 8888
	db, err := formosa.Connect(ip, port)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("id:", id)
	var val interface{}
	db.Set("a", "xxx")
	val, err = db.Get("a")
	fmt.Printf("get A:%s\n", val)
	val, err = db.Get("a")
	fmt.Printf("Get val %s\n", val)
	db.Set("count", "1")
	fmt.Printf("Set Count:%d\n", 1)
	val, err = db.Incr("count", 5)
	fmt.Printf("Get Count:%s\n", val)
	//val, err = db.SetX("expeir","test",5)
	if id == 0 {
		db.SetX("expireT", "Test", 5)
	}
	val, err = db.Exists("expireT")
	fmt.Printf("Exists expireT:%v\n", val)
	if id == 1 {
		val, err = db.Expire("expireT", 5)
		fmt.Printf("Expire expireT:%v Error:%v\n", val, err)
	}
	val, err = db.KeyTTL("expireT")
	fmt.Printf("KeyTTL expireT:%v Error:%v\n", val, err)

	val, err = db.SetNew("b", "TestB")
	fmt.Printf("SetNew B:%v\n", val)

	val, err = db.GetSet("a", "new_value")
	fmt.Printf("GetSet A status:%v\n", val)
	val, err = db.Get("a")
	fmt.Printf("Get val A %s\n", val)
	val, err = db.Get("expireT")
	fmt.Printf("Get expireT:%s\n", val)
	val, err = db.Scan("", "", 10)
	for k, v := range val.(map[string]interface{}) {
		fmt.Printf("Scan[%s]%s\n", k, v)
	}
	fmt.Printf("HashGetAll:%v\n", val)
	val, err = db.HashSet("mdz-2014", "test", "10")
	val, err = db.HashSet("mdz-2014", "1231-0800", "5")
	val, err = db.HashSet("mdz-2014", "1231-0900", "1")
	val, err = db.HashSet("mdz-2014", "1231-1000", "10")
	val, err = db.HashSet("mdz-2015", "1231-1100", "5")
	val, err = db.HashSet("mdz-2015", "1231-1200", "1")
	val, err = db.HashGet("mdz-2014", "test")
	fmt.Printf("HashGet:%s\n", val)
	val, err = db.HashIncr("mdz-2014", "test", 5)
	fmt.Printf("HashIncr:%s\n", val)
	val, err = db.HashExists("hash", "test")
	fmt.Printf("HashExists:%v\n", val)
	val, err = db.HashSize("mdz-2014")
	fmt.Printf("HashSize:%d\n", val)
	val, err = db.HashScan("mdz-2014", "1230", "1231-2", 10)
	for k, v := range val.(map[string]interface{}) {
		fmt.Printf("HashScan[%s]%s\n", k, v)
	}
	fmt.Printf("HashScan:%v\n", val)
	multiSet := make(map[string]interface{})
	multiSet["A"] = 1
	multiSet["B"] = 2
	multiSet["C"] = 3
	val, err = db.HashMultiSet("mdz-2014", multiSet)
	fmt.Printf("HashMultiSet:%v\n", val)

	val, err = db.HashMultiGet("mdz-2014", []string{"A", "B"})
	for k, v := range val.(map[string]interface{}) {
		fmt.Printf("HashGetAll[%s]%s\n", k, v)
	}
	fmt.Printf("HashGetAll:%v\n", val)
	return 0
}
