package app

import (
	"log"
	"nsparser/internal/config"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
)

type App struct {
	c    *cli.App
	conf *config.Config
}

func Init(conf *config.Config) *App {

	res := &App{conf: conf}
	c := &cli.App{
		Name:                 "nsparser",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add show",
				Action:  res.add,
			},
			{
				Name:    "all",
				Aliases: []string{"s"},
				Usage:   "start all torrents",
				Action: func(*cli.Context) error {
					return conf.StartAll()
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "start all torrents",
				Action:  res.del,
			},
		},
	}
	res.c = c
	return res
}

func (x *App) Run() {
	if err := x.c.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func (x *App) add(*cli.Context) error {
	prompt := promptui.Select{Label: "Select translator", Items: x.conf.GetTranslators()}
	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	prompt2 := promptui.Prompt{Label: "Title"}
	title, err := prompt2.Run()
	if err != nil {
		return err
	}

	prompt3 := promptui.Prompt{Label: "Present"}
	present, err := prompt3.Run()
	if err != nil {
		return err
	}

	return x.conf.AddShow(i, title, present)
}

func (x *App) del(*cli.Context) error {
	prompt := promptui.Select{Label: "Select translator", Items: x.conf.GetTranslators()}
	i, _, err := prompt.Run()
	if err != nil {
		return err
	}
	prompt2 := promptui.Select{Label: "Select translator", Items: x.conf.GetShows(i)}
	j, _, err := prompt2.Run()
	if err != nil {
		return err
	}

	return x.conf.DeleteShow(i, j)
}
