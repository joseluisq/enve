package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	cli "github.com/joseluisq/cline"
)

// Build-time application values
var (
	versionNumber string = "devel"
	buildTime     string
	buildCommit   string
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
	app.BuildCommit = buildCommit
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
		cli.FlagBool{
			Name:    "overwrite",
			Aliases: []string{"w"},
			Value:   false,
			Summary: "Overwrite environment variables if already set",
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

	// 1. Load a .env file if it's available
	fileFound := false
	file, err := flags.String("file")
	if err != nil {
		return err
	}
	filePath := file.Value()
	fileProvided := file.IsProvided()
	if fileProvided {
		if err := fileExists(filePath); err != nil {
			return err
		}
		fileFound = true
	}

	overwrite, err := flags.Bool("overwrite")
	if err != nil {
		return err
	}
	overwriteValue, err := overwrite.Value()
	if err != nil {
		return err
	}

	if fileFound {
		if overwriteValue {
			if err := godotenv.Overload(filePath); err != nil {
				return fmt.Errorf("cannot load env file (overwrite): %v", err)
			}
		} else {
			if err := godotenv.Load(filePath); err != nil {
				return fmt.Errorf("cannot load env file: %v", err)
			}
		}
	}

	tailArgs := ctx.TailArgs

	// 2. Print all env variables in text format by default
	providedFlags := len(flags.GetProvided())
	if (providedFlags == 0 && len(tailArgs) == 0) ||
		(providedFlags <= 2 && len(tailArgs) == 0 && fileProvided) {
		return printEnvText()
	}

	// 3. Output
	output, err := flags.String("output")
	if err != nil {
		return err
	}
	if output.IsProvided() {
		out := output.Value()
		switch out {
		case "json":
			return printEnvJSON()
		case "xml":
			return printEnvXML()
		case "text":
			return printEnvText()
		default:
			if out == "" {
				return fmt.Errorf("output format was empty or not provided")
			}
			return fmt.Errorf("format `%s` is not a supported output", out)
		}
	}

	// 4. Execute the given command if there is tail args passed
	if len(tailArgs) > 0 {
		return execProdivedCmd(tailArgs)
	}

	return nil
}

func fileExists(filename string) error {
	if filename == "" {
		return fmt.Errorf("file path was empty or not provided")
	}
	if info, err := os.Stat(filename); err != nil {
		return fmt.Errorf("cannot access file: %s; %v", filename, err)
	} else {
		if info.IsDir() {
			return fmt.Errorf("file path is a directory: %s", filename)
		}
		return nil
	}
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
		val = strings.ReplaceAll(val, "\\", "\\\\")
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
