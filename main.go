package main

import (
	"sigolang/cmd"
)

var AppVersion string = "local"

func main() {
	cmd.Execute(AppVersion)
}
