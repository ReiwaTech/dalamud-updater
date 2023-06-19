package api

import (
	"dalamud-updater/pkg/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Version struct {
	AssemblyVersion  string `json:"AssemblyVersion"`
	SupportedGameVer string `json:"SupportedGameVer"`
	RuntimeVersion   string `json:"RuntimeVersion"`
	GitSha           string `json:"GitSha"`
	Revision         string `json:"Revision"`
}

func GetVersion(channel string) Version {
	resp, err := http.Get(util.GetVersionUrl(channel))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	var result Version
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	return result
}
