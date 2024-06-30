package main

import (
	"cp-connect/pkg/lib"
)

var version string = "n/a" // Placeholder constant
var commit string = "unk"  // Placeholder constant

func main() {
	lib.HandleArguments(version, commit)
}
