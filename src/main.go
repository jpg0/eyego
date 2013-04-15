package main

import (
	"net/http"
	"eyefi"
	"flag"
//	"os"
)

func main() {

//	f, err := os.OpenFile("/tmp/in.tar", os.O_RDONLY, 0600)
//
//	_, _, err = eyefi.WriteFiles(f)
//
//	if err != nil {
//		panic(err)
//	}
//
//	return

	var conf = flag.String("c", "~/.eyego.conf", "EyeGo configuration file")
	flag.Parse()

	config, err := eyefi.ConfigFrom(*conf)

	if err != nil {
		panic(err)
	}

	eyefi.AddCardConfigs(config.Cards)

	http.HandleFunc("/", eyefi.Handler)
	http.ListenAndServe(":59279", nil)
}
