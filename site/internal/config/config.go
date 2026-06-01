package config

import (
	"flag"
	"log"
	"os"
	"strings"
)

type Config struct {
	Addr            string
	RecipesDir      string
	DBPath          string
	WebhookSecret   string
	AnthropicAPIKey string
}

func Load() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Addr, "addr", ":8080", "listen address")
	flag.StringVar(&cfg.RecipesDir, "recipes", "..", "path to recipes directory")
	flag.StringVar(&cfg.DBPath, "db", "recipes.db", "path to SQLite database file")

	secretFile := ""
	flag.StringVar(&secretFile, "webhook-secret-file", "", "path to file containing webhook secret")
	anthropicKeyFile := ""
	flag.StringVar(&anthropicKeyFile, "anthropic-key-file", "", "path to file containing Anthropic API key")
	flag.Parse()

	if secretFile != "" {
		data, err := os.ReadFile(secretFile)
		if err != nil {
			log.Printf("warning: could not read webhook secret file: %v", err)
		} else {
			cfg.WebhookSecret = strings.TrimSpace(string(data))
		}
	}

	if anthropicKeyFile != "" {
		data, err := os.ReadFile(anthropicKeyFile)
		if err != nil {
			log.Printf("warning: could not read Anthropic API key file: %v", err)
		} else {
			cfg.AnthropicAPIKey = strings.TrimSpace(string(data))
		}
	}

	if env := os.Getenv("WEBHOOK_SECRET"); env != "" && cfg.WebhookSecret == "" {
		cfg.WebhookSecret = env
	}

	if env := os.Getenv("ANTHROPIC_API_KEY"); env != "" && cfg.AnthropicAPIKey == "" {
		cfg.AnthropicAPIKey = env
	}

	return cfg
}
