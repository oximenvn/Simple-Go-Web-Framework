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

var version = 0.01

func main() {
	initConfig()
	routing()
	http.ListenAndServe(":8080", nil)
}

func init() {
	logFile, error := os.OpenFile(time.Now().Format("20060102")+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	core.Check(error)
	//defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.Println("Started. PID:" + strconv.Itoa(os.Getpid()))

	flag.Usage = Usage
	reloadCmd := flag.NewFlagSet("reload", flag.ExitOnError)
	pidServer := reloadCmd.Int("pid", 0, "Request server reload configuration")
	//helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
	//versionCmd := flag.NewFlagSet("version", flag.ExitOnError)
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
	case "help":
		Usage()
		os.Exit(1)
	case "version":
		fmt.Fprintln(os.Stdout, "Version:"+strconv.FormatFloat(version, 'f', -1, 32))
		os.Exit(1)
	default:
		Usage()
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

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <command>\n", os.Args[0])
	fmt.Fprintln(flag.CommandLine.Output(), "\nCommands:")
	fmt.Fprintln(flag.CommandLine.Output(), "  \033[1mreload\033[0m -pid=<PID>   Request server reload configuration")
	fmt.Fprintln(flag.CommandLine.Output(), "  \033[1mhelp\033[0m                Display this help and exit")
	fmt.Fprintln(flag.CommandLine.Output(), "  \033[1mversion\033[0m             Display version information.")
	flag.PrintDefaults()
}
