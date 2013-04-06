package eyefi

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
)

type CardConfig struct {
	UploadKey string `json:"upload_key"`
	MacAddress string `json:"mac_address"`
}

type EyegoConfig struct {
	Cards []CardConfig `json:"cards"`
}

func ConfigFrom(path string) (config EyegoConfig, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		fmt.Print("bad json ", err)
	}

	return config, nil
}

func DumpConfig(config EyegoConfig) {
	data, _ := json.Marshal(config)
	fmt.Printf("%s", data)
}
