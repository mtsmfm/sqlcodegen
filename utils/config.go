package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Out struct {
		Package string
		File    string
	}
	Tags    []string
	Typemap map[string]string
}

var ConfigFileName = "sqlcodegen.yml"

func FindConfigPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	for dir != filepath.Dir(dir) {
		path := filepath.Join(dir, ConfigFileName)
		_, err := os.Stat(path)

		if err != nil {
			dir = filepath.Dir(dir)
		} else {
			return path, nil
		}
	}

	return "", fmt.Errorf(ConfigFileName + " is not found. Run `sqlcodegen init`")
}

func LoadConfig(path string) (*Config, error) {
	var cfg *Config
	file, err := ioutil.ReadFile(path)

	err = yaml.Unmarshal(file, &cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
