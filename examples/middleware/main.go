package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lytics/gogobot"
)

func main() {
	// Create a bot detector
	detector := gogobot.Load()

	fmt.Println("Starting server with bot detection middleware...")

	// Example 1: Basic middleware that blocks bots
	protectedConfig := gogobot.MiddlewareConfig{
		BlockBots:         true,
		BlockedStatusCode: http.StatusForbidden,
		BlockedMessage:    "Bot traffic is not allowed on this endpoint",
		OnBotDetected: func(w http.ResponseWriter, r *http.Request, result *gogobot.BotDetectionResult) {
			log.Printf("Bot detected: %s from %s (Kind: %s)",
				r.Header.Get("User-Agent"),
				r.RemoteAddr,
				result.BotKind)
		},
		OnError: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Bot detection error: %v", err)
		},
	}

	// Apply middleware to specific routes
	http.Handle("/protected", detector.MiddlewareWithConfig(protectedConfig)(
		http.HandlerFunc(protectedHandler),
	))

	// Example 2: API endpoint that returns detection results
	http.Handle("/api/detect", detector.Middleware()(
		http.HandlerFunc(detectHandler),
	))

	// Example 3: Public endpoint without bot blocking
	publicConfig := gogobot.MiddlewareConfig{
		BlockBots: false,
		OnBotDetected: func(w http.ResponseWriter, r *http.Request, result *gogobot.BotDetectionResult) {
			log.Printf("Bot detected on public endpoint: %s (Kind: %s)",
				r.Header.Get("User-Agent"),
				result.BotKind)
		},
	}

	http.Handle("/public", detector.MiddlewareWithConfig(publicConfig)(
		http.HandlerFunc(publicHandler),
	))

	// Example 4: Conditional bot blocking
	conditionalConfig := gogobot.MiddlewareConfig{
		SkipFunc: func(r *http.Request) bool {
			// Skip detection for requests with a special header
			return r.Header.Get("X-Skip-Bot-Detection") == "true"
		},
		OnBotDetected: func(w http.ResponseWriter, r *http.Request, result *gogobot.BotDetectionResult) {
			// Only block malicious bots, allow search engine crawlers
			if result.BotKind == gogobot.BotKindCrawler {
				log.Printf("Search engine crawler allowed: %s", r.Header.Get("User-Agent"))
				return
			}

			log.Printf("Malicious bot blocked: %s (Kind: %s)",
				r.Header.Get("User-Agent"),
				result.BotKind)
			http.Error(w, "Automated traffic blocked", http.StatusForbidden)
		},
	}

	http.Handle("/smart-protection", detector.MiddlewareWithConfig(conditionalConfig)(
		http.HandlerFunc(smartProtectedHandler),
	))

	// Health check endpoint without any protection
	http.HandleFunc("/health", healthHandler)

	fmt.Println("Server starting on :8080")
	fmt.Println("Try these endpoints:")
	fmt.Println("  curl http://localhost:8080/protected (should be blocked)")
	fmt.Println("  curl http://localhost:8080/api/detect (returns detection results)")
	fmt.Println("  curl http://localhost:8080/public (allowed but logged)")
	fmt.Println("  curl -H 'User-Agent: Googlebot/2.1' http://localhost:8080/smart-protection (crawler allowed)")
	fmt.Println("  curl -H 'X-Skip-Bot-Detection: true' http://localhost:8080/smart-protection (skipped)")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Welcome! You passed the bot detection.",
		"endpoint": "/protected",
	})
}

func detectHandler(w http.ResponseWriter, r *http.Request) {
	// Get detection results from context
	result, ok := gogobot.GetResultFromContext(r.Context())
	if !ok {
		http.Error(w, "Detection results not available", http.StatusInternalServerError)
		return
	}

	components, _ := gogobot.GetComponentsFromContext(r.Context())

	response := map[string]interface{}{
		"bot":            result.Bot,
		"userAgent":      components.UserAgent.GetValue(),
		"headerCount":    components.HeaderCount.GetValue(),
		"missingHeaders": components.MissingCommonHeaders.GetValue(),
	}

	if result.Bot {
		response["botKind"] = result.BotKind
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func publicHandler(w http.ResponseWriter, r *http.Request) {
	result, _ := gogobot.GetResultFromContext(r.Context())

	response := map[string]interface{}{
		"message":  "This is a public endpoint",
		"endpoint": "/public",
	}

	if result != nil && result.Bot {
		response["note"] = fmt.Sprintf("Bot detected (%s) but allowed", result.BotKind)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func smartProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Smart protection: crawlers allowed, bots blocked",
		"endpoint": "/smart-protection",
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "gogobot",
	})
}
