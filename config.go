package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	//"os/signal"
	"sync"
	//"syscall"
)

type config struct {
	database struct {
		host   string
		user   string
		pass   string
		name   string
		driver string
	}
	app struct {
		ENV string
	}
}

var (
	sysConfig  config
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

	var temp config
	log.Println(file)
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

func GetConfig() config {
	configLock.RLock()
	defer configLock.RUnlock()
	return sysConfig
}
