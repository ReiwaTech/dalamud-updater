package util

import "fmt"

const distrib = "https://reiwatech.github.io/dalamud-distrib"
const dotnet = "https://dotnetcli.azureedge.net/dotnet"

const DalamudAssets = "https://goatcorp.github.io/DalamudAssets/asset.json"
const NotoSansCJKscRegular = "https://mirrors.tuna.tsinghua.edu.cn/ctan/fonts/notocjksc/NotoSansCJKsc-Regular.otf"
const NotoSansCJKscMedium = "https://mirrors.tuna.tsinghua.edu.cn/ctan/fonts/notocjksc/NotoSansCJKsc-Medium.otf"

func GetVersionUrl(channel string) string {
	return fmt.Sprintf("%s/%s/version", distrib, channel)
}

func GetArchiveUrl(channel string) string {
	return fmt.Sprintf("%s/%s/latest.zip", distrib, channel)
}

func GetDotNetRuntimeUrl(version string) string {
	return fmt.Sprintf("%s/Runtime/%s/dotnet-runtime-%s-win-x64.zip", dotnet, version, version)
}

func GetDotNetWDRuntimeUrl(version string) string {
	return fmt.Sprintf("%s/WindowsDesktop/%s/windowsdesktop-runtime-%s-win-x64.zip", dotnet, version, version)
}
