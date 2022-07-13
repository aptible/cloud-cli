package main

import (
	"github.com/aptible/cloud-cli/cmd"
)

func main() {
	root := cmd.NewRootCmd()
	cmd.Execute(root)
}
