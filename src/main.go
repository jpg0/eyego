package main

import (
	"net/http"
	"eyefi"
	"flag"
	"os"
	"fmt"
)

func main() {

	var conf = flag.String("c", "~/.eyego.conf", "EyeGo configuration file")
	flag.Parse()

	eyefi.ConfigFrom(*conf)

	test()

	http.HandleFunc("/", eyefi.Handler)
	http.ListenAndServe(":59279", nil)
}

func test() {
	f, _ := os.OpenFile("/tmp/eyego/P1050237.JPG.log", os.O_RDONLY, 0)
	p, _ := eyefi.ParseLog(f)

	aps := p.AccessPoints("P1050237.JPG")

	fmt.Println(aps)

	res, err := eyefi.GPSCoordinates(aps)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)

	panic("test complete")
}
