package main

import (
	"net/http"
	"eyefi"
	"flag"
)

func main() {

	var conf = flag.String("c", "~/.eyego.conf", "EyeGo configuration file")
	flag.Parse()

	config, err := eyefi.ConfigFrom(*conf)

	if err != nil {
		panic(err)
	}

	eyefi.AddCardConfigs(config.Cards)

	http.HandleFunc("/", eyefi.Handler)
	http.ListenAndServe(":59278", nil)
}
