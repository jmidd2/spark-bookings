package config

import (
	"fmt"
	"log"
	"os"
	"spark-bookings/client"
	"spark-bookings/utils"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     int
	Domain   string
	Secure   bool
	Env      string
	Bookings client.Config
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	utils.PrintInfo(fmt.Sprintf("Running in %s mode", utils.Bold(env)))
	//fmt.Printf("%sRunning in %s mode%s\n", utils.ColorBlue, utils.Bold(env), utils.ColorReset)

	config := Config{
		Secure: os.Getenv("HTTP_SECURE") == "true",
		Env:    env,
	}

	domain := os.Getenv("DOMAIN")
	if domain != "" {
		config.Domain = domain
	} else {
		if config.Env == "development" {
			config.Domain = "localhost"
		} else {
			return &Config{}, fmt.Errorf("DOMAIN environment variable is required")
		}
	}

	var url string
	var port int
	if config.Secure {
		url = "https://"
		portStr := os.Getenv("PORT")
		port = 443
		if portStr != "" {
			var err error
			port, err = strconv.Atoi(portStr)
			if err != nil {
				return &Config{}, err
			}
		}
	} else {
		url = "http://"
		portStr := os.Getenv("PORT")
		if config.Env == "development" {
			port = 8080
		} else {
			port = 80
		}
		if portStr != "" {
			var err error
			port, err = strconv.Atoi(portStr)
			if err != nil {
				return &Config{}, err
			}
		}
	}
	config.Port = port
	config.Domain = url + config.Domain

	//config.Bookings.ClientID = os.Getenv("BOOKINGS_CLIENT_ID")
	config.Bookings.ClientID = GetEnvStringRequired("BOOKINGS_CLIENT_ID")
	//config.Bookings.ClientSecret = os.Getenv("BOOKINGS_CLIENT_SECRET")
	config.Bookings.ClientSecret = GetEnvStringRequired("BOOKINGS_CLIENT_SECRET")
	//config.Bookings.TenantID = os.Getenv("BOOKINGS_TENANT_ID")
	config.Bookings.TenantID = GetEnvStringRequired("BOOKINGS_TENANT_ID")
	//config.Bookings.BusinessID = os.Getenv("BOOKINGS_BUSINESS_ID")
	config.Bookings.BusinessID = GetEnvStringRequired("BOOKINGS_BUSINESS_ID")
	//authUrl := os.Getenv("BOOKINGS_AUTH_URL")
	authUrl := GetEnvStringOrDefault("BOOKINGS_AUTH_URL", "https://login.microsoftonline.com")

	config.Bookings.AuthURL = authUrl + "/" + config.Bookings.TenantID
	config.Bookings.Scopes = []string{"https://graph.microsoft.com/.default"}

	return &config, nil
}
