package main

import (
	"fmt"
	"os"

	"github.com/FalcoSuessgott/kubectl-vault-login/cmd"
)

var version string

func main() {
	cmd.Version = version

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v.\n", err)

		os.Exit(1)
	}
}
