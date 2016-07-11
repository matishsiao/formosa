package formosa

import (
	"log"
	"os"
	"runtime"
	"time"
)

var (
	Version      string = "0.0.1"
	ConfigPath   string = "configs.json"
	CONFIGS      Configs
	modTime      time.Time
	GlobalClient ServerClient
	DM           DBManager
	Binlogs      DBManager
	memprofile   string = "mempprof.log"
	memFile      *os.File
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
	DM.Construct("data")
	Binlogs.Construct("meta")
	go Listen(CONFIGS.Host, CONFIGS.Port)
	go WebServer()
	timeCounter := 0
	timePrint := 240
	if CONFIGS.Debug {
		timePrint = 20
	}
	for {
		configWatcher()
		//one min ping mirror DBs
		timeCounter++
		if timeCounter%timePrint == 0 {
			//GlobalClient.DBPool.CheckStatus()
			//GlobalClient.DBPool.Status()
			PrintGCSummary()
			timeCounter = 0
		}
		time.Sleep(250 * time.Millisecond)
	}
}
