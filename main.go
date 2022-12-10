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
	cfg, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	err = Run(cfg)
	if err != nil {
		log.Fatalf("SaveRTSP run err: %v\n", err)
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
