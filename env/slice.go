package env

import (
	"encoding/json"
	"encoding/xml"
	"strings"
)

type Slice []string

func (e Slice) Text() string {
	return strings.Join(e, "\n")
}

func (e Slice) Environ() Environment {
	var environ Environment
	for _, s := range e {
		// NOTE: skip non-key=value pair
		pair := strings.SplitN(s, "=", 2)
		if len(pair) < 2 {
			continue
		}
		v := EnvironmentVar{Name: pair[0], Value: pair[1]}
		environ.Env = append(environ.Env, v)
	}
	return environ
}

func (e Slice) JSON() ([]byte, error) {
	environ := e.Environ()
	jsonb, err := json.Marshal(environ)
	if err != nil {
		return []byte(nil), err
	}
	return jsonb, nil
}

func (e Slice) XML() ([]byte, error) {
	environ := e.Environ()
	xmlb, err := xml.Marshal(environ)
	if err != nil {
		return []byte(nil), err
	}
	return xmlb, nil
}
