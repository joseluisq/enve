package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

// Environment defines JSON/XML data structure
type Environment struct {
	Env []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"environment"`
}

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
				Usage:   "load environment variables from a file path (optional)",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "text",
				Usage:   "output environment variables using text, json or xml format",
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

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return !info.IsDir()
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
		if exist := fileExists(f); exist {
			err := godotenv.Load(f)

			if err != nil {
				return err
			}
		}
	}

	// 3. Output flag
	output := ctx.String("output")

	if ctx.NArg() == 0 {
		switch output {
		case "json":
			return jsonPrintAction(ctx)
		case "xml":
			return xmlPrintAction(ctx)
		default:
			return textPrintAction(ctx)
		}
	}

	// 4. Execute the given command
	if ctx.NArg() > 0 {
		return execCmdAction(ctx)
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

// textPrintAction prints all environment variables in plain text
func textPrintAction(ctx *cli.Context) (err error) {
	for _, s := range os.Environ() {
		fmt.Println(s)
	}

	return nil
}

// parseJSONFromEnviron decodes (Unmarshal) system environment variables into a JSON struct
func parseJSONFromEnviron() (jsonu Environment, err error) {
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
		jsonstr += fmt.Sprintf("{\"name\":\"%s\",\"value\":\"%s\"}%s", pairs[0], val, sep)
	}

	jsonb := []byte("{\"environment\":[" + jsonstr + "]}")
	err = json.Unmarshal(jsonb, &jsonu)

	if err != nil {
		return jsonu, err
	}

	return jsonu, nil
}

// jsonPrintAction prints all environment variables in JSON format
func jsonPrintAction(ctx *cli.Context) error {
	jsonu, err := parseJSONFromEnviron()
	if err != nil {
		return err
	}

	jsonb, err := json.Marshal(jsonu)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonb))

	return nil
}

// xmlPrintAction prints all environment variables in XML format
func xmlPrintAction(ctx *cli.Context) error {
	jsonu, err := parseJSONFromEnviron()
	if err != nil {
		return err
	}

	xmlb, err := xml.Marshal(jsonu)
	if err != nil {
		return err
	}

	fmt.Println("<?xml version=\"1.0\" encoding=\"UTF-8\"?>" + string(xmlb))

	return nil
}
