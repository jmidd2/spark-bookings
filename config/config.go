package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"spark-bookings/client"
	"spark-bookings/utils"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     int           `yaml:"port"`
	Domain   string        `yaml:"domain,omitempty"`
	Secure   bool          `yaml:"https"`
	Env      string        `yaml:"env"`
	Bookings client.Config `yaml:"bookings"`
}

const DefaultEnv = "development"

func loadYamlConfigFile(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}

	file, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Config file not found")
			return nil, err
		}
		return nil, err
	}

	config := Config{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	if config.Bookings.AuthURL == "" {
		config.Bookings.AuthURL = "https://login.microsoftonline.com"
	}

	return &config, nil
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	configFile, err := loadYamlConfigFile("config.yaml")
	if err != nil {
		return nil, err
	}

	fmt.Printf("Config: %+v\n", configFile)

	env := GetEnvStringOrDefault("APP_ENV", DefaultEnv)
	utils.PrintInfo(fmt.Sprintf("Running in %s mode", utils.Bold(env)))

	secure := GetEnvBool("HTTP_SECURE")
	utils.PrintInfo(fmt.Sprintf("Secure: %t", secure))

	config := Config{
		Secure: secure,
		Env:    env,
	}

	defaultDomain := "localhost"
	if config.Env != "development" {
		defaultDomain, err = os.Hostname()
		if err != nil {
			return &Config{}, err
		}
	}
	domain := GetEnvStringOrDefault("DOMAIN", defaultDomain)
	config.Domain = domain

	var port int
	portStr := GetEnvInt("PORT")
	if portStr == 0 {
		if config.Secure {
			port = 443
		} else {
			if config.Env == "development" {
				port = 8080
			} else {
				port = 80
			}
		}
	}
	config.Port = port

	var url string
	if config.Secure {
		url = "https://"
	} else {
		url = "http://"
	}
	config.Domain = url + config.Domain

	config.Bookings.ClientID = GetEnvStringOrDefault("BOOKINGS_CLIENT_ID", configFile.Bookings.ClientID)
	config.Bookings.ClientSecret = GetEnvStringOrDefault("BOOKINGS_CLIENT_SECRET", configFile.Bookings.ClientSecret)
	config.Bookings.TenantID = GetEnvStringOrDefault("BOOKINGS_TENANT_ID", configFile.Bookings.TenantID)
	config.Bookings.BusinessID = GetEnvStringOrDefault("BOOKINGS_BUSINESS_ID", configFile.Bookings.BusinessID)
	authUrl := GetEnvStringOrDefault("BOOKINGS_AUTH_URL", configFile.Bookings.AuthURL)

	config.Bookings.AuthURL = authUrl + "/" + config.Bookings.TenantID
	config.Bookings.Scopes = []string{"https://graph.microsoft.com/.default"}

	return &config, nil
}
