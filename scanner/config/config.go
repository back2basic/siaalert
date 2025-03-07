package config

import (
	"fmt"
	"log"
	"os"

	"github.com/back2basic/siaalert/scanner/db"
	"github.com/back2basic/siaalert/scanner/logger"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// Config represents the application's configuration
type Config struct {
	DB       *db.MongoDB
	MongoDB  MongoDBConfig  `yaml:"mongodb"`
	Network  NetworkConfig  `yaml:"network"`
	External ExternalConfig `yaml:"external"`
	Mail     MailConfig     `yaml:"mail"`
}

type MongoDBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// ServerConfig holds the server-related configuration
type NetworkConfig struct {
	Name string `yaml:"name"`
}

type ExternalConfig struct {
	// BenchUrl    string `yaml:"benchUrl"`
	ExploredUrl string `yaml:"exploredUrl"`
}

type MailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var appConfig *Config

// LoadConfig loads the configuration from the given file path
func LoadConfig(filepath string) *Config {
	log := logger.GetLogger()
	defer logger.Sync()
	file, err := os.Open(filepath)
	if err != nil {
		log.Error("Failed to open config file:", zap.Error(err))
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Error("Failed to decode config file:", zap.Error(err))
	}

	// Initialize MongoDB connection
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", config.MongoDB.Username, config.MongoDB.Password, config.MongoDB.Host, config.MongoDB.Port)
	// mongoDB, err := db.NewMongoDB(uri, config.MongoDB.Database, "collectionHost", "collectionScan", "collectionApi", "collectionAlert", "collectionRhp")
	mongoDB, err := db.NewMongoDB(
		uri,
		config.MongoDB.Database,
		config.Network.Name+"_host",
		config.Network.Name+"_scan",
		config.Network.Name+"_api",
		config.Network.Name+"_alert",
		config.Network.Name+"_rhp",
	)
	if err != nil {
		log.Error("Failed to initialize MongoDB connection:", zap.Error(err))
	} else {
		config.DB = mongoDB
		log.Info("MongoDB connection initialized")
	}

	appConfig = &config
	log.Info("Config loaded:", zap.String("filepath", filepath))
	return appConfig
}

// GetConfig returns the loaded configuration
func GetConfig() *Config {
	if appConfig == nil {
		log.Fatalf("Config not loaded")
	}
	return appConfig
}

func (c *Config) Close(log *zap.Logger) {
	if c.DB != nil {
		err := c.DB.Close()
		if err != nil {
			log.Error("Failed to close MongoDB connection", zap.Error(err))
		} else {
			log.Info("MongoDB connection closed")
		}
	}
}
