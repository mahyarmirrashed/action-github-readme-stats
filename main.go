package main

import (
	"github.com/mahyarmirrashed/github-readme-stats/cmd"
	_ "golang.org/x/crypto/x509roots/fallback" // CA bundle for FROM Scratch
)

func main() {
	cmd.Execute()
}
