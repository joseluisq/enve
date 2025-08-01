package cmd

import (
	"fmt"
	"os"

	cli "github.com/joseluisq/cline"

	"github.com/joseluisq/enve/env"
	"github.com/joseluisq/enve/fs"
)

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

	// overwrite option
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
				goto ContinueEnvProc
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

			goto ContinueEnvProc
		}
	}

	if !ignoreEnv {
		if noFile {
			envVars = env.Slice(os.Environ())
			goto ContinueEnvProc
		}

		// .env file processing
		envf, err := env.FromPath(filePath)
		if err != nil {
			return err
		}
		defer envf.Close()

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

ContinueEnvProc:
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
	hasNoArgs := noFlags && !hasTailArgs

	if hasNoArgs {
		fmt.Println(envVars.Text())
		return nil
	}

	// 3. Output
	output, err := flags.String("output")
	if err != nil {
		return err
	}

	if output.IsProvided() {
		if hasTailArgs {
			return fmt.Errorf("error: output format cannot be used when executing a command")
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
	}

	// 4. Execute the given command if there is tail args passed
	if hasTailArgs {
		return execCmd(tailArgs, chdirPath, newEnv, envVars)
	}

	return nil
}
