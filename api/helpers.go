package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

// respond reads the 'f' url parameter ('f' stands for 'format'), formats the given data
// accordingly and sets the required content-type header. Default format is json.
func respond(res http.ResponseWriter, req *http.Request, code int, data interface{}) {
	var err error
	var errMesg []byte
	var out []byte

	f := "json"
	format := req.URL.Query()["f"]
	if len(format) > 0 {
		f = format[0]
	}

	if f == "yaml" {
		res.Header().Set("Content-Type", "text/yaml; charset=utf-8")
		out, err = yaml.Marshal(data)
		errMesg = []byte("--- error: failed while rendering data to yaml")
	} else {
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		out, err = json.Marshal(data)
		errMesg = []byte("{ 'error': 'failed while rendering data to json' }")
	}

	if err != nil {
		out = errMesg
		res.WriteHeader(http.StatusInternalServerError)
	} else {
		res.Write(out)
	}
	return
}

func getBodyAsBytes(body io.ReadCloser) ([]byte, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

// parseBody reads the payload of a request and formats unmarshals it
// accoding to the format specified in the 'd' url parameter ('d' stands
// for 'data'). Default format is json.
func parseBody(req *http.Request, s interface{}) error {
	d := "json"
	format := req.URL.Query()["d"]
	if len(format) > 0 {
		d = format[0]
	}

	b, err := getBodyAsBytes(req.Body)
	if err != nil {
		return err
	}

	if string(b) == "[]" || string(b) == "{}" || string(b) == "" {
		return nil
	}

	if d == "yaml" {
		err = unmarshalYAML(b, s)
	} else {
		err = json.Unmarshal(b, s)
	}

	return err
}

// unmarshalYAML helps out with the YAML unmarshalling in order to provide
// compatibility to JSON marshalling. The default behavior fo go-yaml is to
// unmarshal nested maps to map[interface{}]interface{} values, and such
// values can't be marshalled as JSON.
//
// The keys of map[interface{}]interface{} maps will be converted to
// strings with a %v format string, as will any scalar values that
// aren't already strings (i.e. numbers and boolean values).
func unmarshalYAML(in []byte, s interface{}) error {
	var dest *map[string]interface{}
	var ok bool
	if dest, ok = s.(*map[string]interface{}); !ok {
		return errors.New("Expecting *map[string]interface{}")
	}

	var result map[interface{}]interface{}
	if err := yaml.Unmarshal(in, &result); err != nil {
		return err
	}
	m := cleanUpInterfaceMap(result)

	*dest = m
	return nil
}

func cleanUpInterfaceArray(in []interface{}) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = cleanUpMapValue(v)
	}
	return out
}

func cleanUpInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range in {
		out[fmt.Sprintf("%v", k)] = cleanUpMapValue(v)
	}
	return out
}

func cleanUpMapValue(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return cleanUpInterfaceArray(v)
	case map[interface{}]interface{}:
		return cleanUpInterfaceMap(v)
	default:
		return v
	}
}
