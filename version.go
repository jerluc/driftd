package main

import (
	"fmt"
)

// Auto-filled by build
var Version string
var Commitish string

func PrintVersionInfo() {
	fmt.Printf("%s (%s)\n", Version, Commitish)
}
