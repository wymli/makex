package config

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/wymli/makex/code"
	"gopkg.in/yaml.v2"
)

const (
	CONFIG_NAME = "makex_config"
	CONFIG_TYPE = "yaml"

	CONFIG_DIR_NAME = ".makex"
	SHELL_DIR_NAME  = "shell"
)

var (
	USER_HOME_PATH, _ = os.UserHomeDir()

	CONFIG_DIR_PATH = filepath.Join(USER_HOME_PATH, CONFIG_DIR_NAME)
	SHELL_DIR_PATH  = filepath.Join(CONFIG_DIR_PATH, SHELL_DIR_NAME)

	CONFIG_PATH = filepath.Join(CONFIG_DIR_PATH, CONFIG_NAME+"."+CONFIG_TYPE)
)

//go:embed config.yaml
var CofigTpl []byte

type Config struct {
	Makexfile string `yaml:"makexfile,omitempty"`
	Template  string `yaml:"template,omitempty"`
}

func InitMakexConfig() error {
	// 1. check exists
	stat, err := os.Stat(CONFIG_DIR_PATH)

	if os.IsExist(err) && !stat.IsDir() {
		log.Infof("[init] detected existed config dir '%s' is not a dir, remove it and have a try again", CONFIG_DIR_PATH)
		return fmt.Errorf("config dir is not a dir")
	}

	if os.IsNotExist(err) {
		log.Debugf("[init] %s not found, create default", CONFIG_DIR_PATH)

		// 2. mkdir
		if err := os.MkdirAll(CONFIG_DIR_PATH, os.ModePerm); err != nil {
			return err
		}

		// 3. create config
		if err := os.WriteFile(CONFIG_PATH, CofigTpl, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write config file to %s, err: %w", CONFIG_PATH, err)
		}

		// 4. Move shell
		return MoveShells()
	}

	log.Debugf("[init] %s found, skip", CONFIG_DIR_PATH)
	return nil
}

func ReadMakexConfig() (*Config, error) {
	// 1. load config file
	data, err := ioutil.ReadFile(CONFIG_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file in %s, err: %v", CONFIG_PATH, err)
	}

	// 2. unmarshall
	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarhsall config file in %s to go struct, err: %v", CONFIG_PATH, err)
	}

	return &config, nil
}

func WriteDefaultMakexfile(c *Config, makexfile string) error {
	log.Debugf("[template init] write makexfile: %s", makexfile)

	err := os.WriteFile(makexfile, []byte(c.Template), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create makexfile %s, err: %v", makexfile, err)
	}

	return nil
}

func MoveShells() error {
	dirEntries, err := code.CodeFS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read embedfs dir, err: %v", err)
	}

	if err := os.MkdirAll(SHELL_DIR_PATH, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create dir '%s', err: %v", SHELL_DIR_PATH, err)
	}

	for _, entry := range dirEntries {
		if entry.Type().IsDir() {
			continue
		}

		data, err := code.CodeFS.ReadFile(entry.Name())
		if err != nil {
			return fmt.Errorf("failed to open embedfs file, err: %v", err)
		}

		err = os.WriteFile(filepath.Join(SHELL_DIR_PATH, entry.Name()), data, os.ModePerm)
		if err != nil {
			return fmt.Errorf("faild to write file to SHELLDIR '%s', err: %v", SHELL_DIR_PATH, err)
		}
	}
	return nil
}
