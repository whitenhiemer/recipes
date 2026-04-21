package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	lastDeploy   time.Time
	deployMu     sync.Mutex
	deployWindow = 30 * time.Second
)

func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if s.cfg.WebhookSecret == "" {
		http.Error(w, "webhook not configured", http.StatusServiceUnavailable)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	sig := r.Header.Get("X-Hub-Signature-256")
	if !verifySignature(body, sig, s.cfg.WebhookSecret) {
		http.Error(w, "invalid signature", http.StatusForbidden)
		return
	}

	deployMu.Lock()
	if time.Since(lastDeploy) < deployWindow {
		deployMu.Unlock()
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	lastDeploy = time.Now()
	deployMu.Unlock()

	go func() {
		log.Println("webhook: pulling latest recipes...")
		cmd := exec.Command("git", "-C", s.cfg.RecipesDir, "pull", "origin", "main")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("webhook: git pull failed: %v\n%s", err, output)
			return
		}
		log.Printf("webhook: git pull: %s", output)

		imported, err := s.db.ImportFromMarkdown(s.cfg.RecipesDir)
		if err != nil {
			log.Printf("webhook: import failed: %v", err)
			return
		}
		log.Printf("webhook: imported/updated %d recipes from markdown", imported)

		if err := s.rebuildIndex(); err != nil {
			log.Printf("webhook: index rebuild failed: %v", err)
			return
		}
		log.Println("webhook: index rebuilt successfully")
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func verifySignature(payload []byte, signature string, secret string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}
	sig, err := hex.DecodeString(strings.TrimPrefix(signature, "sha256="))
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := mac.Sum(nil)

	return hmac.Equal(sig, expected)
}
