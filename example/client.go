package formosa_test

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"
	"strconv"
	"encoding/json"
	"github.com/matishsiao/goformosa/formosa/client/go/formosa"
)

var (
	h string
	p int
)

func ExampleClient() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&h, "h", "127.0.0.1", "127.0.0.1")
	flag.IntVar(&p, "p", 4001, "port")
	flag.Parse()
	MultiClient()
}

func MultiClient() {
	c, err := formosa.Connect(h, p, "sa23891odi1@8hfn!0932aqiomc9AQjiHH")
	if err != nil {
		log.Fatal("Connect Error:", err)
	}
	KVTest(c)
	//HashTest(c, "test")
	//HashTest(c, "test2")
	//QueueTest(c)
	//BatchWrite(c)
	c.Close()
}

func BatchWrite(c *formosa.Client) {
	var data [][]string
	data = append(data,[]string{"hset", "batchTest","Test", "TestValue"})
	data = append(data,[]string{"hset","batchTest","Test1", "TestValue1"})
	data = append(data,[]string{"hset", "batchTest","Test2", "TestValue2"})
	data = append(data,[]string{"hset", "batchTest","Test3", "TestValue3"})
	data = append(data,[]string{"hset", "batchTest","Test4", "TestValue4"})
	data = append(data,[]string{"set", "batchTest1", "TestValue1"})
	data = append(data,[]string{"set", "batchTest2", "TestValue2"})
	data = append(data,[]string{"set", "batchTest3", "TestValue3"})

	jsonStr,err := json.Marshal(data)
	if err != nil {
		log.Println("error:",err)
		return
	}

	value, err := c.Do("batchwrite", string(jsonStr))
	if err != nil {
		log.Println("qsize Error:", err)
	}
	log.Println("batchwrite:", value)
	value, err = c.Do("hscan","batchTest","","",-1)
	if err != nil {
		log.Println("hscan Error:", err)
	}
	log.Println("hscan:", value)
	value, err = c.Do("scan","","",-1)
	if err != nil {
		log.Println("scan Error:", err)
	}
	log.Println("scan:", value)
}

func QueueTest(c *formosa.Client) {
	value, err := c.Do("qsize", "Test")
	if err != nil {
		log.Println("qsize Error:", err)
	}
	log.Println("qsize", value, err)
	if len(value) == 2 && value[0] == "ok" {
		size := ToInt64(value[1])
		for i:=0;i < int(size);i++ {
			value, err = c.Do("qpop", "Test")
			if err != nil {
				log.Println("qpop Error:", err)
			}
			log.Println("qpop from loop", value, err)
		}
	}
	value, err = c.Do("qpush", "Test", "TestValue")
	if err != nil {
		log.Println("Set Error:", err)
	}
	log.Println("Qpush", value, err)
	value, err = c.Do("qpush", "Test", "TestValue1")
	if err != nil {
		log.Println("Set Error:", err)
	}
	log.Println("Qpush", value, err)
	value, err = c.Do("qpush", "Test", "TestValue2")
	if err != nil {
		log.Println("Set Error:", err)
	}
	log.Println("Qpush", value, err)
	value, err = c.Do("qpop", "Test")
	if err != nil {
		log.Println("qpop Error:", err)
	}
	log.Println("qpop", value, err)
	value, err = c.Do("qpop", "Test")
	if err != nil {
		log.Println("qpop Error:", err)
	}
	log.Println("qpop", value, err)
	value, err = c.Do("qpop", "Test")
	if err != nil {
		log.Println("qpop Error:", err)
	}
	log.Println("qpop", value, err)
	value, err = c.Do("qpop", "Test")
	if err != nil {
		log.Println("qpop Error:", err)
	}
	log.Println("qpop", value, err)
	value, err = c.Do("qpush", "Test", "TestValue3")
	if err != nil {
		log.Println("Set Error:", err)
	}
	log.Println("Qpush", value, err)
	value, err = c.Do("qpop", "Test")
	if err != nil {
		log.Println("qpop Error:", err)
	}
	log.Println("qpop", value, err)
	value, err = c.Do("qpop", "Test")
	if err != nil {
		log.Println("qpop Error:", err)
	}
	log.Println("qpop", value, err)
	value, err = c.Do("qpush", "Test", "TestValue4")
	if err != nil {
		log.Println("Set Error:", err)
	}
	log.Println("Qpush", value, err)
	value, err = c.Do("qpush", "Test", "TestValue5")
	log.Println("Qpush", value, err)

	value, err = c.Do("qsize", "Test")
	if err != nil {
		log.Println("qsize Error:", err)
	}
	log.Println("qsize", value, err)
}

func GlobalTest(c *formosa.Client) {
	value, err := c.Do("globalgetall")
	if err != nil {
		log.Fatal("globalgetall Error:", err)
	}
	log.Println("globalgetall", value)
}

func KVTest(c *formosa.Client) {
	/*value, err := c.Do("set", "Test", "TestValue")
	if err != nil {
		log.Fatal("Set Error:", err)
	}
	value, err = c.Do("set", "Test1", "TestValue1")
	if err != nil {
		log.Fatal("Set Error:", err)
	}
	value, err = c.Do("set", "Test2", "TestValue2")
	if err != nil {
		log.Fatal("Set Error:", err)
	}
	log.Println("Set", value)
	value, err = c.Do("get", "Test")
	if err != nil {
		log.Fatal("get Error:", err)
	}
	log.Println("get", value)
	value, err = c.Do("scan", "", "", -1)
	if err != nil {
		log.Fatal("scan Error:", err)
	}
	log.Println("scan", value)
	value, err = c.Do("size")
	if err != nil {
		log.Fatal("size Error:", err)
	}
	log.Println("size", value)*/
	value, err := c.Do("incr", "Incr", 1)
	if err != nil {
		log.Fatal(c.Id, "incr Error:", err)
	}
	log.Println(c.Id, "incr", value)
	value, err = c.Do("incr", "Incr", 3)
	if err != nil {
		log.Fatal(c.Id, "incr Error:", err)
	}
	log.Println(c.Id, "incr", value)
	value, err = c.Do("get", "Incr")
	if err != nil {
		log.Fatal("get Error:", err)
	}
	log.Println(c.Id, "incr get", value)
}

func HashInitTest(c *formosa.Client) {
	if err != nil {
		log.Fatal("hset Error:", err)
	}
	log.Println("hset", value)
	/*value, err = c.Do("hget", "Test", "A")
	if err != nil {
		log.Fatal("hget Error:", err)
	}
	log.Println("hget", value)*/
}

func HashTest(c *formosa.Client, hash string) {
	value, err := c.Do("hset", hash, hash+"A", hash+"ValueA")
	if err != nil {
		log.Fatal("hset Error:", err)
	}
	value, err = c.Do("hset", hash, hash+"B", hash+"ValueB")
	if err != nil {
		log.Fatal("hset Error:", err)
	}
	log.Println(hash, "hset", value)
	value, err = c.Do("hset", hash, hash+"C", hash+"ValueC")
	if err != nil {
		log.Fatal("hset Error:", err)
	}
	log.Println(hash, "hset", value)

	value, err = c.Do("hget", hash, "A")
	if err != nil {
		log.Fatal("hget Error:", err)
	}
	log.Println(hash, "hget", value)

	value, err = c.Do("hincr", hash, hash+"Incr", 1)
	if err != nil {
		log.Fatal(c.Id, "hincr Error:", err)
	}

	value, err = c.Do("hget", hash, hash+"Incr")
	if err != nil {
		log.Fatal("hget Error:", err)
	}
	log.Println(hash, c.Id, "hget incr", value)
	value, err = c.Do("hscan", hash, "", "", -1)
	if err != nil {
		log.Fatal("hscan Error:", err)
	}
	for k, v := range value {
		log.Println(hash, "hscan", k, v)
	}

	value, err = c.Do("hsize", hash)
	if err != nil {
		log.Fatal("hsize Error:", err)
	}
	log.Println(hash, "hsize", value)
}
func BenchmarkHashGet(c *formosa.Client) {
	//c.Do("zip", 1)
	start := time.Now().UnixNano()
	for i := 1; i <= 30; i++ {
		data, err := c.Do("hget", "Test", "A")
		if err != nil {
			log.Fatal("hset Error:", err)
		}
		log.Printf("data[%v]:%v\n", i, len(data[1]))
	}

	//result, err := c.Exec()
	end := (time.Now().UnixNano() - start) / 1000000
	//log.Printf("use time:%d ms, Result:%v Error:%v\n", end, result, err)
	log.Printf("use time:%d ms\n", end)

}

func Benchmark(c *formosa.Client) {
	start := time.Now().UnixNano()
	for i := 1; i <= 1000; i++ {
		name := fmt.Sprintf("HashTest-%d", i)
		//c.BatchAppend("hset", "HashTest", name, name)
		_, err := c.Do("hset", "HashTest", name, name)
		if err != nil {
			log.Fatal("hset Error:", err)
		}
	}

	//result, err := c.Exec()
	end := (time.Now().UnixNano() - start) / 1000000
	//log.Printf("use time:%d ms, Result:%v Error:%v\n", end, result, err)
	log.Printf("use time:%d ms\n", end)
}

func BenchmarkBatch(c *formosa.Client) {
	start := time.Now().UnixNano()
	for i := 1; i <= 1000; i++ {
		name := fmt.Sprintf("HashTest-%d", i)
		c.BatchAppend("hset", "HashTest", name, name)
	}

	result, err := c.Exec()
	end := (time.Now().UnixNano() - start) / 1000000
	log.Printf("use time:%d ms, Result:%v Error:%v\n", end, result, err)
}

func ToInt64(data string) int64 {
	val, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		log.Println("Ticker ParseInt error", err, data)
		return 0
	}
	return val
}