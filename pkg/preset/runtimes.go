package preset

import (
	"dalamud-updater/pkg/util"
	"fmt"
)

func GetRuntimes(version string) []string {
	return []string{
		fmt.Sprintf("shared/Microsoft.NETCore.App/%s", version),
		util.GetDotNetRuntimeUrl(version),
		fmt.Sprintf("shared/Microsoft.WindowsDesktop.App/%s", version),
		util.GetDotNetWDRuntimeUrl(version),
	}
}
