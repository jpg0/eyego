package main

import (
	"net/http"
	"eyefi"
	"flag"
)

func main() {

//
////	_, _, err = eyefi.WriteFiles(f)
//
//	if err != nil {
//		panic(err)
//	}
//
//	var block []byte
//	block = make([]byte, 512)
//
//	n, err := f.Read(block)
//
//	if n != 512 { panic(n) }
//
//	fmt.Println(eyefi.tcp_checksum(block[0:512]))
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
