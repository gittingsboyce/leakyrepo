package main

import (
	"github.com/lgboyce/leakyrepo/cmd"
)

// These will be set during build via ldflags
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Set version info in cmd package
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}

