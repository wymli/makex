package config

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/wymli/makex/shell"
	"gopkg.in/yaml.v2"
)

const (
	CONFIG_NAME   = "makex_config"
	CONFIG_TYPE   = "yaml"
	MAKEXFILE_KEY = "makexfile"
	MAKEXFILE     = "makex.yaml"
)

var (
	HOME, _     = os.UserHomeDir()
	CONFIG_DIR  = filepath.Join(HOME, ".makex")
	CONFIG_PATH = filepath.Join(CONFIG_DIR, CONFIG_NAME+"."+CONFIG_TYPE)

	SHELL_DIR = filepath.Join(CONFIG_DIR, "shell")
)

//go:embed config.yaml
var CofigTpl []byte

type Config struct {
	Makexfile string `yaml:"makexfile,omitempty"`
	Template  string `yaml:"template,omitempty"`
}

func InitMakexConfig() error {
	// 1. check exists
	stat, err := os.Stat(CONFIG_DIR)
	if os.IsNotExist(err) || !stat.IsDir() {
		log.Debugf("[init] %s not found, create default", CONFIG_DIR)
		// 2. mkdir
		if err := os.MkdirAll(CONFIG_DIR, os.ModePerm); err != nil {
			return err
		}

		// 3. create config
		if err := os.WriteFile(CONFIG_PATH, CofigTpl, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write config file to %s, err: %w", CONFIG_PATH, err)
		}

		// 4. Move shell
		return MoveShells()
	}

	log.Debugf("[init] %s found, skip", CONFIG_DIR)
	return nil
}

func ReadMakexConfig() (*Config, error) {
	// 1. init config file
	_, err := os.Stat(CONFIG_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to stat config file in %s, err: %v", CONFIG_PATH, err)
	}

	// 2. load config file
	data, err := ioutil.ReadFile(CONFIG_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file in %s, err: %v", CONFIG_PATH, err)
	}

	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to marhsall config file in %s to go struct, err: %v", CONFIG_PATH, err)
	}

	return &config, nil
}

func WriteMakexfile(c *Config, makexfile string) error {
	log.Debugf("[template init] write makexfile: %s", makexfile)

	err := os.WriteFile(makexfile, []byte(c.Template), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create makexfile %s, err: %v", makexfile, err)
	}

	return nil
}

func MoveShells() error {
	dirEntries, err := shell.ShellFS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read embedfs dir, err: %v", err)
	}

	if err := os.MkdirAll(SHELL_DIR, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create dir '%s', err: %v", SHELL_DIR, err)
	}

	for _, entry := range dirEntries {
		if entry.Type().IsDir() {
			continue
		}

		data, err := shell.ShellFS.ReadFile(entry.Name())
		if err != nil {
			return fmt.Errorf("failed to open embedfs file, err: %v", err)
		}

		err = os.WriteFile(filepath.Join(SHELL_DIR, entry.Name()), data, os.ModePerm)
		if err != nil {
			return fmt.Errorf("faild to write file to SHELLDIR '%s', err: %v", SHELL_DIR, err)
		}
	}
	return nil
}
