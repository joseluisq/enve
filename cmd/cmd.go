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
		cli.FlagString{
			Name:    "chdir",
			Aliases: []string{"c"},
			Summary: "Change currrent working directory",
		},
		cli.FlagBool{
			Name:    "new-environment",
			Aliases: []string{"n"},
			Value:   false,
			Summary: "Start a new environment containing only variables from the .env file",
		},
		cli.FlagBool{
			Name:    "ignore-environment",
			Aliases: []string{"i"},
			Value:   false,
			Summary: "Start with an empty environment",
		},
		cli.FlagBool{
			Name:    "no-file",
			Aliases: []string{"z"},
			Value:   false,
			Summary: "Do not load a .env file",
		},
	}
	app.Handler = appHandler

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func appHandler(ctx *cli.AppContext) error {
	var flags = ctx.Flags

	// ignore-environment option
	ignoreEnvF, err := flags.Bool("ignore-environment")
	if err != nil {
		return err
	}
	ignoreEnv, err := ignoreEnvF.Value()
	if err != nil {
		return err
	}

	// no-file option
	noFileF, err := flags.Bool("no-file")
	if err != nil {
		return err
	}
	noFile, err := noFileF.Value()
	if err != nil {
		return err
	}

	// 1. Load a .env file if it's available
	file, err := flags.String("file")
	if err != nil {
		return err
	}
	fileProvided := file.IsProvided()
	filePath := file.Value()

	// new-environment option
	newEnvF, err := flags.Bool("new-environment")
	if err != nil {
		return err
	}
	newEnv, err := newEnvF.Value()
	if err != nil {
		return err
	}

	var envVars []string

	if !ignoreEnv {
		if noFile {
			envVars = os.Environ()
			goto ContinueEnvProcessing
		}

		// .env file processing

		if err := fileExists(filePath); err != nil {
			return err
		}

		if newEnv {
			if envVars, err = parseEnvFile(filePath); err != nil {
				return err
			}
		} else {
			// Overwrite option
			overwrite, err := flags.Bool("overwrite")
			if err != nil {
				return err
			}
			if overwriteValue, err := overwrite.Value(); err != nil {
				return err
			} else if overwriteValue {
				if err := godotenv.Overload(filePath); err != nil {
					return fmt.Errorf("cannot load env file (overwrite): %v", err)
				}
			} else {
				if err := godotenv.Load(filePath); err != nil {
					return fmt.Errorf("cannot load env file: %v", err)
				}
			}

			envVars = os.Environ()
		}
	}

ContinueEnvProcessing:

	// chdir option
	chdirPath := ""
	chdir, err := flags.String("chdir")
	if err != nil {
		return err
	}
	if chdir.IsProvided() {
		chdirPath = chdir.Value()
		if err := dirExists(chdirPath); err != nil {
			return err
		}
	}

	tailArgs := ctx.TailArgs

	// 2. Print all env variables in text format by default
	providedFlags := len(flags.GetProvided())
	if (providedFlags == 0 && len(tailArgs) == 0) ||
		(providedFlags <= 2 && len(tailArgs) == 0 && fileProvided) {
		return printEnvText(envVars)
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
			return printEnvJSON(envVars)
		case "xml":
			return printEnvXML(envVars)
		case "text":
			return printEnvText(envVars)
		default:
			if out == "" {
				return fmt.Errorf("output format was empty or not provided")
			}
			return fmt.Errorf("format `%s` is not a supported output", out)
		}
	}

	// 4. Execute the given command if there is tail args passed
	if len(tailArgs) > 0 {
		return execProdivedCmd(tailArgs, chdirPath, newEnv, envVars)
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

func dirExists(dirname string) error {
	if dirname == "" {
		return fmt.Errorf("directory path was empty or not provided")
	}
	if info, err := os.Stat(dirname); err != nil {
		return fmt.Errorf("cannot access directory: %s; %v", dirname, err)
	} else {
		if !info.IsDir() {
			return fmt.Errorf("directory path is a file: %s", dirname)
		}
		return nil
	}
}

func parseEnvFile(filePath string) ([]string, error) {
	vars := []string{}
	f, err := os.Open(filePath)
	if err != nil {
		return vars, err
	}
	defer f.Close()

	if envMap, err := godotenv.Parse(f); err != nil {
		return vars, fmt.Errorf("cannot read env file: %v", err)
	} else {
		for k, v := range envMap {
			vars = append(vars, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return vars, nil
}

// printEnvText prints all environment variables in plain text.
func printEnvText(envVars []string) (err error) {
	for _, s := range envVars {
		fmt.Println(s)
	}
	return nil
}

// parseJSONFromEnviron decodes (Unmarshal) a list of environment variables into a JSON struct.
func parseJSONFromEnviron(envs []string) (jsonu Environment, err error) {
	jsonstr := ""
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

// printEnvJSON prints all environment variables in JSON format.
func printEnvJSON(envVars []string) error {
	jsonu, err := parseJSONFromEnviron(envVars)
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

// printEnvXML prints all environment variables in XML format.
func printEnvXML(envVars []string) error {
	jsonu, err := parseJSONFromEnviron(envVars)
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
