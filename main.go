// main.go
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/oximenvn/Simple-Go-Web-Framework/core"
)

var version = 0.01

func main() {
	//core.initConfig()
	http.Handle("/static/", //final url can be anything
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static/")))) //Go looks in the relative static directory first, then matches it to a
	//url of our choice as shown in http.Handle("/static/").
	//This url is what we need when referencing our css files
	//once the server begins. Our html code would therefore be <link rel="stylesheet"  href="/static/stylesheet/...">
	//It is important to note the final url can be whatever we like, so long as we are consistent.
	http.HandleFunc("/", core.Routing)
	http.ListenAndServe(":8000", nil)
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
		core.SendSignal(*pidServer)
		os.Exit(1)
	case "migrate":
		core.Migrate(Tables{})
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

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <command>\n", os.Args[0])
	fmt.Fprintln(flag.CommandLine.Output(), "\nCommands:")
	fmt.Fprintln(flag.CommandLine.Output(), "  \033[1mreload\033[0m -pid=<PID>   Request server reload configuration")
	fmt.Fprintln(flag.CommandLine.Output(), "  \033[1mmigrate\033[0m -pid=<PID>  Migrate database")
	fmt.Fprintln(flag.CommandLine.Output(), "  \033[1mhelp\033[0m                Display this help and exit")
	fmt.Fprintln(flag.CommandLine.Output(), "  \033[1mversion\033[0m             Display version information.")
	flag.PrintDefaults()
}
