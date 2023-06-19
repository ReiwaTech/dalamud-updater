package api

import (
	"dalamud-updater/pkg/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type AssetsVersion struct {
	Version int `json:"Version"`
	Assets  []struct {
		URL      string `json:"Url"`
		FileName string `json:"FileName"`
		Hash     string `json:"Hash,omitempty"`
	} `json:"Assets"`
}

func GetAssetsVersion() AssetsVersion {
	resp, err := http.Get(util.DalamudAssets)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	var result AssetsVersion
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	return result
}
