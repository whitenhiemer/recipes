package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type generateRequest struct {
	Prompt   string `json:"prompt"`
	Category string `json:"category"`
}

type generateResponse struct {
	Markdown string `json:"markdown"`
	Title    string `json:"title"`
	Category string `json:"category"`
}

const generateSystemPrompt = `You are a creative chef and recipe writer for a family recipe website. Generate a complete, practical, delicious recipe based on the user's request.

Return ONLY the recipe in this exact markdown format — no preamble, no explanation, no code fences, just the raw markdown:

---
tags: [tag1, tag2]
---
# Recipe Title

**Prep Time:** X min | **Cook Time:** X min | **Servings:** X

## Ingredients

- 1 cup ingredient
- 2 tbsp ingredient

## Instructions

1. Step one.
2. Step two.

## Notes

- Tip or variation.

Rules:
- Tags must be lowercase and hyphenated (e.g. quick, chicken, comfort-food, one-pan)
- Use realistic, specific quantities and times
- Include at least 6 ingredients
- Include at least 5 clear instruction steps
- Include 2-3 useful notes (tips, substitutions, storage)
- Servings is a number or range like 4 or 4-6
- Prep/cook time is in minutes like 15 min or 1 hr 30 min`

func (s *Server) handleGenerateRecipe(w http.ResponseWriter, r *http.Request) {
	if s.cfg.AnthropicAPIKey == "" {
		http.Error(w, "AI generation is not configured on this server", http.StatusServiceUnavailable)
		return
	}

	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	req.Prompt = strings.TrimSpace(req.Prompt)
	if req.Prompt == "" {
		http.Error(w, "prompt is required", http.StatusBadRequest)
		return
	}

	userMsg := req.Prompt
	if req.Category != "" {
		userMsg += "\n\nThis recipe is for the category: " + req.Category
	}

	apiBody, _ := json.Marshal(map[string]any{
		"model":      "claude-haiku-4-5-20251001",
		"max_tokens": 1024,
		"system":     generateSystemPrompt,
		"messages":   []map[string]string{{"role": "user", "content": userMsg}},
	})

	apiReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost,
		"https://api.anthropic.com/v1/messages", bytes.NewReader(apiBody))
	if err != nil {
		http.Error(w, "failed to build request", http.StatusInternalServerError)
		return
	}
	apiReq.Header.Set("Content-Type", "application/json")
	apiReq.Header.Set("x-api-key", s.cfg.AnthropicAPIKey)
	apiReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(apiReq)
	if err != nil {
		http.Error(w, "AI request failed: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "AI API error: "+string(respBody), http.StatusBadGateway)
		return
	}

	var apiResp struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil || len(apiResp.Content) == 0 {
		http.Error(w, "failed to parse AI response", http.StatusInternalServerError)
		return
	}

	markdown := strings.TrimSpace(apiResp.Content[0].Text)

	title := ""
	for _, line := range strings.Split(markdown, "\n") {
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimPrefix(line, "# ")
			break
		}
	}

	category := req.Category
	if category == "" {
		category = "dinners"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(generateResponse{
		Markdown: markdown,
		Title:    title,
		Category: category,
	})
}
