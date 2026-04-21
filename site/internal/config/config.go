package config

import (
	"flag"
	"log"
	"os"
	"strings"
)

type Config struct {
	Addr          string
	RecipesDir    string
	WebhookSecret string
}

func Load() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Addr, "addr", ":8080", "listen address")
	flag.StringVar(&cfg.RecipesDir, "recipes", "..", "path to recipes directory")

	secretFile := ""
	flag.StringVar(&secretFile, "webhook-secret-file", "", "path to file containing webhook secret")
	flag.Parse()

	if secretFile != "" {
		data, err := os.ReadFile(secretFile)
		if err != nil {
			log.Printf("warning: could not read webhook secret file: %v", err)
		} else {
			cfg.WebhookSecret = strings.TrimSpace(string(data))
		}
	}

	if env := os.Getenv("WEBHOOK_SECRET"); env != "" && cfg.WebhookSecret == "" {
		cfg.WebhookSecret = env
	}

	return cfg
}
