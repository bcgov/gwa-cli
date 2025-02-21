package main

import (
	"os"

	"github.com/bcgov/gwa-cli/cmd"
	"github.com/bcgov/gwa-cli/pkg"
)

var ApiHost string
var ApiVersion string
var ClientId string
var DefaultOrg string
var DefaultOrgUnit string
var Version string

func main() {
	cwd, _ := os.Getwd()
	ctx := &pkg.AppContext{
		ApiHost:        ApiHost,
		ApiVersion:     ApiVersion,
		ClientId:       ClientId,
		Cwd:            cwd,
		DefaultOrg:     DefaultOrg,
		DefaultOrgUnit: DefaultOrgUnit,
		Version:        Version,
	}
	cmd.Execute(ctx)
}
