package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	app := &cli.App{
		Name:        "enve",
		Usage:       "run a program in a modified environment using .env files",
		Description: "Set all environment variables of one .env file and run `command`.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   ".env",
				Usage:   "read in a file of environment variables",
			},
			VersionFlag(),
		},
		Action: onCommand,
	}

	err := app.Run(os.Args)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func onCommand(ctx *cli.Context) error {
	// 1. Version flag
	v := ctx.Bool("version")

	if v {
		return VersionAction(ctx)
	}

	// 2. File flag
	f := ctx.String("file")

	if f != "" {
		err := godotenv.Load(f)

		if err != nil {
			return err
		}
	}

	if ctx.NArg() > 0 {
		args := ctx.Args().Slice()
		cmdIn := args[0]

		_, err := exec.LookPath(cmdIn)

		if err != nil {
			return fmt.Errorf("executable \"%s\" was not found\n%s", cmdIn, err)
		}

		cmd := exec.Command(cmdIn, args[1:]...)

		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout

		return cmd.Run()
	}

	return nil
}
