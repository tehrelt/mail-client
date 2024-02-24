package main

import (
	"flag"
	"mail-client/internal/api"
	"mail-client/internal/config"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "config/config.yaml", "PATH TO YAML-CONFIG")
}

func main() {
	flag.Parse()

	cfg := config.LoadConfig(configPath)

	if err := api.Start(cfg); err != nil {
		panic(err)
	}
}
