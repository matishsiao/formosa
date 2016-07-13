package formosa

import (
	"log"
	"os"
	"runtime"
	"time"
)

var (
	Version    string = "0.0.1"
	ConfigPath string = "configs.json"
	CONFIGS    Configs
	modTime    time.Time
	SyncClient ServerClient
	DM         DBManager
	memprofile string = "mempprof.log"
	memFile    *os.File
)

func Run(path string) {
	ConfigPath = path
	config, err := loadConfigs(ConfigPath)
	if err != nil {
		log.Println("Load config file error:", err)
		os.Exit(1)
	}
	CONFIGS = config
	SetUlimit(1002000)
	useCPU := runtime.NumCPU() - 1
	if useCPU <= 0 {
		useCPU = 1
	}
	runtime.GOMAXPROCS(useCPU)
	DM.Construct(CONFIGS.DBPath)
	SyncClient.Init()
	go Listen(CONFIGS.Host, CONFIGS.Port)
	go WebServer()
	go ConfigWatcher()

	//Main Goroutine for check other db status and db status
	for {
		SyncClient.ServerPool.CheckStatus()
		PrintGCSummary()
		time.Sleep(60 * time.Second)
	}
}
