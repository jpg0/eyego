package eyefi

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"promise"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	body_bytes, _ := ioutil.ReadAll(r.Body)

	fmt.Println(string(body_bytes))

	body := string(body_bytes)

	actionHeader := r.Header.Get("SOAPAction")
	action := actionHeader[5:len(actionHeader) - 1]
	var rv promise.Promise

	fmt.Print(action)

	switch (action) {
	case "StartSession":
		soap := new(SoapStartSession)
		ParseSoap(body, soap)
		rv = doStartSession(*soap)
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
			SoapStartSessionResponse{
				Credential:card.Credential(body.CNonce),
				SNonce:GenerateSNonce(),
				TransferMode: "2",
				TransferModeTimestamp: body.TransferModeTimestamp,
				UpSyncAllowed:"false"})})
	return rv
}
