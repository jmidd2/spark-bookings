package config

import (
	"fmt"
	"log"
	"os"
	"spark-bookings/client"
	"spark-bookings/utils"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     int
	Domain   string
	Secure   bool
	Env      string
	Bookings client.Config
}

const DefaultEnv = "development"

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	config.Bookings.ClientID = GetEnvStringRequired("BOOKINGS_CLIENT_ID")
	config.Bookings.ClientSecret = GetEnvStringRequired("BOOKINGS_CLIENT_SECRET")
	config.Bookings.TenantID = GetEnvStringRequired("BOOKINGS_TENANT_ID")
	config.Bookings.BusinessID = GetEnvStringRequired("BOOKINGS_BUSINESS_ID")
	authUrl := GetEnvStringOrDefault("BOOKINGS_AUTH_URL", "https://login.microsoftonline.com")

	config.Bookings.AuthURL = authUrl + "/" + config.Bookings.TenantID
	config.Bookings.Scopes = []string{"https://graph.microsoft.com/.default"}

	return &config, nil
}
