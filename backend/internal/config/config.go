package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DBDSN        string
	JWTSecret    string
	Port         string
	CookieDomain string
	CookieSecure bool
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	viper.AutomaticEnv()

	dsn := viper.GetString("DB_DSN")
	if dsn == "" {
		log.Fatal("FATAL: DB_DSN not set in .env")
	}

	secret := viper.GetString("JWT_SECRET")
	if secret == "" || secret == "CHANGE_ME_MIN_32_CHARS" {
		log.Fatal("FATAL: JWT_SECRET not set in .env")
	}
	if len(secret) < 32 {
		log.Fatal("FATAL: JWT_SECRET must be at least 32 characters")
	}

	port := viper.GetString("SERVER_PORT")
	if port == "" {
		port = ":8080"
	}

	return &Config{
		DBDSN:        dsn,
		JWTSecret:    secret,
		Port:         port,
		CookieDomain: viper.GetString("COOKIE_DOMAIN"),
		CookieSecure: viper.GetBool("COOKIE_SECURE"),
	}
}