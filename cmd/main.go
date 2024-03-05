package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/pellepelster/snex/pkg"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func main() {
	log.SetReportTimestamp(false)

	app := &cli.App{
		Name:           "snex",
		Usage:          "snex keep the code snippets inside your documentation in sync with real code from your sources",
		DefaultCommand: "replace",
		Commands: []*cli.Command{
			{
				Name:  "show-templates",
				Usage: "show default templates for snippet replacements",
				Action: func(cCtx *cli.Context) error {
					for _, template := range pkg.DefaultSnippetTemplates {
						log.Infof("template for file extension(s) %s: '%s'", strings.Join(template.Extensions, ", "), strings.ReplaceAll(template.Template, "\n", "\\n"))
					}

					return cli.Exit("", 4)
				},
			},
			{
				Name:      "replace",
				Usage:     "replace snippets in all source folders and files",
				ArgsUsage: "[source folders or files...]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "template",
						Usage: fmt.Sprintf("set custom snippet template to use for replacements, available variables are:\n%s", pkg.TemplateHelp),
					},
				},
				Action: func(context *cli.Context) error {
					if context.IsSet("template") {
						err := pkg.ValidateTemplate(context.String("template"))
						if err != nil {
							return cli.Exit(fmt.Sprintf("validating the template failed: %s", err), 2)
						}
					}

					if context.NArg() == 0 {
						return cli.Exit("no source folders provided", 3)
					}

					for _, folderOrFile := range context.Args().Slice() {
						if !fileOrDirExists(folderOrFile) {
							return cli.Exit(fmt.Sprintf("folder or file '%s' not found", folderOrFile), 5)
						}
					}

					return processFiles(context.Args().Slice(), context.String("template"))
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
