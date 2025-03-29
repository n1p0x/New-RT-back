package config

import (
	"log"
	"path/filepath"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	Mode         string       `json:"mode"`
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
	Migrate  struct {
		Enable bool   `json:"enable"`
		Dir    string `json:"dir"`
	} `json:"migrate"`
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

func Load(configPath string) (*Config, error) {
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

	var cfg Config
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "json", FlatPaths: false}); err != nil {
		log.Printf("failed to unmarshal with conf; err: %v", err)
		return nil, err
	}
	// koanf.UnmarshalConf преобразует данные, собранные из всех источников в структуру Config
	// UnmarshalWithConf(path string, o interface{}, conf UnmarshalConf)
	// path string - путь в конфигурации, с которого начинается размаршалинг.
	// Если пустой, используется вся конфикурация от корна
	// o interface{} - указатель на объект (обычно структура), в который будут записаны данные
	// conf UnmarshalConf - структура с настройками процесса размаршалинга
	// type UnmarshalConf struct {
	//	 Tag         string            // Тег для маппинга (например, "json")
	//	 FlatPaths   bool              // Обрабатывать ли ключи как плоские или вложенные
	//	 Decoder     Decoder           // Пользовательский декодер (опционально)
	//	 TagFallback []string          // Дополнительные теги (если основной не найден)
	// }

	return &cfg, err
}
