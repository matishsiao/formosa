package formosa

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	_ "io"
	"log"
	"net"
	_ "runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pquerna/ffjson/ffjson"
)

//Service Client for Proxy
type SrvClient struct {
	Conn        *net.TCPConn
	mu          *sync.Mutex
	RemoteAddr  string
	RequestTime int64
	recvBuf     bytes.Buffer
	Auth        bool
	Zip         bool
	TmpResult   []string
	Connected   bool
}

//Client Init
func (cl *SrvClient) Init(conn *net.TCPConn) {
	cl.Conn = conn
	cl.Connected = true
	if CONFIGS.Password == "" {
		cl.Auth = true
	}
	cl.RemoteAddr = strings.Split(cl.Conn.RemoteAddr().String(), ":")[0]
	go cl.HealthCheck()
	cl.Read()
	//PrintGCSummary()
}

//Close Client connection
func (cl *SrvClient) Close() {
	cl.mu.Lock()
	cl.Conn.Close()
	cl.Connected = false
	cl.Auth = false
	cl.mu.Unlock()
	ProxyConn--
	if ProxyConn <= 0 {
		ProxyConn = 0
	}
}

func (cl *SrvClient) HealthCheck() {
	cl.RequestTime = time.Now().Unix()
	timeout := 1
	for cl.Connected {
		if time.Now().Unix()-cl.RequestTime >= CONFIGS.Timeout {
			log.Println("HealthCheck Service Connection by Timeout:", cl.Conn.RemoteAddr())
			cl.Close()
			break
		}
		time.Sleep(time.Duration(timeout) * time.Second)
	}

	//receyle client
	cl = nil
}

func (cl *SrvClient) Read() {
	//timeout := 100

	for cl.Connected {
		data, err := cl.Recv()
		if err != nil {
			if CONFIGS.Debug {
				log.Printf("Srv Client Receive Error:%v RemoteAddr:%s\n", err, cl.RemoteAddr)
			}
			cl.Close()
			break
		} else {
			cl.RequestTime = time.Now().Unix()
			if len(data) > 0 {
				//start := time.Now().UnixNano()
				cl.Process(data)
				//end := (time.Now().UnixNano() - start) / 1000000
				//log.Printf("use time:%d ms", end)
			}
			//timeout = 10
		}
		//time.Sleep(time.Duration(timeout) * time.Microsecond)
	}
}

func (cl *SrvClient) Process(req []string) {

	if len(req) == 0 {
		//ok, not_found, error, fail, client_error
		cl.Send([]string{"error", "request format incorrect."}, false)
	} else {
		switch req[0] {
		case "auth":
			if CONFIGS.Password != "" {
				if len(req) == 2 {
					if req[1] == CONFIGS.Password {
						cl.Auth = true
						cl.Send([]string{"ok", "1"}, false)
					} else {
						cl.Send([]string{"fail", "password incorrect."}, false)
					}
				} else {
					cl.Send([]string{"fail", "request format incorrect"}, false)
				}
			} else {
				cl.Auth = true
				cl.Send([]string{"ok", "1"}, false)
			}
		case "ping":
			cl.Send([]string{"ok", "1"}, false)
		case "batchexec":
			if cl.Auth {
				if len(req) == 2 {
					var cmdlist [][]string
					err := ffjson.Unmarshal([]byte(req[1]), &cmdlist)
					if err != nil {
						cl.Send([]string{"fail", "batchexec need use json format."}, false)
						return
					}
					var resultlist [][]string
					async := false
					if len(cmdlist) > 0 && len(cmdlist[0]) >= 1 && cmdlist[0][0] == "async" {
						async = true
						cl.Send([]string{"ok", "batchexec use async mode."}, false)
					}
					//err = DM.Batch(cmdlist)
					for _, v := range cmdlist {
						res, err := cl.Query(v)
						if err != nil {
							resultlist = append(resultlist, []string{"error", err.Error()})
							continue
						}
						if res == nil {
							resultlist = append(resultlist, []string{"not_found"})
							continue
						}
						resultlist = append(resultlist, res)
					}
					if !async {
						resultjson, err := ffjson.Marshal(resultlist)
						if err != nil {
							cl.Send([]string{"fail", "batchexec send json result failed." + err.Error()}, false)
							return
						}
						res := []string{"ok", string(resultjson)}
						if cl.Zip {
							cl.Send(res, true)
						} else {
							cl.Send(res, false)
						}
					}
				}
			} else {
				cl.Send([]string{"fail", "not auth"}, false)
			}
		case "batchwrite":
			if cl.Auth {
				if len(req) == 2 {
					var cmdlist [][]string
					err := ffjson.Unmarshal([]byte(req[1]), &cmdlist)
					if err != nil {
						cl.Send([]string{"fail", "batchwrite need use json format."}, false)
						return
					}
					async := false
					if len(cmdlist) > 0 && len(cmdlist[0]) >= 1 && cmdlist[0][0] == "async" {
						async = true
						cl.Send([]string{"ok", "batchwrite use async mode."}, false)
					}
					err = DM.Batch(cmdlist)
					if !async {
						var res []string
						if err != nil {
							res = []string{"fail", err.Error()}
						} else {
							res = []string{"ok", "1"}
						}
						if cl.Zip {
							cl.Send(res, true)
						} else {
							cl.Send(res, false)
						}
					}
				}
			} else {
				cl.Send([]string{"fail", "not auth"}, false)
			}
		case "zip":
			if cl.Auth {
				if len(req) == 2 {
					if req[1] == "1" {
						cl.Zip = true
					} else {
						cl.Zip = false
					}
					cl.Send([]string{"ok", req[1]}, false)
				}
			} else {
				cl.Send([]string{"fail", "not auth"}, false)
			}
		default:
			if cl.Auth {
				res, err := cl.Query(req)
				if err != nil {
					cl.Send([]string{"error", err.Error()}, false)
					return
				}
				if CONFIGS.Debug {
					log.Println("Response:", res)
				}
				if res == nil {
					cl.Send([]string{"not_found"}, false)
					return
				} else {
					if cl.Zip {
						cl.Send(res, true)
					} else {
						cl.Send(res, false)
					}
					return
				}
			} else {
				cl.Send([]string{"error", "you need login first"}, false)
				return
			}
		}
	}
}

func (cl *SrvClient) Query(args []string) ([]string, error) {
	find := false
	if CONFIGS.Debug {
		log.Println("Query:", args)
	}
	if len(args) == 0 {
		return nil, fmt.Errorf("bad request:request args length incorrect.")
	}
	var mapList map[string]string
	var tmpList []string
	var response []string
	var counter int
	errFlag := false
	var errMsg error
	if len(args) > 0 {
		//log.Println("Args:", args, binlog.seq)
		switch args[0] {
		case "globalgetall":
			if len(args) == 1 {
				return DM.GlobalGetAll()
			} else {
				return response, fmt.Errorf("Args length not equl 2.")
			}
		case "hset":
			if len(args) == 4 {
				err := DM.HashSet(args[1], args[2], args[3])
				if err != nil {
					return response, err
				} else {
					//
					response = append(response, "ok")
					response = append(response, "1")
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 4.")
			}
		case "hdel":
			if len(args) == 3 {
				err := DM.HashDel(args[1], args[2])
				if err != nil {
					return response, err
				} else {
					response = append(response, "ok")
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 3.")
			}
		case "hget":
			if len(args) == 3 {
				data, err := DM.HashGet(args[1], args[2])
				if err != nil {
					return response, err
				} else {
					response = append(response, "ok")
					response = append(response, string(data))
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 3.")
			}
		case "hsize":
			if len(args) == 2 {
				data, err := DM.HashSize(args[1])
				if err != nil {
					return response, err
				} else {
					response = append(response, "ok")
					response = append(response, fmt.Sprintf("%d", data))
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 2.")
			}
		case "hscan":
			if len(args) == 5 {
				return DM.HashScan(args[1], args[2], args[3], ToInt64(args[4]))
			} else {
				return response, fmt.Errorf("Args length not equl 5.")
			}
		case "hexists":
			if len(args) == 3 {
				data, err := DM.HashExists(args[1], args[2])
				if err != nil {
					return response, err
				} else {
					response = append(response, "ok")
					if data {
						response = append(response, "1")
					} else {
						response = append(response, "0")
					}
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 4.")
			}
		case "hincr":
			if len(args) == 4 {
				data, err := DM.HashIncr(args[1], args[2], args[3])
				if err != nil {
					return response, err
				} else {
					response = append(response, "ok")
					response = append(response, data)
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 4.")
			}
		case "scan":
			if len(args) == 4 {
				return DM.Scan(args[1], args[2], ToInt64(args[3]))
			} else {
				return response, fmt.Errorf("Args length not equl 4.")
			}
		case "set":
			if len(args) == 3 {
				err := DM.Set(args[1], args[2])
				if err != nil {
					return response, err
				} else {

					response = append(response, "ok")
					response = append(response, "1")
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 3.")
			}
		case "del":
			if len(args) == 2 {
				err := DM.Del(args[1])
				if err != nil {
					return response, err
				} else {
					response = append(response, "ok")
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 2.")
			}
		case "get":
			if len(args) == 2 {
				data, err := DM.Get(args[1])
				if err != nil {
					return response, err
				} else {
					response = append(response, "ok")
					response = append(response, string(data))
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 2.")
			}
		case "exists":
			if len(args) == 2 {
				data, err := DM.Exists(args[1])
				if err != nil {
					return response, err
				} else {
					response = append(response, "ok")
					if data {
						response = append(response, "1")
					} else {
						response = append(response, "0")
					}
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 2.")
			}
		case "incr":
			if len(args) == 3 {
				data, err := DM.Incr(args[1], args[2])
				if err != nil {
					return response, err
				} else {
					response = append(response, "ok")
					response = append(response, data)
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 3.")
			}
		case "size":
			if len(args) == 1 {
				data, err := DM.Size()
				if err != nil {
					return response, err
				} else {
					response = append(response, "ok")
					response = append(response, fmt.Sprintf("%d", data))
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 1.")
			}
		case "qpush":
			if len(args) == 3 {
				err := DM.QueuePush(args[1], args[2])
				if err != nil {
					return response, err
				} else {
					response = append(response, "ok")
					response = append(response, "1")
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 3.")
			}
		case "qpop":
			if len(args) == 2 {
				data, err := DM.QueuePop(args[1])
				if err != nil {
					return response, err
				} else {
					response = append(response, "ok")
					response = append(response, string(data))
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 2.")
			}
		case "qsize":
			if len(args) == 2 {
				data, err := DM.QueueSize(args[1])
				if err != nil {
					return response, err
				} else {
					response = append(response, "ok")
					response = append(response, fmt.Sprintf("%d", data))
					return response, err
				}
			} else {
				return response, fmt.Errorf("Args length not equl 2.")
			}
		}

		/*for _, v := range cl.DBNodes {

			db := v.Client
			if CONFIGS.Debug {
				log.Printf("Process:%v Mirror:%v Args:%v Info:%v\n", process, mirror, args, v.Info)
			}
			if sync {
				if v.Info.Mode == "main" {
					if len(args) > 1 {
						val, err := db.Do(args[1:])
						if err != nil {
							errFlag = true
							errMsg = err
							quit = true
						}
						if !errFlag && val[0] == "ok" {
							response = val
							if CONFIGS.Sync {
								log.Printf("Query Mirror Sync args:%v Response:%v Error:%v Server:%s Port:%d RemoteAddr:%s\n", args, val, err, v.Info.Host, v.Info.Port, cl.RemoteAddr)
							}
							break
						} else if len(response) == 0 {
							response = val
						}
					}
				}
			} else if syncDel {
				if len(args) > 1 && v.Info.Mode != "mirror" {
					val, err := db.Do(args[1:])
					if err != nil {
						errFlag = true
						errMsg = err
					}
					log.Printf("Mirror Sync Del Process args:%v Response:%v Error:%v Server:%s Port:%d RemoteAddr:%s\n", args, val, err, v.Info.Host, v.Info.Port, cl.RemoteAddr)

					if !errFlag && val[0] == "ok" {
						response = val
						if CONFIGS.Sync {
							log.Printf("Query Mirror SyncDel args:%v Response:%v Error:%v Server:%s Port:%d RemoteAddr:%s\n", args, val, err, v.Info.Host, v.Info.Port, cl.RemoteAddr)
						}
					} else if len(response) == 0 {
						response = val
					}
				}
			} else if mirror && !process {
				if v.Info.Mode == "main" {
					val, err := db.Do(args)
					if err != nil {
						errFlag = true
						errMsg = err
					}
					if !errFlag && val[0] == "ok" {
						response = val
						var mirror_args []string
						if len(args) > 2 {
							SendSubInterface(args[1], args)
						}
						mirror_args = append(mirror_args, "mirror")
						mirror_args = append(mirror_args, args...)
						GlobalClient.Append(mirror_args)
						break
					} else {
						response = val
						break
					}
					if CONFIGS.Sync {
						log.Printf("Query Main Need Mirror args:%v Response:%v Error:%v Server:%s Port:%d RemoteAddr:%s\n", args, val, err, v.Info.Host, v.Info.Port, cl.RemoteAddr)
					}
				}
			} else if mirror && process {

				val, err := db.Do(args)
				if err != nil {
					errFlag = true
					errMsg = err
				}
				//log.Printf("Main Mirror Process args:%v Response:%v Error:%v Server:%s Port:%d RemoteAddr:%s\n", args, val, err, v.Info.Host, v.Info.Port, cl.RemoteAddr)

				if !errFlag && val[0] == "ok" {
					response = val
				} else if len(response) == 0 {
					response = val
				}
				if v.Info.Mode == "main" {
					var mirror_args []string
					mirror_args = append(mirror_args, "mirror_del")
					mirror_args = append(mirror_args, args...)
					GlobalClient.Append(mirror_args)
				}
			} else if v.Info.Mode != "mirror" {
				val, err := db.Do(args)
				if err != nil {
					errFlag = true
					errMsg = err
				}
				if CONFIGS.Debug {
					log.Println("args:", args, " Do Response:", val, "error:", err)
				}
				if !errFlag && !process && len(val) >= 1 {
					if val[0] == "ok" {
						response = val
						break
					} else if len(response) == 0 {
						response = val
					}
				}
				if !errFlag && len(val) >= 1 && val[0] != "not_found" {
					find = true
					switch args[0] {
					case "hsize":
						size, err := strconv.Atoi(val[1])
						if err != nil {
							log.Println("hsize change fail:", err, val[1])
						}
						counter += size
					case "hkeys", "keys", "rkeys", "hlist", "hrlist":
						val = val[1:]
						if CONFIGS.Debug {
							log.Println("keys val:", val)
						}
						for _, kv := range val {
							kfind := false
							for _, rv := range tmpList {
								if kv == rv {
									kfind = true
									break
								}
							}
							if !kfind {
								tmpList = append(tmpList, kv)
							}
						}
					case "del", "multi_del", "hclear", "hdel", "multi_hdel":
						response = val
					case "exists", "hexists":
						response = val
						if val[1] == "1" {
							quit = true
						}
					default:
						length := len(val[1:])
						if length%2 == 0 {
							data := val[1:]
							for i := 0; i < length; i += 2 {
								if _, ok := mapList[data[i]]; !ok {
									mapList[data[i]] = data[i+1]
								}
							}
						} else {
							log.Println("query failed:", args, "Return:", val)
							response = val
							errMsg = fmt.Errorf("bad request:request length incorrect.")
							quit = true
							errFlag = true
						}
					}
				}
			}
			if quit {
				break
			}
		}*/

	} else {
		errFlag = true
		errMsg = fmt.Errorf("bad request:request length incorrect.")
	}

	if errFlag {
		return nil, errMsg
	}

	if find {
		limit := -1
		switch args[0] {
		case "hscan", "hrscan", "scan", "rscan", "hkeys", "keys", "rkeys", "hlist", "hrlist":
			argsLimit, err := strconv.Atoi(args[len(args)-1])
			if err != nil {
				log.Println("limit parser error:", err)
			} else {
				limit = argsLimit
			}
			if CONFIGS.Debug {
				log.Println("argsLimit:", argsLimit, " args len:", args[len(args)-1])
			}
			break
		}
		switch args[0] {
		case "hgetall", "hscan", "hrscan", "multi_hget", "multi_get", "scan", "rscan":
			response = append(response, "ok")
			if len(mapList) > 0 {
				keylist := sortedKeys(mapList)
				if args[0] == "rscan" || args[0] == "hrscan" {
					sort.Sort(sort.Reverse(sort.StringSlice(keylist)))
				}
				if CONFIGS.Debug {
					log.Println("keylist:", keylist, " limit:", limit)
				}

				//if data length > limit ,cut it
				if limit != -1 && len(keylist) >= limit {
					keylist = keylist[:limit]
				}
				for _, v := range keylist {
					response = append(response, v)
					response = append(response, mapList[v])
				}
			}
			break
		case "hsize":
			response = append(response, "ok")
			response = append(response, fmt.Sprintf("%d", counter))
			break
		case "hkeys", "keys", "rkeys", "hlist", "hrlist":
			response = append(response, "ok")
			sort.Strings(tmpList)
			if args[0] == "rkeys" || args[0] == "hrlist" {
				sort.Sort(sort.Reverse(sort.StringSlice(tmpList)))
			}
			if limit != -1 {
				if len(tmpList) < limit {
					limit = len(tmpList)
				}
				tmpList = tmpList[:limit]
			}
			response = append(response, tmpList...)
			break
		}
		mapList = nil
		tmpList = nil
		return response, nil
	}
	mapList = nil
	tmpList = nil
	return response, nil

}

func (cl *SrvClient) Send(args []string, zip bool) {
	var buf bytes.Buffer
	if zip {
		buf.WriteString("3")
		buf.WriteByte('\n')
		buf.WriteString("zip")
		buf.WriteByte('\n')
		var zipbuf bytes.Buffer
		w := gzip.NewWriter(&zipbuf)
		for _, s := range args {
			w.Write([]byte(fmt.Sprintf("%d", len(s))))
			w.Write([]byte("\n"))
			w.Write([]byte(s))
			w.Write([]byte("\n"))
		}
		w.Close()
		zipbuff := base64.StdEncoding.EncodeToString(zipbuf.Bytes())
		buf.WriteString(fmt.Sprintf("%d", len(zipbuff)))
		buf.WriteByte('\n')
		buf.WriteString(zipbuff)
		buf.WriteByte('\n')
		buf.WriteByte('\n')
	} else {
		for _, s := range args {
			buf.WriteString(fmt.Sprintf("%d", len(s)))
			buf.WriteByte('\n')
			buf.WriteString(s)
			buf.WriteByte('\n')
		}
		buf.WriteByte('\n')
	}
	tmpBuf := buf.Bytes()
	_, err := cl.Conn.Write(tmpBuf)
	if err != nil {
		if CONFIGS.Debug {
			log.Printf("Srv Client Send Error:%v RemoteAddr:%s\n", err, cl.RemoteAddr)
		}
		cl.Close()
	}
}

func (cl *SrvClient) Recv() ([]string, error) {
	return cl.recv()
}

func (cl *SrvClient) recv() ([]string, error) {
	tmp := make([]byte, 102400)
	for {
		resp := cl.parse()
		if resp == nil || len(resp) > 0 {
			return resp, nil
		}
		n, err := cl.Conn.Read(tmp)
		if err != nil {
			return nil, err
		}
		cl.recvBuf.Write(tmp[0:n])
	}
}

func (cl *SrvClient) parse() []string {
	resp := []string{}
	buf := cl.recvBuf.Bytes()
	var idx, offset int
	idx = 0
	offset = 0

	for {
		idx = bytes.IndexByte(buf[offset:], '\n')
		if idx == -1 {
			break
		}
		p := buf[offset : offset+idx]
		offset += idx + 1
		//fmt.Printf("> [%s]\n", p);
		if len(p) == 0 || (len(p) == 1 && p[0] == '\r') {
			if len(resp) == 0 {
				continue
			} else {
				cl.recvBuf.Next(offset)
				return resp
			}
		}

		size, err := strconv.Atoi(string(p))
		if err != nil || size < 0 {
			return nil
		}
		if offset+size >= cl.recvBuf.Len() {
			break
		}

		v := buf[offset : offset+size]
		resp = append(resp, string(v))
		offset += size + 1
	}

	//fmt.Printf("buf.size: %d packet not ready...\n", len(buf))
	return []string{}
}
