package eyefi

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"promise"
	"encoding/xml"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	body_bytes, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body_bytes))

	body := ParseSoap(string(body_bytes))

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

	fmt.Println(rv.Get().(string))


	fmt.Fprintf(w, rv.Get().(string))
}

func doStartSession(body SoapStartSession) promise.Promise {

	card := GetCard(body.MacAddress)

	rv := promise.NewEagerPromise(func() interface{} {
		return CreateSoap(
		SoapBody{
			StartSessionResponse:&SoapStartSessionResponse{
				Credential:card.Credential(body.CNonce),
				SNonce:GenerateSNonce(),
				TransferMode: body.TransferMode,
				TransferModeTimestamp: body.TransferModeTimestamp,
				UpSyncAllowed:"false",
				XMLName:xml.Name{
					Space:"http://localhost/api/soap/eyefilm",
					Local:"StartSessionResponse"}}})
	})
	return rv
}
