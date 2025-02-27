package config

import (
	"fmt"
	"log"
	"os"

	"github.com/appwrite/sdk-for-go/models"
	"gopkg.in/yaml.v2"
)

// Config represents the application's configuration
type Config struct {
	Appwrite AppwriteConfig `yaml:"appwrite"`
	Network  NetworkConfig  `yaml:"network"`
	External ExternalConfig `yaml:"external"`
	Mail     MailConfig     `yaml:"mail"`
}

// DatabaseConfig holds the database-related configuration
type AppwriteConfig struct {
	Endpoint string `yaml:"endpoint"`
	Project  string `yaml:"project"`
	Key      string `yaml:"key"`

	// Appwrite SDK
	Database  *models.Database
	ColHosts  *models.Collection
	ColStatus *models.Collection
	ColAlert  *models.Collection
	ColCheck  *models.Collection
	ColRhp2   *models.Collection
	ColRhp3   *models.Collection
	ColRhp4   *models.Collection
}

// ServerConfig holds the server-related configuration
type NetworkConfig struct {
	Name string `yaml:"name"`
}

type ExternalConfig struct {
	BenchUrl    string `yaml:"benchUrl"`
	ExploredUrl string `yaml:"exploredUrl"`
}

type MailConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var appConfig *Config

// LoadConfig loads the configuration from the given file path
func LoadConfig(filepath string) *Config {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	appConfig = &config
	fmt.Println("Config loaded:", filepath)
	fmt.Println("Appwrite Endpoint:", appConfig.Appwrite.Endpoint)
	return appConfig
}

// GetConfig returns the loaded configuration
func GetConfig() *Config {
	if appConfig == nil {
		log.Fatalf("Config not loaded")
	}
	return appConfig
}
