package cmd

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/cyrushiker/pigo/models"
)

var (
	Crawl = cli.Command{
		Name:        "crawl",
		Usage:       "scrapy things from internet",
		Description: `You can choose what you are want.`,
		Subcommands: []cli.Command{
			subcmdCrawlNovel,
		},
	}

	subcmdCrawlNovel = cli.Command{
		Name:   "novel",
		Usage:  "scrapy novels",
		Action: runCrawlNovels,
		Flags: []cli.Flag{
			boolFlag("proxy, p", "whether use proxy or not."),
			stringFlag("book, b", "", "set the book number to scratch"),
		},
	}
)

func runCrawlNovels(c *cli.Context) error {
	proxy := c.Bool("proxy")
	book := c.String("book")
	if err := models.NovelScrapy(proxy, book); err != nil {
		fmt.Fprintf(c.App.Writer, "crawl novel: %v", err)
	}

	return nil
}
