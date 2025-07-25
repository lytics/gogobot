package main

import (
	"fmt"
	"net/http"

	"gogobot"
)

func main() {
	// Example 1: Parse browser from user agent string
	fmt.Println("=== Browser Parsing Examples ===\n")

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 14; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"curl/7.81.0",
	}

	for i, ua := range userAgents {
		fmt.Printf("Example %d:\n", i+1)
		fmt.Printf("User Agent: %s\n", ua)

		browserInfo := gogobot.ParseBrowser(ua)

		fmt.Printf("Browser: %s\n", browserInfo.Name)
		fmt.Printf("Version: %s\n", browserInfo.Version)
		fmt.Printf("Major Version: %s\n", browserInfo.GetMajorVersion())
		fmt.Printf("Browser Family: %s\n", browserInfo.GetBrowserFamily())
		fmt.Printf("Is Mobile: %v\n", browserInfo.IsMobile())
		fmt.Printf("Is Bot: %v\n", browserInfo.IsBot)
		if browserInfo.IsBot {
			fmt.Printf("Bot Kind: %s\n", browserInfo.BotKind)
		}
		fmt.Println()
	}

	// Example 2: HTTP Request Analysis
	fmt.Println("=== HTTP Request Analysis ===\n")

	// Simulate an HTTP request
	req, _ := http.NewRequest("GET", "/api/data", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// Get comprehensive browser and bot information
	browserInfo, botResult, err := gogobot.GetBrowserInfo(req)
	if err != nil {
		fmt.Printf("Error analyzing request: %v\n", err)
		return
	}

	fmt.Printf("Request Analysis Results:\n")
	fmt.Printf("Browser: %s %s\n", browserInfo.Name, browserInfo.Version)
	fmt.Printf("Browser Family: %s\n", browserInfo.GetBrowserFamily())
	fmt.Printf("Is Mobile: %v\n", browserInfo.IsMobile())
	fmt.Printf("Is Bot (Browser Analysis): %v\n", browserInfo.IsBot)
	fmt.Printf("Is Bot (Full Detection): %v\n", botResult.Bot)
	if botResult.Bot {
		fmt.Printf("Bot Type: %s\n", botResult.BotKind)
	}
	fmt.Println()

	// Example 3: Browser Support Checking
	fmt.Println("=== Browser Support Checking ===\n")

	// Define minimum browser versions for your application
	minVersions := map[gogobot.BrowserName]string{
		gogobot.BrowserChrome:  "100.0.0.0",
		gogobot.BrowserFirefox: "100.0.0.0",
		gogobot.BrowserSafari:  "15.0",
		gogobot.BrowserEdge:    "100.0.0.0",
	}

	testBrowsers := []struct {
		name      string
		userAgent string
	}{
		{
			name:      "Supported Chrome",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		},
		{
			name:      "Unsupported Chrome",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36",
		},
		{
			name:      "Supported Safari",
			userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		},
		{
			name:      "Unsupported Browser",
			userAgent: "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; Trident/4.0)",
		},
	}

	for _, test := range testBrowsers {
		req.Header.Set("User-Agent", test.userAgent)
		browserInfo := gogobot.ParseBrowserFromHTTPRequest(req)
		isSupported := gogobot.IsSupportedBrowser(req, minVersions)

		fmt.Printf("%s:\n", test.name)
		fmt.Printf("  Browser: %s %s\n", browserInfo.Name, browserInfo.Version)
		fmt.Printf("  Supported: %v\n", isSupported)
		fmt.Println()
	}

	// Example 4: Mobile Detection
	fmt.Println("=== Mobile Detection ===\n")

	mobileTestCases := []struct {
		name      string
		userAgent string
	}{
		{
			name:      "iPhone Safari",
			userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
		},
		{
			name:      "Android Chrome",
			userAgent: "Mozilla/5.0 (Linux; Android 14; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
		},
		{
			name:      "Desktop Chrome",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		},
		{
			name:      "iPad Safari",
			userAgent: "Mozilla/5.0 (iPad; CPU OS 17_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
		},
	}

	for _, test := range mobileTestCases {
		req.Header.Set("User-Agent", test.userAgent)
		isMobile := gogobot.IsMobileBrowser(req)
		browserInfo := gogobot.ParseBrowserFromHTTPRequest(req)

		fmt.Printf("%s:\n", test.name)
		fmt.Printf("  Browser: %s %s\n", browserInfo.Name, browserInfo.Version)
		fmt.Printf("  Is Mobile: %v\n", isMobile)
		fmt.Printf("  Browser Family: %s\n", browserInfo.GetBrowserFamily())
		fmt.Println()
	}

	// Example 5: Combined Bot Detection and Browser Parsing
	fmt.Println("=== Combined Analysis ===\n")

	// Create a bot detector
	detector := gogobot.NewDetector()

	combinedTestCases := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"Mozilla/5.0 (Unknown; Linux x86_64) AppleWebKit/538.1 (KHTML, like Gecko) PhantomJS/2.1.1 Safari/538.1",
		"curl/7.81.0",
	}

	for i, ua := range combinedTestCases {
		req.Header.Set("User-Agent", ua)

		// Parse browser information
		browserInfo := gogobot.ParseBrowserFromHTTPRequest(req)

		// Perform bot detection
		botResult, err := detector.DetectFromRequest(req)
		if err != nil {
			fmt.Printf("Error detecting bot: %v\n", err)
			continue
		}

		fmt.Printf("Test Case %d:\n", i+1)
		fmt.Printf("  User Agent: %s\n", ua)
		fmt.Printf("  Browser: %s %s\n", browserInfo.Name, browserInfo.Version)
		fmt.Printf("  Browser Family: %s\n", browserInfo.GetBrowserFamily())
		fmt.Printf("  Is Bot (Browser Parser): %v\n", browserInfo.IsBot)
		fmt.Printf("  Is Bot (Full Detector): %v\n", botResult.Bot)
		if botResult.Bot {
			fmt.Printf("  Bot Type: %s\n", botResult.BotKind)
		}
		fmt.Println()
	}
}
