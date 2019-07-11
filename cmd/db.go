package cmd

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/cyrushiker/pigo/models"
	"github.com/cyrushiker/pigo/pkg/setting"
)

var (
	Db = cli.Command{
		Name:        "db",
		Usage:       "Create of rebuild database in mysql",
		Description: `You must take care of this command, cause it will destroy all you datas in database.`,
		Subcommands: []cli.Command{
			subcmdCreateTables,
		},
	}

	subcmdCreateTables = cli.Command{
		Name:   "create",
		Usage:  "Create tables in mysql",
		Action: runCreateTables,
		Flags: []cli.Flag{
			stringFlag("config, c", "conf/app.yml", "Custom configuration file path"),
			boolFlag("drop, d", "Drop the table if exixts."),
		},
	}
)

func runCreateTables(c *cli.Context) error {
	if c.IsSet("config") {
		setting.CustomConf = c.String("config")
	}

	setting.NewContext()
	models.NewGormDB()

	flag := c.Bool("drop")
	if err := models.CreateTables(flag); err != nil {
		fmt.Fprintf(c.App.Writer, "CreateTables: %v", err)
	}

	return nil
}
