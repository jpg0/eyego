package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	"encoding/xml"
	"soap"
	"promise"
	"eyefi"
	"encoding/hex"
)

func handler(w http.ResponseWriter, r *http.Request) {
	body_bytes, _ := ioutil.ReadAll(r.Body)
	body := soap.ParseSoap(string(body_bytes))

	actionHeader := r.Header.Get("SOAPAction")
	action := actionHeader[5:len(actionHeader) - 1]
	var rv promise.Promise

	fmt.Print(action)

	switch (action) {
	case "StartSession":
		rv = doStartSession(*body.StartSession)
	default:
		panic("unknown soap format")
	}

	r.Header.Set("Content-Type", "application/xml")

	//	fmt.Print(string(body))
	fmt.Fprintf(w, rv.Get().(string))
}

func doStartSession(body soap.SoapStartSession) promise.Promise {
	rv := promise.NewEagerPromise(func() interface{} {
		return soap.CreateSoap(
			soap.SoapBody{
				StartSessionResponse:&soap.SoapStartSessionResponse{
					Credential:"test",
					SNonce:"",
					TransferMode: body.TransferMode,
					TransferModeTimestamp: body.TransferModeTimestamp,
					UpSyncAllowed:"true"}})
	})
	return rv
}

func main() {

	eyefi.AddCard(eyefi.Card{mac_address: "t"})

	fmt.Print(hex.DecodeString("1234abcd"))
	fmt.Print(eyefi.Card{}.Credential("test"))

	//test()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":59278", nil)
}



func test() {
	s := "<?xml version=\"1.0\" encoding=\"UTF-8\"?><SOAP-ENV:Envelope xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:ns1=\"EyeFi/SOAP/EyeFilm\"><SOAP-ENV:Body><ns1:StartSession><macaddress>001856433ca9</macaddress><cnonce>01a90d5c197af66d41f95e08aaad740a</cnonce><transfermode>546</transfermode><transfermodetimestamp>1358227288</transfermodetimestamp></ns1:StartSession></SOAP-ENV:Body></SOAP-ENV:Envelope>"

	parser := xml.NewDecoder(bytes.NewBufferString(s))

	envelope := new(soap.SoapEnvelope)

	parser.DecodeElement(&envelope, nil)

	fmt.Print(envelope.Body.StartSession.MacAddress)
}
