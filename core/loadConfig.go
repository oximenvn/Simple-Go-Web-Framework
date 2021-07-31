package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	sysConfig  Config
	configLock = new(sync.RWMutex)
)

func loadConfig(fail bool) {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("open config: ", err)
		if fail {
			os.Exit(1)
		}
	}

	var temp Config
	if err = json.Unmarshal(file, &temp); err != nil {
		log.Println("parse config: ", err)
		if fail {
			os.Exit(1)
		}
	}
	configLock.Lock()
	sysConfig = temp
	configLock.Unlock()
}

func GetConfig() Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return sysConfig
}

func SendSignal(pid int) {
	log.Println(pid, syscall.SIGUSR2)
	syscall.Kill(pid, syscall.SIGUSR2)
}

func init() {
	loadConfig(true)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR2)
	go func() {
		for {
			<-s
			loadConfig(false)
			log.Println("Reloaded")
			log.Println(sysConfig)
		}
	}()
}
