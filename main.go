// main.go
package main

import (
	//"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	log.Println("Starting...\nPID:" + strconv.Itoa(os.Getpid()))
	initConfig()
	http.HandleFunc("/", welcome)
	http.ListenAndServe(":8080", nil)
}

func initConfig() {
	loadConfig(true)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR2)
	go func() {
		for {
			<-s
			loadConfig(false)
			log.Println("Reloaded")
		}
	}()
}
