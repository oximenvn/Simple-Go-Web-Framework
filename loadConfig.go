package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
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
