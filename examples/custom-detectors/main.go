package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lytics/gogobot"
)

func main() {
	fmt.Println("=== Custom Detectors Example ===")

	// Define custom detectors
	customDetectors := map[string]gogobot.DetectorFunc{
		"suspiciousIp":      detectSuspiciousIP,
		"rapidRequests":     detectRapidRequests,
		"missingReferer":    detectMissingReferer,
		"automationHeaders": detectAutomationHeaders,
	}

	// Create detector with custom detectors
	detector := gogobot.LoadWithCustomDetectors(customDetectors)

	// Test various scenarios
	scenarios := []struct {
		name         string
		setupRequest func() *http.Request
	}{
		{
			name: "Normal browser request",
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
				req.Header.Set("Referer", "https://google.com/")
				req.RemoteAddr = "192.168.1.100:12345"
				return req
			},
		},
		{
			name: "Request from datacenter IP",
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")
				req.RemoteAddr = "54.239.123.45:80" // AWS IP range
				return req
			},
		},
		{
			name: "Request with automation headers",
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")
				req.Header.Set("X-Requested-With", "XMLHttpRequest")
				req.Header.Set("Chrome-Proxy", "frfr")
				return req
			},
		},
		{
			name: "Direct access without referer",
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/important-page", nil)
				req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")
				// No referer header
				return req
			},
		},
		{
			name: "Python requests library",
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/", nil)
				req.Header.Set("User-Agent", "python-requests/2.25.1")
				req.Header.Set("Accept", "*/*")
				req.Header.Set("Accept-Encoding", "gzip, deflate")
				return req
			},
		},
	}

	for _, scenario := range scenarios {
		fmt.Printf("\n--- %s ---\n", scenario.name)

		req := scenario.setupRequest()
		result, err := detector.DetectFromRequest(req)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Bot detected: %t\n", result.Bot)
		if result.Bot {
			fmt.Printf("Bot kind: %s\n", result.BotKind)
		}

		// Show individual detector results
		detections := detector.GetDetections()
		if detections != nil {
			if detections.UserAgent.Bot {
				fmt.Printf("  - User Agent: Bot (%s)\n", detections.UserAgent.BotKind)
			}
			if detections.Headers.Bot {
				fmt.Printf("  - Headers: Bot (%s)\n", detections.Headers.BotKind)
			}
			if detections.MissingHeaders.Bot {
				fmt.Printf("  - Missing Headers: Bot\n")
			}
			// Note: Custom detectors would need to be added to DetectionDict
			// or we'd need a different way to access their results
		}
	}

	// Example of adding a detector at runtime
	fmt.Printf("\n=== Runtime Detector Addition ===\n")

	detector.AddDetector("lowEntropy", detectLowEntropyUserAgent)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0") // Very generic UA

	result, _ := detector.DetectFromRequest(req)
	fmt.Printf("Low entropy UA - Bot detected: %t\n", result.Bot)

	// Show available detectors
	fmt.Printf("\nActive detectors: %v\n", detector.GetDetectorNames())
}

// Custom detector: Check if IP is from a datacenter/cloud provider
func detectSuspiciousIP(components *gogobot.ComponentDict) *gogobot.BotDetectionResult {
	if components.RemoteAddr.GetState() != gogobot.StateSuccess {
		return &gogobot.BotDetectionResult{Bot: false}
	}

	remoteAddr := components.RemoteAddr.GetValue()
	ip := strings.Split(remoteAddr, ":")[0]

	// Simplified check for known datacenter IP ranges
	datacenters := []string{
		"54.239.",  // AWS
		"52.0.",    // AWS
		"104.154.", // Google Cloud
		"35.184.",  // Google Cloud
		"40.76.",   // Azure
		"13.107.",  // Azure
	}

	for _, dcRange := range datacenters {
		if strings.HasPrefix(ip, dcRange) {
			return &gogobot.BotDetectionResult{
				Bot:     true,
				BotKind: gogobot.BotKindUnknown,
			}
		}
	}

	return &gogobot.BotDetectionResult{Bot: false}
}

// Custom detector: Detect rapid requests (would need state management in real implementation)
func detectRapidRequests(components *gogobot.ComponentDict) *gogobot.BotDetectionResult {
	// This is a simplified example - real implementation would need
	// rate limiting state management
	return &gogobot.BotDetectionResult{Bot: false}
}

// Custom detector: Detect missing referer on important pages
func detectMissingReferer(components *gogobot.ComponentDict) *gogobot.BotDetectionResult {
	path := components.RequestPath.GetValue()

	// Check if this is an important page that should have a referer
	importantPages := []string{"/checkout", "/important-page", "/admin"}

	isImportantPage := false
	for _, page := range importantPages {
		if strings.Contains(path, page) {
			isImportantPage = true
			break
		}
	}

	if !isImportantPage {
		return &gogobot.BotDetectionResult{Bot: false}
	}

	// Check if referer is missing from headers
	headers := components.Headers.GetValue()
	if _, hasReferer := headers["Referer"]; !hasReferer {
		return &gogobot.BotDetectionResult{
			Bot:     true,
			BotKind: gogobot.BotKindUnknown,
		}
	}

	return &gogobot.BotDetectionResult{Bot: false}
}

// Custom detector: Detect automation-specific headers
func detectAutomationHeaders(components *gogobot.ComponentDict) *gogobot.BotDetectionResult {
	if components.Headers.GetState() != gogobot.StateSuccess {
		return &gogobot.BotDetectionResult{Bot: false}
	}

	headers := components.Headers.GetValue()

	// Check for automation-specific headers
	automationHeaders := []string{
		"X-Requested-With",
		"Chrome-Proxy",
		"Purpose",
		"X-DevTools-Emulate-Network-Conditions-Client-Id",
		"X-Chrome-UMA-Enabled",
		"X-Client-Data",
	}

	for _, header := range automationHeaders {
		if _, exists := headers[header]; exists {
			return &gogobot.BotDetectionResult{
				Bot:     true,
				BotKind: gogobot.BotKindUnknown,
			}
		}
	}

	return &gogobot.BotDetectionResult{Bot: false}
}

// Custom detector: Detect low-entropy user agents
func detectLowEntropyUserAgent(components *gogobot.ComponentDict) *gogobot.BotDetectionResult {
	if components.UserAgent.GetState() != gogobot.StateSuccess {
		return &gogobot.BotDetectionResult{Bot: false}
	}

	userAgent := components.UserAgent.GetValue()

	// Very simple entropy check - count unique characters
	seen := make(map[rune]bool)
	for _, char := range userAgent {
		seen[char] = true
	}

	// If user agent has very few unique characters, it might be generic/fake
	uniqueChars := len(seen)
	if len(userAgent) > 20 && uniqueChars < 10 {
		return &gogobot.BotDetectionResult{
			Bot:     true,
			BotKind: gogobot.BotKindUnknown,
		}
	}

	return &gogobot.BotDetectionResult{Bot: false}
}
