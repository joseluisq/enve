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

// Build-time application values
var (
	versionNumber string = "devel"
	buildTime     string
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
	app.Summary = "Run a program in a modified environment using .env files"
	app.Version = versionNumber
	app.BuildTime = buildTime
	app.Flags = []cli.Flag{
		cli.FlagString{
			Name:    "file",
			Aliases: []string{"f"},
			Value:   ".env",
			Summary: "Load environment variables from a file path (optional)",
		},
		cli.FlagString{
			Name:    "output",
			Aliases: []string{"o"},
			Value:   "text",
			Summary: "Output environment variables using text, json or xml format",
		},
	}
	app.Handler = appHandler

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func appHandler(ctx *cli.AppContext) error {
	flags := ctx.Flags
	tailArgs := ctx.TailArgs

	// 1. Load a .env file if it's available
	var err error = nil
	file := flags.String("file")
	found := fileExists(file)
	fileProvided := flags.IsProvidedFlag("file")

	if fileProvided {
		if file == "" {
			return fmt.Errorf("File path was empty or not provided")
		}
		if !found {
			return fmt.Errorf("File path was not found or inaccessible")
		}
		err = godotenv.Load(file)
	} else if found {
		err = godotenv.Load(file)
	}
	if err != nil {
		return fmt.Errorf("Env file: %v", err)
	}

	// 2. Print all env variables in text format by default
	providedFlags := len(flags.GetProvidedFlags())
	if (providedFlags == 0 && len(tailArgs) == 0) ||
		(providedFlags == 1 && len(tailArgs) == 0 && fileProvided) {
		return printEnvText()
	}

	// 3. Output flag
	if flags.IsProvidedFlag("output") {
		out := flags.String("output")
		switch out {
		case "json":
			return printEnvJSON()
		case "xml":
			return printEnvXML()
		case "text":
			return printEnvText()
		default:
			if out == "" {
				return fmt.Errorf("Output format was empty or not provided")
			}
			return fmt.Errorf("Format `%s` is not supported output", out)
		}
	}

	// 4. Execute the given command if there is tail args passed
	if len(tailArgs) > 0 {
		return execProdivedCmd(tailArgs)
	}

	return nil
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

// execProdivedCmd executes a command along with its env variables
func execProdivedCmd(tArgs []string) (err error) {
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

// printEnvText prints all environment variables in plain text
func printEnvText() (err error) {
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
	if err := json.Unmarshal(jsonb, &jsonu); err != nil {
		return jsonu, err
	}
	return jsonu, nil
}

// printEnvJSON prints all environment variables in JSON format
func printEnvJSON() error {
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

// printEnvXML prints all environment variables in XML format
func printEnvXML() error {
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
