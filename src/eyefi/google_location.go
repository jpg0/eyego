package eyefi

import (
	"net/http"
	"fmt"
	"encoding/json"
	"bytes"
	"io/ioutil"
)

var (
	LOC_BASE_URL = "https://www.googleapis.com/geolocation/v1/geolocate"
)

type WifiAccessPoints struct {
	AccessPoints []AccessPointSightingInfo `json:"wifiAccessPoints"`
}

type Location struct {
	Latitude string `json:"lat"`
	Longitude string `json:"lon"`
}

type LocationResult struct {
	Location Location `json:"location"`
	Accuracy string `json:"accuracy"`

}


func GPSCoordinates(aps []AccessPointSightingInfo) (lr LocationResult, err error) {

	req, err := json.Marshal(WifiAccessPoints{AccessPoints:aps})
	if err != nil { return }

	resp, err := http.Post(fmt.Sprintf("%s?key=%s", LOC_BASE_URL, Config().GoogleAPIKey), "application/json", bytes.NewBuffer(req))
	if err != nil { return }

	defer resp.Body.Close()

	lr = LocationResult{}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil { return }

	err = json.Unmarshal(bytes, &lr)
	if err != nil { return }

	return
}
