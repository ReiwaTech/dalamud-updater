package util

import "fmt"

const distrib = "https://github.com/ReiwaTech/dalamud-distrib/raw/main"
const dotnet = "https://dotnetcli.azureedge.net/dotnet"

const DalamudAssets = "https://raw.githubusercontent.com/goatcorp/DalamudAssets/master/asset.json"
const NotoSansCJKscRegular = "https://github.com/notofonts/noto-cjk/raw/main/Sans/OTF/SimplifiedChinese/NotoSansCJKsc-Regular.otf"
const NotoSansCJKscMedium = "https://github.com/notofonts/noto-cjk/raw/main/Sans/OTF/SimplifiedChinese/NotoSansCJKsc-Medium.otf"

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
