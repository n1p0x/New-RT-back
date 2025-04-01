package config

import (
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/env"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/knadh/koanf"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	Mode         string       `json:"mode"`
	Origin       string       `json:"origin"`
	ServerConfig ServerConfig `json:"server"`
	DBConfig     DBConfig     `json:"db"`
	TgConfig     TgConfig     `json:"tg"`
	TonConfig    TonConfig    `json:"ton"`
}

type ServerConfig struct {
	Host             string        `json:"host"`
	Port             int           `json:"port"`
	WriteTimeout     time.Duration `json:"writeTimeout"`
	GracefulShutdown time.Duration `json:"gracefulShutdown"`
}

type DBConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Port     int    `json:"port"`
}

type TgConfig struct {
	BotToken      string `json:"botToken"`
	AdminBotToken string `json:"adminBotToken"`
	ClientID      int    `json:"clientID"`
	ClientHash    string `json:"clientHash"`
	ClientPhone   string `json:"clientPhone"`
}

type TonConfig struct {
	IsTestnet              bool   `json:"isTestnet"`
	TonCenterApiKey        string `json:"tonCenterApiKey"`
	TonCenterApiKeyTestnet string `json:"tonCenterApiKeyTestnet"`
	AdminWallet            string `json:"adminWallet"`
	Mnemonic               string `json:"mnemonic"`
}

func Load(configPath string, envPath string) (*Config, error) {
	k := koanf.New(".")

	err := k.Load(confmap.Provider(defaultConfig, "."), nil)
	if err != nil {
		log.Printf("failed to load default config; err: %v", err)
		return nil, err
	}

	if configPath != "" {
		path, err := filepath.Abs(configPath)
		if err != nil {
			log.Printf("failed to get absolute config path; configPath: %s, err: %v", configPath, err)
			return nil, err
		}
		log.Printf("load config file from %s", path)
		if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
			log.Printf("failed to load config from file; err: %v", err)
			return nil, err
		}
	}

	if envPath != "" {
		path, err := filepath.Abs(envPath)
		if err != nil {
			log.Printf("failed to get absolute env path; envPath: %s, err: %v", envPath, err)
			return nil, err
		}

		if err := godotenv.Load(path); err != nil {
			log.Printf("failed to load env from file; err: %v", err)
			return nil, err
		}

		log.Printf("load env file from %s", path)
		if err := k.Load(env.Provider("ENV_", ".", func(s string) string {
			return strings.Replace(strings.ToLower(
				strings.TrimPrefix(s, "ENV_")), "_", ".", -1)
		}), nil); err != nil {
			log.Printf("failed to load env from file; err: %v", err)
			return nil, err
		}
	}

	var cfg Config
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "json", FlatPaths: false}); err != nil {
		log.Printf("failed to unmarshal with conf; err: %v", err)
		return nil, err
	}

	return &cfg, err
}
