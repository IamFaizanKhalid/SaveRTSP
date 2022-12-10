package main

import (
	"flag"
	"fmt"
	"github.com/IamFaizanKhalid/SaveRTSP/download"
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

	err = download.Start(cfg)
	if err != nil {
		log.Fatalf("SaveRTSP run err: %v\n", err)
	}
}

func readConfig() (download.Config, error) {
	execDir, err := os.Getwd()
	if err != nil {
		return download.Config{}, err
	}

	var configPath string

	flag.StringVar(&configPath, "config", "", "Config path (optional)")
	flag.Parse()

	var cfg download.Config
	if configPath == "" {
		configPath = "config.yaml"
	}

	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(execDir, configPath)
	}
	_, err = os.Stat(configPath)
	if err != nil {
		return download.Config{}, fmt.Errorf("can't load config file: %w", err)
	}

	err = configor.Load(&cfg, configPath)
	if err != nil {
		return download.Config{}, fmt.Errorf("can't load config file: %w", err)
	}

	return cfg, nil
}
