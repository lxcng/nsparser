package app

import (
	"log"
	"nsparser/internal/adapter"
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
				Usage:   "delete show",
				Action:  res.del,
			},
			{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "delete show",
				Action:  res.remove,
			},
			// {
			// 	Name:    "list",
			// 	Aliases: []string{"l"},
			// 	Usage:   "list torrents",
			// 	Action:  res.list,
			// },
			{
				Name:    "on",
				Aliases: []string{"on"},
				Usage:   "start transmission",
				Action: func(ctx *cli.Context) error {
					return adapter.Start()
				},
			},
			{
				Name:    "off",
				Aliases: []string{"off"},
				Usage:   "stop transmission",
				Action: func(ctx *cli.Context) error {
					return adapter.Stop()
				},
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

func (x *App) remove(*cli.Context) error {
	return adapter.Flush()
}

func (x *App) list(*cli.Context) error {
	// TODO list
	return nil
}
