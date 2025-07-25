package main

import (
	"fmt"
	"net/http"

	"github.com/lytics/gogobot"
)

func main() {
	// Example 1: Basic detection from a request
	fmt.Println("=== Basic Bot Detection Example ===")

	// Create a new bot detector
	detector := gogobot.Load()

	// Simulate a request from Chrome browser
	req1, _ := http.NewRequest("GET", "/", nil)
	req1.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req1.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req1.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req1.Header.Set("Accept-Encoding", "gzip, deflate")
	req1.Header.Set("Connection", "keep-alive")

	result1, err := detector.DetectFromRequest(req1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Chrome browser request - Bot: %t\n", result1.Bot)
	if result1.Bot {
		fmt.Printf("Bot kind: %s\n", result1.BotKind)
	}

	// Simulate a request from curl
	req2, _ := http.NewRequest("GET", "/", nil)
	req2.Header.Set("User-Agent", "curl/7.68.0")
	req2.Header.Set("Accept", "*/*")

	result2, err := detector.DetectFromRequest(req2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Curl request - Bot: %t\n", result2.Bot)
	if result2.Bot {
		fmt.Printf("Bot kind: %s\n", result2.BotKind)
	}

	// Simulate a request from Selenium
	req3, _ := http.NewRequest("GET", "/", nil)
	req3.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/91.0.4472.124 Safari/537.36")

	result3, err := detector.DetectFromRequest(req3)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Headless Chrome request - Bot: %t\n", result3.Bot)
	if result3.Bot {
		fmt.Printf("Bot kind: %s\n", result3.BotKind)
	}

	// Example 2: Quick user agent check
	fmt.Println("\n=== Quick User Agent Check ===")

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"curl/7.68.0",
		"python-requests/2.25.1",
		"Googlebot/2.1 (+http://www.google.com/bot.html)",
		"PhantomJS/2.1.1",
	}

	for _, ua := range userAgents {
		isBot, botKind := gogobot.IsBotUserAgent(ua)
		fmt.Printf("UA: %s -> Bot: %t", ua, isBot)
		if isBot {
			fmt.Printf(" (Kind: %s)", botKind)
		}
		fmt.Println()
	}

	// Example 3: Analyze headers
	fmt.Println("\n=== Header Analysis ===")

	suspiciousHeaders := map[string][]string{
		"User-Agent": {"Mozilla/5.0"},
		// Missing common headers like Accept-Language, Accept-Encoding
	}

	headerResult := gogobot.AnalyzeHeaders(suspiciousHeaders)
	fmt.Printf("Suspicious headers - Bot: %t\n", headerResult.Bot)
	if headerResult.Bot {
		fmt.Printf("Bot kind: %s\n", headerResult.BotKind)
	}
}
