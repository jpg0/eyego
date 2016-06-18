package eyego

import (
	"net/http"
	"fmt"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"errors"
)

var (
	LOC_BASE_URL = "https://www.googleapis.com/geolocation/v1/geolocate"
)

type WifiAccessPoints struct {
	AccessPoints []AccessPointSightingInfo `json:"wifiAccessPoints"`
	ConsiderIp bool `json:"considerIp"`
}

type Location struct {
	Latitude float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type LocationResult struct {
	Location Location `json:"location"`
	Accuracy float64 `json:"accuracy"`

}

func GPSCoordinates(aps []AccessPointSightingInfo) (lr LocationResult, err error) {

	req, err := json.Marshal(WifiAccessPoints{AccessPoints:aps,ConsiderIp:false})
	if err != nil { return }

	resp, err := http.Post(fmt.Sprintf("%s?key=%s", LOC_BASE_URL, Config().GoogleAPIKey), "application/json", bytes.NewBuffer(req))
	if err != nil { return }

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		err = errors.New(fmt.Sprintf("Request failed: %s", body))
		return
	}

	defer resp.Body.Close()

	lr = LocationResult{}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil { return }

	err = json.Unmarshal(bytes, &lr)
	if err != nil { return }

	return
}
