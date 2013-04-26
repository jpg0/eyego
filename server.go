package eyego

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	var response string

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
			response = doStartSession(*soap)
		case "GetPhotoStatus":
			soap := new(GetPhotoStatus)
			ParseSoap(body, soap)
			response = doGetPhotoStatus(*soap)
		case "MarkLastPhotoInRoll":
			soap := new(MarkLastPhotoInRoll)
			ParseSoap(body, soap)
			response = CreateSoap(MarkLastPhotoInRollResponse{})
		default:
			panic("unknown soap format")
		}
	case "/api/soap/eyefilm/v1/upload":
		Info("Upload started...")
		var err error
		response, err = doPhotoUpload(r)
		if err != nil {
			panic(fmt.Sprintf("Upload failed %s", err))
		}
		Info("Upload complete.")
	}

	w.Header().Add("Content-Length", fmt.Sprint(len(response)))
	fmt.Fprintf(w, response)
}

func doGetPhotoStatus(body GetPhotoStatus) string {
	return CreateSoap(
		GetPhotoStatusResponse{
			FileID:"1",
			Offset:"0"})
}

func doStartSession(body SoapStartSession) string {

	card := GetCard(body.MacAddress)

	return CreateSoap(
		SoapStartSessionResponse{
			Credential:card.Credential(body.CNonce),
			SNonce:GenerateSNonce(),
			TransferMode: "2",
			TransferModeTimestamp: body.TransferModeTimestamp,
			UpSyncAllowed:"false"})
}
