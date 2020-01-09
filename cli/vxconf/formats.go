// Copyright © 2019 Hedzr Yeh.

package vxconf

import (
	"encoding/json"
	"github.com/hedzr/njuone/cli/vxconf/cfg"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

// JSON -----------------------------------------------------------------------

// ParseJSON reads a JSON configuration from the given string.
func ParseJSON(cfg string) (*cfg.AppConfig, error) {
	return parseJSON([]byte(cfg))
}

// ParseJSONFile reads a JSON configuration from the given filename.
func ParseJSONFile(filename string) (*cfg.AppConfig, error) {
	cfg, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parseJSON(cfg)
}

// parseJSON performs the real JSON parsing.
func parseJSON(cfg []byte) (*cfg.AppConfig, error) {
	var out interface{}
	var err error
	if err = json.Unmarshal(cfg, &out); err != nil {
		return nil, err
	}
	if out, err = cfg.normalizeValue(out); err != nil {
		return nil, err
	}
	return &cfg.AppConfig{Root: out}, nil
}

// RenderJSON renders a JSON configuration.
func RenderJSON(cfg interface{}) (string, error) {
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// YAML -----------------------------------------------------------------------

// ParseYaml reads a YAML configuration from the given string.
func ParseYaml(cfg string) (*cfg.AppConfig, error) {
	return parseYaml([]byte(cfg))
}

// ParseYamlFile reads a YAML configuration from the given filename.
func ParseYamlFile(filename string) (*cfg.AppConfig, error) {
	cfg, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parseYaml(cfg)
}

// parseYaml performs the real YAML parsing.
func parseYaml(cfg []byte) (*cfg.AppConfig, error) {
	var out interface{}
	var err error
	if err = yaml.Unmarshal(cfg, &out); err != nil {
		return nil, err
	}
	if out, err = cfg.normalizeValue(out); err != nil {
		return nil, err
	}
	return &cfg.AppConfig{Root: out}, nil
}

// RenderYaml renders a YAML configuration.
func RenderYaml(cfg interface{}) (string, error) {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return UnescapeUnicode(b), nil
}
