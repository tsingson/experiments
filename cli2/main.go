package main

import (
	"fmt"
	"github.com/tsingson/fastweb/fasthttputils"
	"gopkg.in/urfave/cli.v2"
	"os"
)

func main() {
	path, _ := fasthttputils.GetCurrentExecDir()
	app := &cli.App{
		Name:        "greet",
		Version:     "0.1.0",
		Description: "This is how we describe greet the app",
		Authors: []*cli.Author{
			{Name: "Harrison", Email: "harrison@lolwut.com"},
			{Name: "Oliver Allen", Email: "oliver@toyshop.com"},
		},
		Flags: []cli.Flag{
			//	&cli.StringFlag{Name: "name", Value: "bob", Usage: "a name to say"},
		}}
	app.Commands = []*cli.Command{
		{
			Name:    "hello",
			Aliases: []string{"hi"},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "name",
					Value: path,
					Usage: "Name of the person to greet",
				},
			},
			Usage:       "use it to see a description",
			Description: "This is how we describe hello the function",
			//	Subcommands: []*cli.Command{
			//		{
			/**
			Name:        "english",
			Aliases:     []string{"en"},
			Usage:       "sends a greeting in english",
			Description: "greets someone in english",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "name",
					Value: "Bob",
					Usage: "Name of the person to greet",
				},
			},
			*/
			Action: func(c *cli.Context) error {
				fmt.Println("Hello,", c.String("name"))
				return nil
				//		},
				//	},
			},
		},
	}

	app.Run(os.Args)
}
