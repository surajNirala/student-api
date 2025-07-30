package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required`
}

// env-defult : "production"
type Config struct {
	Env         string       `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string       `yaml:"storage_path"`
	MySQL       *MySQLConfig `yaml:"mysql"`
	HTTPServer  `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "Please Pass the config file")
		flag.Parse()
		configPath = *flags
		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Can not read config file %s", err.Error())
	}
	return &cfg
}

type MySQLConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
}

// type MysqlConfig struct {
// 	Env        string      `yaml:"env" env:"ENV" env-required:"true"`
// 	MySQL      MySQLConfig `yaml:"mysql"`
// 	HTTPServer HTTPServer  `yaml:"http_server"`
// }

func MustLoadMySQL() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "Please Pass the config file")
		flag.Parse()
		configPath = *flags
		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Cannot read config file: %s", err.Error())
	}

	return &cfg
}
