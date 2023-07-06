package main

import (
	"os"

	"github.com/bcgov/gwa-cli/cmd"
	"github.com/bcgov/gwa-cli/pkg"
)

var ApiHost string
var ClientId string

func main() {
	cwd, _ := os.Getwd()
	ctx := &pkg.AppContext{
		ApiHost:  ApiHost,
		ClientId: ClientId,
		Cwd:      cwd,
	}
	cmd.Execute(ctx)
}
