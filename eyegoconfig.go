package eyego

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
)
type EyegoConfig struct {
	TargetDir string `json:"target_dir"`
	GoogleAPIKey string `json:"google_api_key"`
	Cards []Card `json:"cards"`
}


type Card struct {
	UploadKey string `json:"upload_key"`
	MacAddress string `json:"mac_address"`
}

var config EyegoConfig

func ConfigFrom(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		fmt.Print(err)
	}
}

func Config() EyegoConfig {
	return config
}

func DumpConfig(config EyegoConfig) {
	data, _ := json.Marshal(config)
	fmt.Printf("%s", data)
}

func GetCard(mac_address string) Card {
	for i := range config.Cards {
		if config.Cards[i].MacAddress == mac_address {
			return config.Cards[i]
		}
	}

	panic("no such card")
}
