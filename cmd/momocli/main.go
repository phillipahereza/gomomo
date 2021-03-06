package main

import (
	"fmt"
	"github.com/phillipahereza/gomomo"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()

	app.Commands = []*cli.Command{
		{
			Name:    "Momo Sandbox",
			Aliases: []string{"sandbox"},
			Usage:   "Create a sandbox environment API user ",
			Action:  createSandboxUserCmd,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "callback",
					Aliases:  []string{"c"},
					Value:    "",
					Usage:    "Your callback host .e.g. http://myapp.com",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "key",
					Aliases:  []string{"k"},
					Value:    "",
					Usage:    "Subscription key which provides access to this API. Found in your profile Primary Key",
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
	client := gomomo.NewClient(c.String("key"), "sandbox", "https://sandbox.momodeveloper.mtn.com/")
	refID, err := client.Sandbox.CreateSandboxUser(c.String("callback"))
	if err != nil {
		log.Fatal(err)
	}
	apiKey, err := client.Sandbox.GenerateSandboxUserAPIKey(refID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("API Key: %s\n", apiKey.APIKey)
	fmt.Printf("User ID: %s\n", refID)
	return nil
}
