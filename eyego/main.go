package main

import (
	"net/http"
	"eyego"
	"flag"
)

func main() {

	var conf = flag.String("c", "~/.eyego.conf", "EyeGo configuration file")
	flag.Parse()

	eyego.ConfigFrom(*conf)

	http.HandleFunc("/", eyego.Handler)
	http.ListenAndServe(":59278", nil)
}
