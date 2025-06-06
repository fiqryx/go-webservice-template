package cmd

import (
	"fmt"
	"runtime"
)

var (
	version = "1.0.0"
	build   = "2025-06-06T00:00:00Z"
)

type Version struct {
	Show     bool   `json:"-" yaml:"-"`
	Version  string `json:"version" yaml:"version"`
	Build    string `json:"build" yaml:"build"`
	Go       string `json:"go" yaml:"go"`
	Compiler string `json:"compiler" yaml:"compiler"`
	Platform string `json:"platform" yaml:"platform"`
}

func GetVersion() *Version {
	return &Version{
		Show:     true,
		Version:  version,
		Build:    build,
		Go:       runtime.Version(),
		Compiler: runtime.Compiler,
		Platform: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
