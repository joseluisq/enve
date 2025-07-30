package cmd

import (
	"fmt"
	"os"

	cli "github.com/joseluisq/cline"

	"github.com/joseluisq/enve/env"
	"github.com/joseluisq/enve/fs"
)

// Build-time application values
var (
	versionNumber string = "devel"
	buildTime     string
	buildCommit   string
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	app := cli.New()
	app.Name = "enve"
	app.Summary = "Run a program in a modified environment providing an optional .env file or variables from stdin"
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
			Summary: "Start a new environment with only variables from the .env file or stdin",
		},
		cli.FlagBool{
			Name:    "ignore-environment",
			Aliases: []string{"i"},
			Value:   false,
			Summary: "Starts with an empty environment, ignoring any existing environment variables",
		},
		cli.FlagBool{
			Name:    "no-file",
			Aliases: []string{"z"},
			Value:   false,
			Summary: "Do not load a .env file",
		},
		cli.FlagBool{
			Name:    "stdin",
			Aliases: []string{"s"},
			Value:   false,
			Summary: "Read only environment variables from stdin and ignore the .env file",
		},
	}
	app.Handler = appHandler

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
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

	// 1. Load a .env file if available
	file, err := flags.String("file")
	if err != nil {
		return err
	}
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

	var envVars env.Slice

	// stdin option
	stdinF, err := flags.Bool("stdin")
	if err != nil {
		return err
	}
	stdin, err := stdinF.Value()
	if err != nil {
		return err
	}

	overwriteF, err := flags.Bool("overwrite")
	if err != nil {
		return err
	}
	overwrite, err := overwriteF.Value()
	if err != nil {
		return err
	}

	if stdin {
		fi, err := os.Stdin.Stat()
		if err != nil {
			return fmt.Errorf("error: cannot read from stdin.\n%v", err)
		}
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			envr := env.FromReader(os.Stdin)

			if ignoreEnv {
				goto ContinueEnvProcessing
			}

			if newEnv {
				vmap, err := envr.Parse()
				if err != nil {
					return err
				}
				envVars = vmap.Array()
			} else {
				if err := envr.Load(overwrite); err != nil {
					str := ""
					if overwrite {
						str = " (overwrite)"
					}
					return fmt.Errorf("error: cannot load env from stdin%s.\n%v", str, err)
				}
				envVars = env.Slice(os.Environ())
			}

			goto ContinueEnvProcessing
		}
	}

	if !ignoreEnv {
		if noFile {
			envVars = env.Slice(os.Environ())
			goto ContinueEnvProcessing
		}

		// .env file processing
		envf, err := env.FromPath(filePath)
		if err != nil {
			return err
		}

		if newEnv {
			vmap, err := envf.Parse()
			if err != nil {
				return err
			}
			envVars = vmap.Array()
		} else {
			if err := envf.Load(overwrite); err != nil {
				str := ""
				if overwrite {
					str = " (overwrite)"
				}
				return fmt.Errorf("error: cannot load env from file '%s'.\n%v", str, err)
			}

			envVars = env.Slice(os.Environ())
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
		if err := fs.DirExists(chdirPath); err != nil {
			return err
		}
	}

	tailArgs := ctx.TailArgs

	// 2. Print all env variables in text format by default
	totalFags := len(flags.GetProvided())
	noFlags := totalFags == 0
	hasTailArgs := len(tailArgs) > 0
	hasNoArgs := noFlags && hasTailArgs

	if hasNoArgs {
		fmt.Println(envVars.Text())
		return nil
	}

	// 3. Execute the given command if there is tail args passed
	if hasTailArgs {
		return execProdivedCmd(tailArgs, chdirPath, newEnv, envVars)
	}

	// 4. Output
	output, err := flags.String("output")
	if err != nil {
		return err
	}

	out := output.Value()
	switch out {
	case "text":
		fmt.Println(envVars.Text())
	case "json":
		if buf, err := envVars.JSON(); err != nil {
			return err
		} else {
			fmt.Println(string(buf))
		}
	case "xml":
		if buf, err := envVars.XML(); err != nil {
			return err
		} else {
			fmt.Println("<?xml version=\"1.0\" encoding=\"UTF-8\"?>" + string(buf))
		}
	default:
		if out == "" {
			return fmt.Errorf("error: output format was empty or not provided")
		}
		return fmt.Errorf("error: output format '%s' is not supported", out)
	}

	return nil
}
