package eyefi

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"promise"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	var rv promise.Promise

	switch (r.URL.Path) {
	case "/api/soap/eyefilm/v1":
		body_bytes, _ := ioutil.ReadAll(r.Body)
		body := string(body_bytes)
		actionHeader := r.Header.Get("SOAPAction")
		action := actionHeader[5:len(actionHeader) - 1]

		switch (action) {
		case "StartSession":
			soap := new(SoapStartSession)
			ParseSoap(body, soap)
			rv = doStartSession(*soap)
		case "GetPhotoStatus":
			soap := new(GetPhotoStatus)
			ParseSoap(body, soap)
			rv = doGetPhotoStatus(*soap)
		default:
			panic("unknown soap format")
		}
	case "/api/soap/eyefilm/v1/upload":
		var err error
		rv, err = doPhotoUpload(r)
		if err != nil {
			panic(fmt.Sprintf("Upload failed %s", err))
		}
	}

	responseString := rv.Get().(string)
	w.Header().Add("Content-Length", fmt.Sprint(len(responseString)))
	fmt.Fprintf(w, responseString)
}

func doGetPhotoStatus(body GetPhotoStatus) promise.Promise {
	rv := promise.NewEagerPromise(func() interface{} {
		return CreateSoap(
		GetPhotoStatusResponse{
			FileID:"1",
			Offset:"0"})})
	return rv
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
