package env

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

type Slice []string

func (e Slice) Text() string {
	return strings.Join(e, "\n")
}

func (e Slice) Environ() (environ Environment, err error) {
	str := ""
	for i, s := range e {
		pairs := strings.SplitN(s, "=", 2)
		sep := ""
		if i < len(e)-1 {
			sep = ","
		}
		val := strings.ReplaceAll(pairs[1], "\"", "\\\"")
		val = strings.ReplaceAll(val, "\n", "\\n")
		val = strings.ReplaceAll(val, "\\", "\\\\")
		val = strings.ReplaceAll(val, "\r", "\\r")
		str += fmt.Sprintf("{\"name\":\"%s\",\"value\":\"%s\"}%s", pairs[0], val, sep)
	}
	jsonb := []byte("{\"environment\":[" + str + "]}")
	if err := json.Unmarshal(jsonb, &environ); err != nil {
		return environ, err
	}
	return environ, nil
}

func (e Slice) JSON() ([]byte, error) {
	jsonenv, err := e.Environ()
	if err != nil {
		return []byte(nil), err
	}
	jsonb, err := json.Marshal(jsonenv)
	if err != nil {
		return []byte(nil), err
	}
	return jsonb, nil
}

func (e Slice) XML() ([]byte, error) {
	jsonenv, err := e.Environ()
	if err != nil {
		return []byte(nil), err
	}
	xmlb, err := xml.Marshal(jsonenv)
	if err != nil {
		return []byte(nil), err
	}
	return xmlb, nil
}
