package main

import (
	"embed"

	"github.com/cvhariharan/flowctl/cmd"
)

//go:embed all:site/build
//go:embed configs
//go:embed migrations
var staticFiles embed.FS

func main() {
	cmd.StaticFiles = staticFiles
	cmd.Execute()
}
