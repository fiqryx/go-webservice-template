package cmd

import (
	"fmt"
	"runtime"

	c "template.go/packages/common"
)

var Version = c.Env("VERSION", "1.0.0")

type VersionInfo struct {
	Show     bool   `json:"-" yaml:"-"`
	Version  string `json:"version"`
	Build    string `json:"build"`
	Go       string `json:"go"`
	Compiler string `json:"compiler"`
	Platform string `json:"platform"`
}

func GetVersion() *VersionInfo {
	return &VersionInfo{
		Show:     true,
		Version:  Version,
		Go:       runtime.Version(),
		Compiler: runtime.Compiler,
		Platform: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
