package gogobot

import (
	"net/http"
	"testing"
)

func TestLoad(t *testing.T) {
	detector := Load()

	if detector == nil {
		t.Fatal("Load() returned nil")
	}

	// Should be equivalent to NewDetector()
	directDetector := NewDetector()
	if len(detector.GetDetectorNames()) != len(directDetector.GetDetectorNames()) {
		t.Error("Load() should return detector with same number of detectors as NewDetector()")
	}
}

func TestLoadWithCustomDetectors(t *testing.T) {
	customDetectors := map[string]DetectorFunc{
		"custom1": func(components *ComponentDict) *BotDetectionResult {
			return &BotDetectionResult{Bot: true, BotKind: BotKindUnknown}
		},
		"custom2": func(components *ComponentDict) *BotDetectionResult {
			return &BotDetectionResult{Bot: false}
		},
	}

	detector := LoadWithCustomDetectors(customDetectors)

	if detector == nil {
		t.Fatal("LoadWithCustomDetectors() returned nil")
	}

	names := detector.GetDetectorNames()
	foundCustom1 := false
	foundCustom2 := false

	for _, name := range names {
		if name == "custom1" {
			foundCustom1 = true
		}
		if name == "custom2" {
			foundCustom2 = true
		}
	}

	if !foundCustom1 || !foundCustom2 {
		t.Error("Custom detectors not found in loaded detector")
	}
}

func TestDetect(t *testing.T) {
	tests := []struct {
		name        string
		userAgent   string
		headers     map[string]string
		expectedBot bool
	}{
		{
			name:      "Browser request",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			headers: map[string]string{
				"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"Accept-Language": "en-US,en;q=0.5",
				"Accept-Encoding": "gzip, deflate",
				"Connection":      "keep-alive",
			},
			expectedBot: false,
		},
		{
			name:        "Curl request",
			userAgent:   "curl/7.68.0",
			expectedBot: true,
		},
		{
			name:        "Python requests",
			userAgent:   "python-requests/2.25.1",
			expectedBot: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := createTestHTTPRequest("GET", "/", map[string]string{
				"User-Agent": test.userAgent,
			})

			// Add additional headers if provided
			for k, v := range test.headers {
				req.Header.Set(k, v)
			}

			result, err := Detect(req)

			if err != nil {
				t.Fatalf("Detect() returned error: %v", err)
			}

			if result.Bot != test.expectedBot {
				t.Errorf("Expected bot=%t, got bot=%t", test.expectedBot, result.Bot)
			}
		})
	}
}

func TestDetectWithCustomDetectors(t *testing.T) {
	// Custom detector that always detects bots
	alwaysBot := func(components *ComponentDict) *BotDetectionResult {
		return &BotDetectionResult{Bot: true, BotKind: BotKindUnknown}
	}

	customDetectors := map[string]DetectorFunc{
		"alwaysBot": alwaysBot,
	}

	req := createTestHTTPRequest("GET", "/", map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language": "en-US,en;q=0.5",
		"Accept-Encoding": "gzip, deflate",
		"Connection":      "keep-alive",
	})

	// Without custom detector
	result1, err := Detect(req)
	if err != nil {
		t.Fatalf("Detect() returned error: %v", err)
	}

	// With custom detector
	result2, err := DetectWithCustomDetectors(req, customDetectors)
	if err != nil {
		t.Fatalf("DetectWithCustomDetectors() returned error: %v", err)
	}

	// The custom detector should make it detect as bot
	if !result2.Bot {
		t.Error("Expected custom detector to detect as bot")
	}

	// Results should be different due to custom detector
	if result1.Bot == result2.Bot {
		t.Error("Expected different results with and without custom detector")
	}
}

func TestQuickCheck(t *testing.T) {
	tests := []struct {
		name        string
		userAgent   string
		headers     map[string]string
		expectedBot bool
	}{
		{
			name:      "Browser request",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			headers: map[string]string{
				"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
				"Accept-Language": "en-US,en;q=0.5",
				"Accept-Encoding": "gzip, deflate",
				"Connection":      "keep-alive",
			},
			expectedBot: false,
		},
		{
			name:        "Curl request",
			userAgent:   "curl/7.68.0",
			expectedBot: true,
		},
		{
			name:        "Missing User-Agent",
			userAgent:   "",
			expectedBot: true, // Should be detected by missing headers check
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := createTestHTTPRequest("GET", "/", map[string]string{})
			if test.userAgent != "" {
				req.Header.Set("User-Agent", test.userAgent)
			}

			// Add additional headers if provided
			for k, v := range test.headers {
				req.Header.Set(k, v)
			}

			result, err := QuickCheck(req)

			if err != nil {
				t.Fatalf("QuickCheck() returned error: %v", err)
			}

			if result.Bot != test.expectedBot {
				t.Errorf("Expected bot=%t, got bot=%t", test.expectedBot, result.Bot)
			}
		})
	}
}

func TestIsBotUserAgent(t *testing.T) {
	tests := []struct {
		userAgent    string
		expectedBot  bool
		expectedKind BotKind
	}{
		{
			userAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			expectedBot: false,
		},
		{
			userAgent:    "curl/7.68.0",
			expectedBot:  true,
			expectedKind: BotKindCurl,
		},
		{
			userAgent:    "python-requests/2.25.1",
			expectedBot:  true,
			expectedKind: BotKindUnknown,
		},
		{
			userAgent:    "PhantomJS/2.1.1",
			expectedBot:  true,
			expectedKind: BotKindPhantomJS,
		},
		{
			userAgent:    "Googlebot/2.1 (+http://www.google.com/bot.html)",
			expectedBot:  true,
			expectedKind: BotKindCrawler, // Changed from BotKindBot to BotKindCrawler
		},
		{
			userAgent:   "",
			expectedBot: false, // Empty user agent won't be detected by user agent check alone
		},
	}

	for _, test := range tests {
		t.Run(test.userAgent, func(t *testing.T) {
			isBot, botKind := IsBotUserAgent(test.userAgent)

			if isBot != test.expectedBot {
				t.Errorf("Expected bot=%t, got bot=%t for UA: %s", test.expectedBot, isBot, test.userAgent)
			}

			if test.expectedBot && botKind != test.expectedKind {
				t.Errorf("Expected bot kind %s, got %s for UA: %s", test.expectedKind, botKind, test.userAgent)
			}
		})
	}
}

func TestAnalyzeHeaders(t *testing.T) {
	tests := []struct {
		name        string
		headers     map[string][]string
		expectedBot bool
	}{
		{
			name: "Normal browser headers",
			headers: map[string][]string{
				"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"},
				"Accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
				"Accept-Language": {"en-US,en;q=0.5"},
				"Accept-Encoding": {"gzip, deflate"},
				"Connection":      {"keep-alive"},
			},
			expectedBot: false,
		},
		{
			name: "Missing common headers",
			headers: map[string][]string{
				"User-Agent": {"Mozilla/5.0"},
				// Missing Accept, Accept-Language, Accept-Encoding, Connection
			},
			expectedBot: true,
		},
		{
			name: "Automation headers",
			headers: map[string][]string{
				"User-Agent":       {"Mozilla/5.0"},
				"X-Requested-With": {"XMLHttpRequest"},
			},
			expectedBot: true,
		},
		{
			name: "Too few headers",
			headers: map[string][]string{
				"User-Agent": {"test"},
			},
			expectedBot: true,
		},
		{
			name: "Suspicious accept pattern",
			headers: map[string][]string{
				"User-Agent": {"Mozilla/5.0"},
				"Accept":     {"*/*"},
				// Missing Accept-Language
			},
			expectedBot: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := AnalyzeHeaders(test.headers)

			if result.Bot != test.expectedBot {
				t.Errorf("Expected bot=%t, got bot=%t", test.expectedBot, result.Bot)
			}
		})
	}
}

// Helper function to create test HTTP requests for API tests
func createTestHTTPRequest(method, path string, headers map[string]string) *http.Request {
	req, _ := http.NewRequest(method, path, nil)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req
}
