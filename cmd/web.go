package cmd

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"

	"github.com/cyrushiker/pigo/pkg/setting"
	"github.com/cyrushiker/pigo/routes"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "Start web server",
	Description: `You will need to set host and ip first before run, otherwise it will runing with default setting.`,
	Action:      runWeb,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "port, p",
			Value: "9090",
			Usage: "Temporary port number to prevent conflict"},
		cli.StringFlag{
			Name:  "config, c",
			Value: "conf/app.yml",
			Usage: "Custom configuration file path"},
	},
}

func newGin() *gin.Engine {
	r := gin.Default()
	// todo: middleware config
	routes.RegHome(r.Group("/"))
	routes.RegUser(r.Group("/user"))
	return r
}

func runWeb(c *cli.Context) {
	if c.IsSet("config") {
		setting.CustomConf = c.String("config")
	}
	if c.IsSet("port") {
		setting.HTTPPort = c.String("port")
	}
	listenAddr := fmt.Sprintf("%s:%s", setting.HTTPAddr, setting.HTTPPort)
	routes.GlobalInit()

	r := newGin()
	r.Run(listenAddr)
}
