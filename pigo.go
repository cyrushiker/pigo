package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/cyrushiker/pigo/cmd"
	"github.com/cyrushiker/pigo/pkg/setting"
)

const APP_VER = "19.07.0310"

func init() {
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "Pigo"
	app.Usage = "Pigo is something amazing"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.Web,
		cmd.Index,
		cmd.Db,
		cmd.Crawl,
	}
	app.Run(os.Args)
}
