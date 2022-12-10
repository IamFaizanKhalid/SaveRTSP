package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

func main() {
	_, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}
}

func readConfig() (Config, error) {
	execDir, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}

	var configPath string

	flag.StringVar(&configPath, "config", "", "Config path (optional)")
	flag.Parse()

	var cfg Config
	if configPath == "" {
		configPath = "config.yaml"
	}

	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(execDir, configPath)
	}
	_, err = os.Stat(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("can't load config file: %w", err)
	}

	err = configor.Load(&cfg, configPath)
	if err != nil {
		return Config{}, fmt.Errorf("can't load config file: %w", err)
	}

	return cfg, nil
}

type Config struct {
	Cameras []Camera `yaml:"cameras"`
}

type Camera struct {
	Name      string `yaml:"name"`
	StreamUrl string `yaml:"stream_url"`
	Split     int    `yaml:"split"`
}
