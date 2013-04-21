package main

import (
	"net/http"
	"eyefi"
	"flag"
	"fmt"
	"os"
)

func main() {

	var conf = flag.String("c", "~/.eyego.conf", "EyeGo configuration file")
	flag.Parse()

	eyefi.ConfigFrom(*conf)

	http.HandleFunc("/", eyefi.Handler)
	http.ListenAndServe(":59278", nil)
}
