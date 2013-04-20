package main

import (
	"net/http"
	"eyefi"
	"flag"
	"os"
	"fmt"
)

func main() {

	f, _ := os.OpenFile("/tmp/eyego/P1050237.JPG.log", os.O_RDONLY, 0)
	p, _ := eyefi.ParseLog(f)

	fmt.Println(p.AccessPoints("P1050237.JPG"))

	return

	var conf = flag.String("c", "~/.eyego.conf", "EyeGo configuration file")
	flag.Parse()

	eyefi.ConfigFrom(*conf)

	http.HandleFunc("/", eyefi.Handler)
	http.ListenAndServe(":59279", nil)
}
