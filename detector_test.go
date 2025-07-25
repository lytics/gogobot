package gogobot

import (
	"net/http"
	"net/url"
	"testing"
)

func TestNewDetector(t *testing.T) {
	detector := NewDetector()

	if detector == nil {
		t.Fatal("NewDetector() returned nil")
	}
	if detector.detectorFuncs == nil {
		t.Error("Expected detector functions to be initialized")
	}
	if len(detector.detectorFuncs) == 0 {
		t.Error("Expected some default detectors to be loaded")
	}
}

func TestNewDetectorWithCustomDetectors(t *testing.T) {
	customDetectors := map[string]DetectorFunc{
		"test": func(components *ComponentDict) *BotDetectionResult {
			return &BotDetectionResult{Bot: true, BotKind: BotKindUnknown}
		},
	}

	detector := NewDetectorWithCustomDetectors(customDetectors)

	if detector == nil {
		t.Fatal("NewDetectorWithCustomDetectors() returned nil")
	}
	if detector.detectorFuncs == nil {
		t.Error("Expected detector functions to be initialized")
	}

	// Should have default detectors plus custom ones
	if len(detector.detectorFuncs) < len(customDetectors) {
		t.Error("Expected custom detectors to be added")
	}

	// Check that custom detector exists
	if _, exists := detector.detectorFuncs["test"]; !exists {
		t.Error("Expected custom detector 'test' to be present")
	}
}

func TestBotDetector_Collect(t *testing.T) {
	detector := NewDetector()

	// Create test request
	req := createTestRequest("GET", "/test", map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		"Accept":     "text/html,application/xhtml+xml",
		"Connection": "keep-alive",
	})

	components, err := detector.Collect(req)

	if err != nil {
		t.Fatalf("Collect() returned error: %v", err)
	}
	if components == nil {
		t.Fatal("Collect() returned nil components")
	}

	// Test that components were collected
	if components.UserAgent.GetState() != StateSuccess {
		t.Error("Expected UserAgent to be collected successfully")
	}
	if components.UserAgent.GetValue() != "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36" {
		t.Error("UserAgent value not collected correctly")
	}
	if components.RequestMethod.GetValue() != "GET" {
		t.Error("Request method not collected correctly")
	}
	if components.RequestPath.GetValue() != "/test" {
		t.Error("Request path not collected correctly")
	}
}

func TestBotDetector_Detect(t *testing.T) {
	detector := NewDetector()

	// Test panic when Detect() called before Collect()
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected Detect() to panic when called before Collect()")
		}
	}()
	detector.Detect()
}

func TestBotDetector_DetectAfterCollect(t *testing.T) {
	detector := NewDetector()

	// Create bot request (curl)
	req := createTestRequest("GET", "/", map[string]string{
		"User-Agent": "curl/7.68.0",
		"Accept":     "*/*",
	})

	_, err := detector.Collect(req)
	if err != nil {
		t.Fatalf("Collect() returned error: %v", err)
	}

	result := detector.Detect()

	if !result.Bot {
		t.Error("Expected curl request to be detected as bot")
	}
	if result.BotKind != BotKindCurl {
		t.Errorf("Expected bot kind %s, got %s", BotKindCurl, result.BotKind)
	}
}

func TestBotDetector_DetectFromRequest(t *testing.T) {
	tests := []struct {
		name         string
		userAgent    string
		headers      map[string]string
		expectedBot  bool
		expectedKind BotKind
	}{
		{
			name:      "Normal browser",
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
			name:         "Curl bot",
			userAgent:    "curl/7.68.0",
			headers:      map[string]string{"Accept": "*/*"},
			expectedBot:  true,
			expectedKind: BotKindCurl,
		},
		{
			name:         "Python requests",
			userAgent:    "python-requests/2.25.1",
			headers:      map[string]string{"Accept": "*/*"},
			expectedBot:  true,
			expectedKind: BotKindUnknown,
		},
		{
			name:         "PhantomJS",
			userAgent:    "Mozilla/5.0 (Unknown; Linux x86_64) AppleWebKit/534.34 (KHTML, like Gecko) PhantomJS/2.1.1 Safari/534.34",
			expectedBot:  true,
			expectedKind: BotKindPhantomJS,
		},
		{
			name:         "Selenium",
			userAgent:    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 selenium/3.141.0",
			expectedBot:  true,
			expectedKind: BotKindSelenium,
		},
		{
			name:         "Headless Chrome",
			userAgent:    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/91.0.4472.124 Safari/537.36",
			expectedBot:  true,
			expectedKind: BotKindHeadlessChrome,
		},
	}

	detector := NewDetector()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := createTestRequest("GET", "/", map[string]string{"User-Agent": test.userAgent})

			// Add additional headers if provided
			for k, v := range test.headers {
				req.Header.Set(k, v)
			}

			result, err := detector.DetectFromRequest(req)

			if err != nil {
				t.Fatalf("DetectFromRequest() returned error: %v", err)
			}

			if result.Bot != test.expectedBot {
				t.Errorf("Expected bot=%t, got bot=%t", test.expectedBot, result.Bot)
			}

			if test.expectedBot && result.BotKind != test.expectedKind {
				t.Errorf("Expected bot kind %s, got %s", test.expectedKind, result.BotKind)
			}
		})
	}
}

func TestBotDetector_GetComponents(t *testing.T) {
	detector := NewDetector()

	// Before collecting
	if detector.GetComponents() != nil {
		t.Error("Expected GetComponents() to return nil before Collect()")
	}

	// After collecting
	req := createTestRequest("GET", "/", map[string]string{"User-Agent": "test"})
	detector.Collect(req)

	components := detector.GetComponents()
	if components == nil {
		t.Error("Expected GetComponents() to return non-nil after Collect()")
	}
}

func TestBotDetector_GetDetections(t *testing.T) {
	detector := NewDetector()

	// Before detecting
	if detector.GetDetections() != nil {
		t.Error("Expected GetDetections() to return nil before Detect()")
	}

	// After detecting
	req := createTestRequest("GET", "/", map[string]string{"User-Agent": "curl/7.68.0"})
	detector.Collect(req)
	detector.Detect()

	detections := detector.GetDetections()
	if detections == nil {
		t.Error("Expected GetDetections() to return non-nil after Detect()")
	}
}

func TestBotDetector_AddRemoveDetector(t *testing.T) {
	detector := NewDetector()

	// Test adding detector
	testDetector := func(components *ComponentDict) *BotDetectionResult {
		return &BotDetectionResult{Bot: true, BotKind: BotKindUnknown}
	}

	initialCount := len(detector.GetDetectorNames())
	detector.AddDetector("test", testDetector)

	if len(detector.GetDetectorNames()) != initialCount+1 {
		t.Error("Expected detector count to increase after AddDetector()")
	}

	// Test that detector was added
	names := detector.GetDetectorNames()
	found := false
	for _, name := range names {
		if name == "test" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find 'test' detector in detector names")
	}

	// Test removing detector
	detector.RemoveDetector("test")
	if len(detector.GetDetectorNames()) != initialCount {
		t.Error("Expected detector count to decrease after RemoveDetector()")
	}

	// Test that detector was removed
	names = detector.GetDetectorNames()
	for _, name := range names {
		if name == "test" {
			t.Error("Expected 'test' detector to be removed")
		}
	}
}

func TestDetectUserAgent(t *testing.T) {
	tests := []struct {
		name         string
		userAgent    string
		expectedBot  bool
		expectedKind BotKind
	}{
		{"Empty UA", "", false, ""},
		{"Normal Chrome", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36", false, ""},
		{"Curl", "curl/7.68.0", true, BotKindCurl},
		{"Wget", "Wget/1.20.3", true, BotKindWget},
		{"Python requests", "python-requests/2.25.1", true, BotKindUnknown},
		{"PhantomJS", "PhantomJS/2.1.1", true, BotKindPhantomJS},
		{"Selenium", "selenium webdriver", true, BotKindSelenium},
		{"Headless", "HeadlessChrome/91.0", true, BotKindHeadlessChrome},
		{"Googlebot", "Googlebot/2.1", true, BotKindCrawler},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			components := &ComponentDict{}
			if test.userAgent == "" {
				components.UserAgent = ErrorComponent[string]{
					State: StateUndefined,
					Error: "missing",
				}
			} else {
				components.UserAgent = SuccessComponent[string]{
					State: StateSuccess,
					Value: test.userAgent,
				}
			}

			result := detectUserAgent(components)

			if result.Bot != test.expectedBot {
				t.Errorf("Expected bot=%t, got bot=%t", test.expectedBot, result.Bot)
			}

			if test.expectedBot && result.BotKind != test.expectedKind {
				t.Errorf("Expected bot kind %s, got %s", test.expectedKind, result.BotKind)
			}
		})
	}
}

func TestDetectHeaders(t *testing.T) {
	tests := []struct {
		name        string
		headers     map[string][]string
		expectedBot bool
	}{
		{
			name:        "Normal headers",
			headers:     map[string][]string{"User-Agent": {"Mozilla/5.0"}, "Accept": {"text/html"}},
			expectedBot: false,
		},
		{
			name:        "Automation header",
			headers:     map[string][]string{"X-Requested-With": {"XMLHttpRequest"}},
			expectedBot: true,
		},
		{
			name:        "Chrome proxy header",
			headers:     map[string][]string{"Chrome-Proxy": {"frfr"}},
			expectedBot: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			components := &ComponentDict{
				Headers: SuccessComponent[map[string][]string]{
					State: StateSuccess,
					Value: test.headers,
				},
			}

			result := detectHeaders(components)

			if result.Bot != test.expectedBot {
				t.Errorf("Expected bot=%t, got bot=%t", test.expectedBot, result.Bot)
			}
		})
	}
}

func TestDetectHeaderCount(t *testing.T) {
	tests := []struct {
		name        string
		count       int
		expectedBot bool
	}{
		{"Too few headers", 2, true},
		{"Normal count", 8, false},
		{"Too many headers", 35, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			components := &ComponentDict{
				HeaderCount: SuccessComponent[int]{
					State: StateSuccess,
					Value: test.count,
				},
			}

			result := detectHeaderCount(components)

			if result.Bot != test.expectedBot {
				t.Errorf("Expected bot=%t, got bot=%t", test.expectedBot, result.Bot)
			}
		})
	}
}

// Helper function to create test HTTP requests
func createTestRequest(method, path string, headers map[string]string) *http.Request {
	req := &http.Request{
		Method: method,
		URL:    &url.URL{Path: path},
		Header: make(http.Header),
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req
}
