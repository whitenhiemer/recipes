package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/whitenhiemer/recipe-site/internal/config"
	"github.com/whitenhiemer/recipe-site/internal/handler"
	"github.com/whitenhiemer/recipe-site/internal/recipe"
)

func main() {
	cfg := config.Load()

	recipes, err := recipe.ParseRecipesDir(cfg.RecipesDir)
	if err != nil {
		log.Fatalf("failed to parse recipes: %v", err)
	}
	log.Printf("loaded %d recipes from %s", len(recipes), cfg.RecipesDir)

	idx := recipe.NewIndex(recipes)

	mux := http.NewServeMux()
	handler.Register(mux, idx, cfg)

	wrapped := handler.Chain(mux,
		handler.Recovery,
		handler.SecurityHeaders,
		handler.RequestLogger,
	)

	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      wrapped,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP)
		for range sigs {
			log.Println("SIGHUP received, reloading index...")
			if err := idx.Reload(cfg.RecipesDir); err != nil {
				log.Printf("reload failed: %v", err)
			} else {
				log.Println("index reloaded")
			}
		}
	}()

	log.Printf("listening on %s", cfg.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
