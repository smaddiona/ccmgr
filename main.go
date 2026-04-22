package main

import (
	"os"

	"github.com/smaddiona/ccmgr/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
