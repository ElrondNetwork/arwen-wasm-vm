package main

import (
	"fmt"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwendebug"
	"github.com/urfave/cli"
)

func initializeCLI(facade *arwendebug.DebugFacade) *cli.App {
	app := cli.NewApp()
	app.Name = "Arwen Debug"
	app.Usage = ""

	args := &cliArguments{}

	flagServerAddress := cli.StringFlag{
		Name:        "address",
		Usage:       "",
		Value:       ":9091",
		Destination: &args.ServerAddress,
	}

	flagDatabase := cli.StringFlag{
		Name:        "database",
		Usage:       "",
		Value:       "./db",
		Destination: &args.Database,
	}

	flagSession := cli.StringFlag{
		Name:        "session",
		Usage:       "",
		Value:       "default",
		Destination: &args.Session,
	}

	app.Flags = []cli.Flag{
		flagSession,
		flagDatabase,
	}

	app.Authors = []cli.Author{
		{
			Name:  "The Elrond Team",
			Email: "contact@elrond.com",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "start debug server",
			Action: func(context *cli.Context) error {
				facade.StartServer(args.ServerAddress)
				return nil
			},
			Flags: []cli.Flag{
				flagServerAddress,
			},
		},
		{
			Name:  "deploy",
			Usage: "deploy a smart contract",
			Action: func(context *cli.Context) error {
				response, err := facade.DeploySmartContract(args.toDeployRequest())
				if err != nil {
					return err
				}

				fmt.Println(response.DebugString())
				return nil
			},
		},
		{
			Name:  "upgrade",
			Usage: "upgrade smart contract",
			Action: func(context *cli.Context) error {
				facade.UpgradeSmartContract(args.toUpgradeRequest())
				return nil
			},
		},
	}

	return app
}
