// Copyright © 2020 Hedzr Yeh.

package vxconf

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// JSONToString convert json object `in` to json text string
func JSONToString(in interface{}, pretty bool) string {
	var b []byte
	var err error
	if pretty {
		b, err = json.MarshalIndent(in, "", "  ")
	} else {
		b, err = json.Marshal(in)
	}
	if err != nil {
		logrus.Errorf("json decode string failed: %v", err)
		return "<error-json-object>"
	}
	return string(b)
}

// YamlToString convert yaml object `in` to yaml text string
func YamlToString(in interface{}, pretty bool) string {
	var b []byte
	var err error
	if pretty {
		b, err = yaml.Marshal(in)
	} else {
		b, err = yaml.Marshal(in)
	}
	if err != nil {
		logrus.Errorf("yaml decode string failed: %v", err)
		return "<error-yaml-object>"
	}
	return string(b)
}
