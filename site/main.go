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
	"github.com/whitenhiemer/recipe-site/internal/store"
)

func main() {
	cfg := config.Load()

	db, err := store.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	imported, err := db.ImportFromMarkdown(cfg.RecipesDir)
	if err != nil {
		log.Printf("warning: markdown import failed: %v", err)
	} else {
		log.Printf("imported/updated %d recipes from markdown files", imported)
	}

	rows, err := db.GetAllRecipes()
	if err != nil {
		log.Fatalf("failed to load recipes from database: %v", err)
	}
	recipes := store.RowsToRecipes(rows)
	log.Printf("loaded %d recipes from database", len(recipes))

	idx := recipe.NewIndex(recipes)

	mux := http.NewServeMux()
	handler.Register(mux, idx, cfg, db)

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
			log.Println("SIGHUP received, reimporting and rebuilding index...")
			if n, err := db.ImportFromMarkdown(cfg.RecipesDir); err != nil {
				log.Printf("reimport failed: %v", err)
			} else {
				log.Printf("reimported %d recipes", n)
			}
			allRows, err := db.GetAllRecipes()
			if err != nil {
				log.Printf("reload failed: %v", err)
				continue
			}
			idx.Rebuild(store.RowsToRecipes(allRows))
			log.Println("index rebuilt")
		}
	}()

	log.Printf("listening on %s (db: %s)", cfg.Addr, cfg.DBPath)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
