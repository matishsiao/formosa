package formosa

import (
	"sort"

	"github.com/matishsiao/formosa/client/go/formosa"
)

const DATATYPE_BINLOG string = "B"
const DATATYPE_HASH string = "H"
const DATATYPE_KV string = "K"
const DATATYPE_KV_END string = "K#"
const DATATYPE_QUEUE string = "Q"
const QUEUE_SIZE int64 = 100000000
const DATATYPE_QUEUE_FRONT string = "F"
const DATATYPE_QUEUE_REAR string = "R"

type Configs struct {
	Debug  bool   `json:"debug"`
	Host   string `json:"host"`
	DBPath string `json:"dbpath"`
	Web    struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"web"`
	Timeout         int64        `json:"timeout"`
	Nodelist        []DBNodeInfo `json:"nodelist"`
	Password        string       `json:"password"`
	Port            int          `json:"port"`
	ConnectionLimit int          `json:"limit"`
}

type DBNodeInfo struct {
	Host     string `json:"host"`
	Id       string `json:"id"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

type DBNode struct {
	Client *formosa.Client
	Info   DBNodeInfo
	Id     string
}

type SrvData struct {
	Key   string
	Value string
}

type sortedMap struct {
	m map[string]string
	s []string
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[string]string) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Strings(sm.s)
	return sm.s
}

type sortedSrvArray []SrvData

func (sm sortedSrvArray) Len() int {
	return len(sm)
}

func (sm sortedSrvArray) Less(i, j int) bool {
	return sm[i].Key > sm[j].Key
}

func (sm sortedSrvArray) Swap(i, j int) {
	sm[i], sm[j] = sm[j], sm[i]
}

func sortedSrvKeys(m []SrvData) []SrvData {
	var sm sortedSrvArray = m
	sort.Sort(sm)
	return sm
}

type sortedRSrvArray []SrvData

func (sm sortedRSrvArray) Len() int {
	return len(sm)
}

func (sm sortedRSrvArray) Less(i, j int) bool {
	return sm[i].Key < sm[j].Key
}

func (sm sortedRSrvArray) Swap(i, j int) {
	sm[i], sm[j] = sm[j], sm[i]
}

func sortedSrvRKeys(m []SrvData) []SrvData {
	var sm sortedRSrvArray = m
	sort.Sort(sm)
	return sm
}
