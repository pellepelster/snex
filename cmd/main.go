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
	app := &cli.App{
		Name:  "snex",
		Usage: "snex keep the code snippets inside your documentation in sync with real code from your sources",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "template",
				Usage: fmt.Sprintf("set custom snippet template to use for replacements, available variables are:\n%s", pkg.TemplateHelp),
			},
			&cli.BoolFlag{
				Name:  "show-templates",
				Usage: fmt.Sprintf("show default templates for snippet replacing"),
			},
		},
		Action: func(context *cli.Context) error {

			if context.IsSet("show-templates") {
				for _, template := range pkg.DefaultSnippetTemplates {
					log.Infof("template for file extension(s) %s: '%s'", strings.Join(template.Extensions, ", "), strings.ReplaceAll(template.Template, "\n", "\\n"))
				}
				os.Exit(4)
			}

			if context.IsSet("template") {
				err := pkg.ValidateTemplate(context.String("template"))
				if err != nil {
					log.Fatalf("validating the template failed: %s", err)
				}
			}

			if context.NArg() == 0 {
				return cli.Exit("no source folders provided", 2)
			}

			for _, folderOrFile := range context.Args().Slice() {
				if !fileOrDirExists(folderOrFile) {
					return cli.Exit(fmt.Sprintf("folder or file '%s' not found", folderOrFile), 3)
				}
			}

			return processFiles(context.Args().Slice(), context.String("template"))
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
