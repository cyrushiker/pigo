package cmd

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/cyrushiker/pigo/models"
	"github.com/cyrushiker/pigo/pkg/setting"
)

var (
	Index = cli.Command{
		Name:        "index",
		Usage:       "Create of rebuild index in Elasticsearch",
		Description: `You may need to take care of this command.`,
		Subcommands: []cli.Command{
			subcmdCreateIndices,
		},
	}

	subcmdCreateIndices = cli.Command{
		Name:   "create",
		Usage:  "Create index in Elastic",
		Action: runCreateIndices,
		Flags: []cli.Flag{
			stringFlag("config, c", "conf/app.yml", "Custom configuration file path"),
			boolFlag("delete, d", "Delete the index if it is already exixts."),
		},
	}
)

func runCreateIndices(c *cli.Context) error {
	if c.IsSet("config") {
		setting.CustomConf = c.String("config")
	}

	setting.NewContext()
	models.NewEsCli()

	flag := c.Bool("delete")
	if err := models.CreateIndex(flag); err != nil {
		fmt.Fprintf(c.App.Writer, "CreateIndices: %v", err)
	}

	return nil
}
