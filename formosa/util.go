package formosa

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"syscall"
	"time"
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

func PaddingLeft(str string, length int) string {
	if len(str) >= length {
		return str
	} else {
		return PaddingLeft("0"+str, length)
	}
}

func memPorfile() {
	log.Println("pprof profile started.")
	StartCPUProfile()
	time.Sleep(300 * time.Second)
	DumpHeap()
	StopCPUProfile()
	log.Println("write pprof profile finished.")
}
func writeMemProfile() {
	pprof.WriteHeapProfile(memFile)
}

func ConfigWatcher() {
	for {
		configWatcher()
		time.Sleep(250 * time.Millisecond)
	}
}
func configWatcher() {
	file, err := os.Open(ConfigPath) // For read access.
	if err != nil {
		log.Println("configWatcher error:", err)
		return
	}
	info, err := file.Stat()
	if err != nil {
		log.Println("configWatcher error:", err)
		file.Close()
		return
	}
	file.Close()
	if modTime.Unix() == -62135596800 {
		log.Println("configWatcher init mod time")
		modTime = info.ModTime()
	}

	if info.ModTime() != modTime {
		log.Printf("Config file changed. Reolad config file.\n")
		modTime = info.ModTime()
		CONFIGS, err = loadConfigs(ConfigPath)
		if err != nil {
			log.Printf("configWatcher error:%v\n", err)
			return
		}
	}

}

func loadConfigs(fileName string) (Configs, error) {
	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Printf("Load config file error: %v\n", e)
		os.Exit(1)
	}

	var config Configs
	err := json.Unmarshal(file, &config)
	if err != nil {
		log.Printf("Config load error:%v \n", err)
		return config, err
	}
	return config, nil
}
