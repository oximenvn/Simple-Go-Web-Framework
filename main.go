// main.go
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"./core"
)

func main() {
	initConfig()
	http.HandleFunc("/", welcome)
	http.ListenAndServe(":8080", nil)
}

func init() {
	logFile, error := os.OpenFile(time.Now().Format("20060102")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	core.Check(error)
	//defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.Println("Started. PID:" + strconv.Itoa(os.Getpid()))
	reloadCmd := flag.NewFlagSet("reload", flag.ExitOnError)
	pidServer := reloadCmd.Int("pid", 0, "Request server reload configuration")

	flag.Parse()
	//log.Println(os.Args)

	if len(os.Args) == 1 {
		return
	}

	switch os.Args[1] {
	case "reload":
		reloadCmd.Parse(os.Args[2:])
		//log.Println("reload ", reloadCmd.Args())
		sendSignal(*pidServer)
		os.Exit(1)
	case "h":
	case "help":
		printHelp()
		os.Exit(1)
	default:
		printHelp()
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func sendSignal(pid int) {
	log.Println(pid, syscall.SIGUSR2)
	syscall.Kill(pid, syscall.SIGUSR2)
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
			log.Println(sysConfig)
		}
	}()
}

func printHelp() {
	fmt.Println("Usage of " + os.Args[0] + ":")
	fmt.Println("The most commonly used commands are:")
	fmt.Println(" reload -pid=<PID>   Request server reload configuration")
}
