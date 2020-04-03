package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Commands = []*cli.Command {
		{
			Name:    "sandbox-user",
			Aliases: []string{"sandbox"},
			Usage:   "Create a sandbox user",
			Action:  createSandboxUserCmd,
			Flags: []cli.Flag {
				&cli.StringFlag{
					Name: "callback",
					Aliases: []string{"c"},
					Value: "",
					Usage: "Your callback host .e.g. http://myapp.com",
					Required: true,
				},
				&cli.StringFlag{
					Name: "key",
					Aliases: []string{"k"},
					Value: "",
					Usage: "Subscription key which provides access to this API. Found in your profile Primary Key",
					Required: true,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func createSandboxUserCmd(c *cli.Context) error {
	fmt.Printf("%s - %s\n", c.String("callback"), c.String("key"))
	return nil
}
