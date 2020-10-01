package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	app := &cli.App{
		Name:        "enve",
		Usage:       "run a program in a modified environment using .env files",
		Description: "Set all environment variables of one .env file and run a `command`.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   ".env",
				Usage:   "load environment variables from a file path",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "text",
				Usage:   "output environment variables in specific format (text, json)",
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

	// 3. Output flag
	output := ctx.String("output")

	if ctx.NArg() == 0 {
		switch output {
		case "json":
			return jsonOutputAction(ctx)
		default:
			return textOutputAction(ctx)
		}
	}

	// 4. Execute the given command
	if ctx.NArg() > 0 {
		return execCmdAction(ctx)
	}

	return nil
}

// textOutputAction prints all environment variables in plain text
func textOutputAction(ctx *cli.Context) (err error) {
	for _, s := range os.Environ() {
		fmt.Println(s)
	}

	return nil
}

// execCmdAction executes a command along with its env variables
func execCmdAction(ctx *cli.Context) (err error) {
	args := ctx.Args().Slice()
	cmdIn := args[0]

	_, err = exec.LookPath(cmdIn)

	if err != nil {
		return fmt.Errorf("executable \"%s\" was not found\n%s", cmdIn, err)
	}

	cmd := exec.Command(cmdIn, args[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func jsonOutputAction(ctx *cli.Context) (err error) {
	jsonstr := ""
	envs := os.Environ()

	for i, s := range envs {
		pairs := strings.SplitN(s, "=", 2)
		sep := ""

		if i < len(envs)-1 {
			sep = ","
		}

		val := strings.ReplaceAll(pairs[1], "\"", "\\\"")
		val = strings.ReplaceAll(val, "\n", "\\n")
		val = strings.ReplaceAll(val, "\r", "\\r")

		jsonstr += fmt.Sprintf("\"%s\":\"%s\"%s", pairs[0], val, sep)
	}

	var out bytes.Buffer
	json.HTMLEscape(&out, []byte("{"+jsonstr+"}"))
	out.WriteTo(os.Stdout)

	return nil
}
