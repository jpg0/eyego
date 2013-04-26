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
//test()
	http.HandleFunc("/", eyego.Handler)
	eyego.Info("Server starting...")
	err = http.ListenAndServe(":59278", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func test() {
	soapString := `<?xml version="1.0" encoding="UTF-8"?><SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns1="EyeFi/SOAP/EyeFilm"><SOAP-ENV:Body><ns1:UploadPhoto><fileid>1</fileid><macaddress>001856433ca9</macaddress><filename>P1050237.JPG.tar</filename><filesize>5599232</filesize><filesignature>353100003000040000000000e0110300</filesignature><encryption>none</encryption><flags>4</flags></ns1:UploadPhoto></SOAP-ENV:Body></SOAP-ENV:Envelope>`
	soap := new(eyego.UploadPhoto)
	eyego.ParseSoap(soapString, soap)

	eyego.Info("Uploading %s", soap.Filename)

	soap.MacAddress = "001856433ca9"
	soap.Filename = "P1050237.JPG.tar"
	soap.FileSize = "5599232"
	soap.FileSignature = "353100003000040000000000e0110300"
	soap.Encryption = "none"
	soap.Flags = "4"
	panic(eyego.CreateSoap(soap))

	panic("ok")

}
