package gogobot

import (
	"net/http"
)

// Load creates and returns a new BotDetector instance
// This mirrors the JavaScript API's load() function
func Load() *BotDetector {
	return NewDetector()
}

// LoadWithCustomDetectors creates a new BotDetector with custom detectors
func LoadWithCustomDetectors(customDetectors map[string]DetectorFunc) *BotDetector {
	return NewDetectorWithCustomDetectors(customDetectors)
}

// Detect performs bot detection on an HTTP request using a default detector
// This is a convenience function for one-off detections
func Detect(req *http.Request) (BotDetectionResult, error) {
	detector := NewDetector()
	return detector.DetectFromRequest(req)
}

// DetectWithCustomDetectors performs bot detection using custom detectors
func DetectWithCustomDetectors(req *http.Request, customDetectors map[string]DetectorFunc) (BotDetectionResult, error) {
	detector := NewDetectorWithCustomDetectors(customDetectors)
	return detector.DetectFromRequest(req)
}

// QuickCheck performs a fast bot detection check focusing on the most reliable signals
func QuickCheck(req *http.Request) (BotDetectionResult, error) {
	// Create a detector with only the most reliable detectors for speed
	quickDetectors := map[string]DetectorFunc{
		"userAgent":      detectUserAgent,
		"missingHeaders": detectMissingHeaders,
	}

	detector := NewDetectorWithCustomDetectors(quickDetectors)
	return detector.DetectFromRequest(req)
}

// IsBotUserAgent checks if a user agent string indicates a bot
// This is a utility function for checking user agents without a full HTTP request
func IsBotUserAgent(userAgent string) (bool, BotKind) {
	// Create a minimal HTTP request just for user agent analysis
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", userAgent)

	detector := NewDetector()
	components, _ := detector.Collect(req)

	result := detectUserAgent(components)
	if result.Bot {
		return true, result.BotKind
	}

	return false, ""
}

// AnalyzeHeaders performs detailed analysis of HTTP headers
func AnalyzeHeaders(headers map[string][]string) BotDetectionResult {
	// Create a minimal HTTP request with the provided headers
	req, _ := http.NewRequest("GET", "/", nil)
	for k, v := range headers {
		for _, val := range v {
			req.Header.Add(k, val)
		}
	}

	detector := NewDetector()
	components, _ := detector.Collect(req)

	// Run header-specific detectors
	headerDetectors := map[string]DetectorFunc{
		"headers":        detectHeaders,
		"headerCount":    detectHeaderCount,
		"missingHeaders": detectMissingHeaders,
		"acceptHeaders":  detectAcceptHeaders,
		"connection":     detectConnection,
	}

	for _, detectorFunc := range headerDetectors {
		result := detectorFunc(components)
		if result != nil && result.Bot {
			return *result
		}
	}

	return BotDetectionResult{Bot: false}
}

// ParseBrowser extracts browser information from a user agent string
// This is a convenience function that wraps ParseBrowserFromUserAgent
func ParseBrowser(userAgent string) BrowserInfo {
	return ParseBrowserFromUserAgent(userAgent)
}

// ParseBrowserFromHTTPRequest extracts browser information from an HTTP request
// This is a convenience function that wraps ParseBrowserFromRequest
func ParseBrowserFromHTTPRequest(req *http.Request) BrowserInfo {
	return ParseBrowserFromRequest(req)
}

// GetBrowserInfo performs comprehensive analysis of an HTTP request
// Returns both bot detection and browser information
func GetBrowserInfo(req *http.Request) (BrowserInfo, BotDetectionResult, error) {
	// Parse browser information
	browserInfo := ParseBrowserFromRequest(req)

	// Perform bot detection
	detector := NewDetector()
	botResult, err := detector.DetectFromRequest(req)

	// Update browser info with bot detection results if needed
	if botResult.Bot && !browserInfo.IsBot {
		browserInfo.IsBot = true
		browserInfo.BotKind = botResult.BotKind
	}

	return browserInfo, botResult, err
}

// IsSupportedBrowser checks if the browser meets minimum version requirements
func IsSupportedBrowser(req *http.Request, minVersions map[BrowserName]string) bool {
	browserInfo := ParseBrowserFromRequest(req)
	return browserInfo.IsSupported(minVersions)
}

// GetBrowserFamily returns the browser family from a user agent string
func GetBrowserFamily(userAgent string) string {
	browserInfo := ParseBrowserFromUserAgent(userAgent)
	return browserInfo.GetBrowserFamily()
}

// IsMobileBrowser checks if the request comes from a mobile browser
func IsMobileBrowser(req *http.Request) bool {
	browserInfo := ParseBrowserFromRequest(req)
	return browserInfo.IsMobile()
}
