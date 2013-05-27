package main

import (
	"net/http"
	"eyego"
	"flag"
	"os"
	"io"
	"log"
)

func main() {

	var conf = flag.String("conf", "~/.eyego.conf", "Configuration file")
	var logfileName = flag.String("logfile", "-", "Log file")
	var loglevelString = flag.String("loglevel", "INFO", "Log level")
	flag.Parse()

	var logfile io.Writer
	var err error

	if *logfileName == "-" {
		logfile = os.Stdout
	} else {
		logfile, err = os.OpenFile(*logfileName, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
	}

	eyego.Init(logfile, eyego.LogLevelFromName(*loglevelString))
	eyego.ConfigFrom(*conf)
	http.HandleFunc("/", eyego.Handler)
	eyego.Info("Server starting...")
	err = http.ListenAndServe(":59278", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
