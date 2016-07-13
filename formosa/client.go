package formosa

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	_ "syscall"
	"time"
)

type Client struct {
	sock      net.Conn
	recv_buf  bytes.Buffer
	process   chan []interface{}
	batchBuf  [][]interface{}
	result    chan ClientResult
	Id        string
	Ip        string
	Port      int
	Password  string
	Connected bool
	Retry     bool
	mu        *sync.Mutex
	Closed    bool
	init      bool
	debug     bool
}

type ClientResult struct {
	Id    string
	Data  []string
	Error error
}

type HashData struct {
	HashName string
	Key      string
	Value    string
}

var clientVersion string = "0.1.6"

const layout = "2006-01-06 15:04:05"

func GetClient(ip string, port int, auth string) (*Client, error) {
	client, err := connect(ip, port, auth)
	if err != nil {
		if client.debug {
			log.Printf("Formosa Client Connect failed:%s:%d error:%v\n", ip, port, err)
		}
		go client.RetryConnect()
		return client, err
	}
	if client != nil {
		return client, nil
	}
	return nil, nil
}

func connect(ip string, port int, auth string) (*Client, error) {
	log.Printf("Formosa Client Version:%s\n", clientVersion)
	var c Client
	c.Ip = ip
	c.Port = port
	c.Password = auth
	c.Id = fmt.Sprintf("Cl-%d", time.Now().UnixNano())
	c.mu = &sync.Mutex{}
	err := c.Connect()
	return &c, err
}

func (c *Client) Debug(flag bool) bool {
	c.debug = flag
	log.Println("Formosa Client Debug Mode:", c.debug)
	return c.debug
}

func (c *Client) Connect() error {
	log.Printf("Client[%s] connect to %s:%d\n", c.Id, c.Ip, c.Port)
	seconds := 60
	timeOut := time.Duration(seconds) * time.Second
	sock, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", c.Ip, c.Port), timeOut)
	if err != nil {
		log.Println("Formosa Client dial failed:", err, c.Id)
		return err
	}
	c.sock = sock
	c.Connected = true
	if c.Retry {
		log.Printf("Client[%s] retry connect to %s:%d success.", c.Id, c.Ip, c.Port)
	} else {
		log.Printf("Client[%s] connect to %s:%d success\n", c.Id, c.Ip, c.Port)
	}
	c.Retry = false
	if !c.init {
		c.process = make(chan []interface{})
		c.result = make(chan ClientResult)
		go c.processDo()
		c.init = true
	}

	if c.Password != "" {
		c.Auth(c.Password)
	}

	return nil
}

func (c *Client) KeepAlive() {
	go c.HealthCheck()
}

func (c *Client) HealthCheck() {
	timeout := 60
	for {
		if c.Connected && !c.Retry && !c.Closed {
			result, err := c.Do("ping")
			if err != nil {
				log.Printf("Client Health Check Failed[%s]:%v\n", c.Id, err)
			} else {
				if c.debug {
					log.Printf("Client Health Check Success[%s]:%v\n", c.Id, result)
				}
			}
		}
		time.Sleep(time.Duration(timeout) * time.Second)
	}
}

func (c *Client) RetryConnect() {
	if !c.Retry {
		c.mu.Lock()
		c.Retry = true
		c.Connected = false
		c.mu.Unlock()
		log.Printf("Client[%s] retry connect to %s:%d Connected:%v Closed:%v\n", c.Id, c.Ip, c.Port, c.Connected, c.Closed)
		for {
			if !c.Connected && !c.Closed {
				log.Printf("Client[%s] retry connect to %s:%d\n", c.Id, c.Ip, c.Port)
				err := c.Connect()
				if err != nil {
					log.Printf("Client[%s] Retry connect to %s:%d Failed. Error:%v\n", c.Id, c.Ip, c.Port, err)
					time.Sleep(5 * time.Second)
				}
			} else {
				log.Printf("Client[%s] Retry connect to %s:%d stop by conn:%v closed:%v\n.", c.Id, c.Ip, c.Port, c.Connected, c.Closed)
				break
			}
		}
	}
}

func (c *Client) CheckError(err error) {
	if err != nil {
		if !c.Closed {
			log.Printf("Check Error:%v Retry connect.\n", err)
			c.sock.Close()
			go c.RetryConnect()
		}

	}
}

func (c *Client) processDo() {
	for args := range c.process {
		runId := args[0].(string)
		runArgs := args[1:]
		result, err := c.do(runArgs)
		c.result <- ClientResult{Id: runId, Data: result, Error: err}
	}
}

func ArrayAppendToFirst(src []interface{}, dst []interface{}) []interface{} {
	tmp := src
	tmp = append(tmp, dst...)
	return tmp
}

func (c *Client) Do(args ...interface{}) ([]string, error) {
	if c.Connected && !c.Retry && !c.Closed {
		runId := fmt.Sprintf("%d", time.Now().UnixNano())
		args = ArrayAppendToFirst([]interface{}{runId}, args)
		c.process <- args
		for result := range c.result {
			if result.Id == runId {
				return result.Data, result.Error
			} else {
				c.result <- result
			}
		}
	}
	return nil, fmt.Errorf("Connection has closed.")
}

func (c *Client) BatchAppend(args ...interface{}) {
	c.batchBuf = append(c.batchBuf, args)
}

func (c *Client) Exec() ([][]string, error) {
	if c.Connected && !c.Retry && !c.Closed {
		if len(c.batchBuf) > 0 {
			runId := fmt.Sprintf("%d", time.Now().UnixNano())
			jsonStr, err := json.Marshal(&c.batchBuf)
			if err != nil {
				return [][]string{}, fmt.Errorf("Exec Json Error:%v", err)
			}
			args := []interface{}{"batchexec", string(jsonStr)}
			args = ArrayAppendToFirst([]interface{}{runId}, args)
			c.batchBuf = c.batchBuf[:0]
			c.process <- args
			for result := range c.result {
				if result.Id == runId {
					if len(result.Data) == 2 && result.Data[0] == "ok" {
						var resp [][]string
						err := json.Unmarshal([]byte(result.Data[1]), &resp)
						if err != nil {
							return [][]string{}, fmt.Errorf("Batch Json Error:%v", err)
						}
						return resp, result.Error
					} else {
						return [][]string{}, result.Error
					}

				} else {
					c.result <- result
				}
			}
		} else {
			return [][]string{}, fmt.Errorf("Batch Exec Error:No Batch Command found.")
		}
	}
	return nil, fmt.Errorf("Connection has closed.")
}

func (c *Client) do(args []interface{}) ([]string, error) {
	if c.Connected {
		err := c.send(args)
		if err != nil {
			if c.debug {
				log.Printf("Formosa Client[%s] Do Send Error:%v Data:%v\n", c.Id, err, args)
			}
			c.CheckError(err)
			return nil, err
		}
		resp, err := c.recv()
		if err != nil {
			if c.debug {
				log.Printf("Formosa Client[%s] Do Receive Error:%v Data:%v\n", c.Id, err, args)
			}
			c.CheckError(err)
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("lost connection")
}

func (c *Client) ProcessCmd(cmd string, args []interface{}) (interface{}, error) {
	if c.Connected {
		args = ArrayAppendToFirst([]interface{}{cmd}, args)
		runId := fmt.Sprintf("%d", time.Now().UnixNano())
		args = ArrayAppendToFirst([]interface{}{runId}, args)
		var err error
		c.process <- args
		var resResult ClientResult
		for result := range c.result {
			if result.Id == runId {
				resResult = result
				break
			} else {
				c.result <- result
			}
		}
		if resResult.Error != nil {
			return nil, resResult.Error
		}

		resp := resResult.Data
		if len(resp) == 2 && resp[0] == "ok" {
			switch cmd {
			case "set", "del":
				return true, nil
			case "expire", "setnx", "auth", "exists", "hexists":
				if resp[1] == "1" {
					return true, nil
				}
				return false, nil
			case "hsize":
				val, err := strconv.ParseInt(resp[1], 10, 64)
				return val, err
			default:
				return resp[1], nil
			}

		} else if len(resp) == 1 && resp[0] == "not_found" {
			return nil, fmt.Errorf("%v", resp[0])
		} else {
			if len(resp) >= 1 && resp[0] == "ok" {
				//fmt.Println("Process:",args,resp)
				switch cmd {
				case "hgetall", "hscan", "hrscan", "multi_hget", "scan", "rscan":
					list := make(map[string]string)
					length := len(resp[1:])
					data := resp[1:]
					for i := 0; i < length; i += 2 {
						list[data[i]] = data[i+1]
					}
					return list, nil
				default:
					return resp[1:], nil
				}
			}
		}
		if len(resp) == 2 && strings.Contains(resp[1], "connection") {
			c.sock.Close()
			go c.RetryConnect()
		}
		log.Printf("Formosa Client Error Response:%v args:%v Error:%v", resp, args, err)
		return nil, fmt.Errorf("bad response:%v args:%v", resp, args)
	} else {
		return nil, fmt.Errorf("lost connection")
	}
}

func (c *Client) Auth(pwd string) (interface{}, error) {
	return c.Do("auth", pwd)
	//return c.ProcessCmd("auth",params)
}

func (c *Client) Set(key string, val string) (interface{}, error) {
	params := []interface{}{key, val}
	return c.ProcessCmd("set", params)
}

func (c *Client) Get(key string) (interface{}, error) {
	params := []interface{}{key}
	return c.ProcessCmd("get", params)
}

func (c *Client) Del(key string) (interface{}, error) {
	params := []interface{}{key}
	return c.ProcessCmd("del", params)
}
func (c *Client) Scan(start string, end string, limit int) (interface{}, error) {
	params := []interface{}{start, end, limit}
	return c.ProcessCmd("scan", params)
}

//incr num to exist number value
func (c *Client) Incr(key string, val int) (interface{}, error) {
	params := []interface{}{key, val}
	return c.ProcessCmd("incr", params)
}

func (c *Client) Exists(key string) (interface{}, error) {
	params := []interface{}{key}
	return c.ProcessCmd("exists", params)
}

func (c *Client) HashSet(hash string, key string, val string) (interface{}, error) {
	params := []interface{}{hash, key, val}
	return c.ProcessCmd("hset", params)
}

func (c *Client) HashGet(hash string, key string) (interface{}, error) {
	params := []interface{}{hash, key}
	return c.ProcessCmd("hget", params)
}

func (c *Client) HashDel(hash string, key string) (interface{}, error) {
	params := []interface{}{hash, key}
	return c.ProcessCmd("hdel", params)
}

func (c *Client) HashIncr(hash string, key string, val int) (interface{}, error) {
	params := []interface{}{hash, key, val}
	return c.ProcessCmd("hincr", params)
}

func (c *Client) HashExists(hash string, key string) (interface{}, error) {
	params := []interface{}{hash, key}
	return c.ProcessCmd("hexists", params)
}

func (c *Client) HashSize(hash string) (interface{}, error) {
	params := []interface{}{hash}
	return c.ProcessCmd("hsize", params)
}

func (c *Client) HashScan(hash string, start string, end string, limit int) (map[string]string, error) {
	params := []interface{}{hash, start, end, limit}
	val, err := c.ProcessCmd("hscan", params)
	if err != nil {
		return nil, err
	} else {
		return val.(map[string]string), err
	}

	return nil, nil
}

func (c *Client) Send(args ...interface{}) error {
	return c.send(args)
}

func (c *Client) send(args []interface{}) error {
	var buf bytes.Buffer
	for _, arg := range args {
		var s string
		switch arg := arg.(type) {
		case string:
			s = arg
		case []byte:
			s = string(arg)
		case []string:
			for _, s := range arg {
				buf.WriteString(fmt.Sprintf("%d", len(s)))
				buf.WriteByte('\n')
				buf.WriteString(s)
				buf.WriteByte('\n')
			}
			continue
		case int:
			s = fmt.Sprintf("%d", arg)
		case int64:
			s = fmt.Sprintf("%d", arg)
		case float64:
			s = fmt.Sprintf("%f", arg)
		case bool:
			if arg {
				s = "1"
			} else {
				s = "0"
			}
		case nil:
			s = ""
		default:
			return fmt.Errorf("bad arguments")
		}
		buf.WriteString(fmt.Sprintf("%d", len(s)))
		buf.WriteByte('\n')
		buf.WriteString(s)
		buf.WriteByte('\n')
	}
	buf.WriteByte('\n')
	_, err := c.sock.Write(buf.Bytes())
	return err
}

func (c *Client) Recv() ([]string, error) {
	return c.recv()
}

func (c *Client) recv() ([]string, error) {
	//tmp := make([]byte, 102400)
	var tmp [102400]byte
	for {
		resp := c.parse()
		if resp == nil || len(resp) > 0 {
			//log.Println("Formosa Receive:",resp)
			if len(resp) > 0 && resp[0] == "zip" {
				//log.Println("Formosa Receive Zip\n",resp)
				zipData, err := base64.StdEncoding.DecodeString(resp[1])
				if err != nil {
					return nil, err
				}
				resp = c.UnZip(zipData)
			}
			return resp, nil
		}
		n, err := c.sock.Read(tmp[0:])
		if err != nil {
			return nil, err
		}
		c.recv_buf.Write(tmp[0:n])
	}
}

func (c *Client) parse() []string {
	resp := []string{}
	buf := c.recv_buf.Bytes()
	var Idx, offset int
	Idx = 0
	offset = 0
	for {
		Idx = bytes.IndexByte(buf[offset:], '\n')
		if Idx == -1 {
			break
		}
		p := buf[offset : offset+Idx]
		offset += Idx + 1
		if len(p) == 0 || (len(p) == 1 && p[0] == '\r') {
			if len(resp) == 0 {
				continue
			} else {
				c.recv_buf.Next(offset)
				return resp
			}
		}
		pIdx := strings.Replace(strconv.Quote(string(p)), `"`, ``, -1)
		size, err := strconv.Atoi(pIdx)
		if err != nil || size < 0 {
			return nil
		}
		if offset+size >= c.recv_buf.Len() {
			break
		}

		v := buf[offset : offset+size]
		resp = append(resp, string(v))
		offset += size + 1
	}
	return []string{}
}

func (c *Client) UnZip(data []byte) []string {
	var buf bytes.Buffer
	buf.Write(data)
	zipReader, err := gzip.NewReader(&buf)
	if err != nil {
		log.Println("[ERROR] New gzip reader:", err)
	}
	defer zipReader.Close()

	zipData, err := ioutil.ReadAll(zipReader)
	if err != nil {
		fmt.Println("[ERROR] ReadAll:", err)
		return nil
	}
	var resp []string

	if zipData != nil {
		Idx := 0
		offset := 0
		hiIdx := 0
		for {
			Idx = bytes.IndexByte(zipData, '\n')
			if Idx == -1 {
				break
			}
			p := string(zipData[:Idx])
			size, err := strconv.Atoi(string(p))
			if err != nil || size < 0 {
				zipData = zipData[Idx+1:]
				continue
			} else {
				offset = Idx + 1 + size
				hiIdx = size + Idx + 1
				resp = append(resp, string(zipData[Idx+1:hiIdx]))
				zipData = zipData[offset:]
			}

		}
	}
	return resp
}

// Close The Client Connection
func (c *Client) Close() error {
	if !c.Closed {
		c.mu.Lock()
		c.Connected = false
		c.Closed = true
		c.mu.Unlock()
		close(c.process)
		c.sock.Close()
		c = nil
	}
	return nil
}
