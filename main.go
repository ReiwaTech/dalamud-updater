package main

import (
	"dalamud-updater/pkg/api"
	"dalamud-updater/pkg/preset"
	"dalamud-updater/pkg/util"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cavaliergopher/grab/v3"
	"github.com/spf13/viper"
)

func getLocalAssetsVersion(dir string) int {
	file := path.Join(dir, "asset.ver")
	return util.ReadNumber(file)
}

const banner = `
 _           ___         | ReiwaTech
|_) _ o     _.| _  _|_   |  Dalamud
| \(/_|\/\/(_||(/_(_| |  |   Updater
`

func main() {
	fmt.Println(banner)

	dir := util.GetWorkingDir()

	// init config
	viper.SetConfigName("dalamud")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			// fmt.Println("Config file not found, skip")
			if !util.Confirm(fmt.Sprintf("Working in: %s", dir), true) {
				fmt.Println("No - Exiting")
				return
			}

			viper.SafeWriteConfig()
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	// check update
	channel := viper.GetString("Channel")
	if channel == "" {
		channel = "cn"
		viper.Set("Channel", channel)
	}

	localVersion := viper.GetString("Version")
	fmt.Println("Requesting Github for latest version ...")
	fmt.Println("  Channel         :", channel)
	fmt.Println("  Local Version   :", localVersion)
	latest := api.GetVersion(channel)

	fmt.Println("  Dalamud version :", latest.AssemblyVersion)
	fmt.Println("  Runtime version :", latest.RuntimeVersion)

	downloadDir := path.Join(dir, "download")
	os.MkdirAll(downloadDir, os.ModePerm)

	// update dalamud
	if localVersion == latest.AssemblyVersion {
		fmt.Println("Local version is already latest")
	} else {
		// download latest file from github
		fmt.Println("Downloading from Github ...")
		fileName := fmt.Sprintf("dalamud-%s-%s.zip", channel, latest.AssemblyVersion)
		resp, err := util.Download(path.Join(downloadDir, fileName), util.GetArchiveUrl(channel))

		if err != nil {
			panic(err)
		}

		fmt.Println("Download saved to", resp.Filename)
		fmt.Println("Unziping ...")

		if err := util.Unzip(resp.Filename, path.Join(dir, "Release")); err != nil {
			panic(err)
		}

		// save current version
		viper.Set("Version", latest.AssemblyVersion)
		if err := viper.WriteConfig(); err != nil {
			panic(err)
		}
	}

	// update runtime
	fmt.Println("Checking runtimes ...")
	runtimes := preset.GetRuntimes(latest.RuntimeVersion)

	for i := 0; i < len(runtimes); i += 2 {
		if util.IsDir(path.Join(dir, "runtime", runtimes[i])) {
			continue
		}

		resp, err := util.Download(downloadDir, runtimes[i+1])
		if err != nil {
			panic(err)
		}

		fmt.Println("Download saved to", resp.Filename)
		fmt.Println("Unziping ...")

		if err := util.Unzip(resp.Filename, path.Join(dir, "runtime")); err != nil {
			panic(err)
		}
	}

	// update assets
	fmt.Println("Requesting Github for assets version ...")
	latestAssets := api.GetAssetsVersion()

	assetsDir := path.Join(dir, "XIVLauncher/dalamudAssets/dev")
	assetsVerFile := path.Join(assetsDir, "asset.ver")
	assetsVer := util.ReadNumber(assetsVerFile)

	fmt.Println("  Local version   :", assetsVer)
	fmt.Println("  Latest version  :", latestAssets.Version)

	if assetsVer == latestAssets.Version {
		fmt.Println("Local assets version is already latest")
	} else {
		fmt.Println("Downloading from Github ...")

		reqs := []*grab.Request{}
	assets:
		for i := 0; i < len(latestAssets.Assets); i++ {
			row := latestAssets.Assets[i]
			filename := path.Join(assetsDir, path.Clean(row.FileName))
			os.MkdirAll(filepath.Dir(filename), os.ModePerm)

			// check overrides
			for j := 0; j < len(preset.AssetsOverride); j += 2 {
				rule := preset.AssetsOverride[j]
				override := preset.AssetsOverride[j+1]

				if rule == row.FileName {
					fmt.Printf("Overriding %s with %s\n", rule, override)

					if strings.HasPrefix(override, "https://") {
						row.URL = override
						row.Hash = ""
						break
					}

					util.WriteString(filename, override)
					continue assets
				}
			}

			if row.Hash != "" {
				hash, err := util.Sha1(filename)
				if err == nil && hash == strings.ToLower(row.Hash) {
					continue
				}
			}

			os.Remove(filename)

			req, err := grab.NewRequest(filename, row.URL)
			if err != nil {
				panic(err)
			}

			req.NoResume = true
			reqs = append(reqs, req)
		}

		if len(reqs) > 0 {
			util.DownloadBatch(10, reqs...)
		}

		util.WriteNumber(assetsVerFile, latestAssets.Version)
	}

	// update fonts
	fmt.Println("Checking fonts ...")
	installedFonts := viper.GetStringSlice("fonts")

	for i := 0; i < len(preset.Fonts); i += 2 {
		font := preset.Fonts[i]
		fileName := path.Join(assetsDir, font)
		if util.Contains[string](installedFonts, font) {
			continue
		}

		if util.IsFile(fileName) {
			os.Remove(fileName)
		}

		resp, err := util.Download(fileName, preset.Fonts[i+1])
		if err != nil {
			panic(err)
		}

		fmt.Println("Download saved to", resp.Filename)
		installedFonts = append(installedFonts, font)
		viper.Set("fonts", installedFonts)
		viper.WriteConfig()
	}

	fmt.Println()
	fmt.Println("All done!")
}
