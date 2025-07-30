package env

import (
	"fmt"
)

type Map map[string]string

func (e Map) Array() []string {
	vars := []string{}
	for k, v := range e {
		vars = append(vars, fmt.Sprintf("%s=%s", k, v))
	}
	return vars
}
