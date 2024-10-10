package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC     GRPCConfig    `yaml:"grpc"`
	DB       DB            `yaml:"db"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DB struct {
	Host       string `yaml:"host" env-required:"true"`
	DBPort     string `yaml:"port" env-required:"true"`
	Username   string `yaml:"username" env-required:"true"`
	DBName     string `yaml:"dbname" env-required:"true"`
	DBPassword string `yaml:"dbpassword" env-required:"true" env:"DB_PASSWORD"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	path := fetchConfigPath()
	if path == "" {
		log.Fatalf("config path is empty: %s", err)
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
//
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
