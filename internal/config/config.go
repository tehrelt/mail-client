package config

import "github.com/ilyakaznacheev/cleanenv"

type SmtpConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Pop3Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type AppConfig struct {
	Smtp SmtpConfig `yaml:"smtp"`
	Pop3 Pop3Config `yaml:"pop3"`
}

func LoadConfig(path string) *AppConfig {
	var config AppConfig
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic(err)
	}
	return &config
}
