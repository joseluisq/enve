package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"

	cli "github.com/joseluisq/cline"
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
	app := cli.New()
	app.Name = "enve"
	app.Summary = "run a program in a modified environment using .env files"
	app.Version = versionNumber
	app.BuildTime = buildTime
	app.Flags = []cli.Flag{
		cli.FlagString{
			Name:    "file",
			Aliases: []string{"f"},
			Value:   ".env",
			Summary: "load environment variables from a file path (optional)",
		},
		cli.FlagString{
			Name:    "output",
			Aliases: []string{"o"},
			Value:   "text",
			Summary: "output environment variables using text, json or xml format",
		},
	}
	app.Handler = onCommand

	if err := app.Run(os.Args); err != nil {
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

func onCommand(ctx *cli.AppContext) error {
	// 2. File flag
	f := ctx.Flags.String("file")
	if f != "" {
		if exist := fileExists(f); exist {
			err := godotenv.Load(f)
			if err != nil {
				return err
			}
		}
	}

	tArgs := ctx.TailArgs

	// 4. Execute the given command
	if len(tArgs) > 0 {
		return execCmdAction(tArgs)
	}

	output := ctx.Flags.String("output")

	// 3. Output flag
	switch output {
	case "json":
		return jsonPrintAction()
	case "xml":
		return xmlPrintAction()
	case "text":
		return textPrintAction()
	}

	return nil
}

// execCmdAction executes a command along with its env variables
func execCmdAction(tArgs []string) (err error) {
	cmdIn := tArgs[0]
	_, err = exec.LookPath(cmdIn)
	if err != nil {
		return fmt.Errorf("executable \"%s\" was not found\n%s", cmdIn, err)
	}

	cmd := exec.Command(cmdIn, tArgs[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// textPrintAction prints all environment variables in plain text
func textPrintAction() (err error) {
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
func jsonPrintAction() error {
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
func xmlPrintAction() error {
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
