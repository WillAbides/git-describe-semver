package main

import (
	"fmt"
	"os"

	"github.com/willabides/git-describe-semver/cmd"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {
	if err := cmd.Execute(cmd.FullVersion{Version: version, Commit: commit, Date: date, BuiltBy: builtBy}); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
