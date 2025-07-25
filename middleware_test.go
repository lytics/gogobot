package gogobot

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultMiddlewareConfig(t *testing.T) {
	config := DefaultMiddlewareConfig()

	if config.SkipFunc != nil {
		t.Error("Expected SkipFunc to be nil by default")
	}
	if config.OnBotDetected != nil {
		t.Error("Expected OnBotDetected to be nil by default")
	}
	if config.OnError != nil {
		t.Error("Expected OnError to be nil by default")
	}
	if config.BlockBots {
		t.Error("Expected BlockBots to be false by default")
	}
	if config.BlockedStatusCode != http.StatusForbidden {
		t.Errorf("Expected BlockedStatusCode to be %d, got %d", http.StatusForbidden, config.BlockedStatusCode)
	}
	if config.BlockedMessage != "Bot traffic is not allowed" {
		t.Errorf("Expected default blocked message, got %s", config.BlockedMessage)
	}
}

func TestBotDetector_Middleware(t *testing.T) {
	detector := NewDetector()
	middleware := detector.Middleware()

	if middleware == nil {
		t.Fatal("Middleware() returned nil")
	}

	// Test with normal request
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Accept-Language", "en-US")

	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "success" {
		t.Errorf("Expected 'success', got %s", w.Body.String())
	}
}

func TestBotDetector_MiddlewareWithBotBlocking(t *testing.T) {
	detector := NewDetector()

	config := MiddlewareConfig{
		BlockBots:         true,
		BlockedStatusCode: http.StatusForbidden,
		BlockedMessage:    "Bot blocked",
	}

	middleware := detector.MiddlewareWithConfig(config)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	wrappedHandler := middleware(handler)

	// Test with bot request (curl)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "curl/7.68.0")

	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
	if w.Body.String() != "Bot blocked\n" {
		t.Errorf("Expected 'Bot blocked\\n', got %s", w.Body.String())
	}
}

func TestBotDetector_MiddlewareWithSkipFunc(t *testing.T) {
	detector := NewDetector()

	config := MiddlewareConfig{
		SkipFunc: func(r *http.Request) bool {
			return r.Header.Get("X-Skip-Detection") == "true"
		},
		BlockBots: true,
	}

	middleware := detector.MiddlewareWithConfig(config)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	wrappedHandler := middleware(handler)

	// Test with bot request but skip header
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "curl/7.68.0")
	req.Header.Set("X-Skip-Detection", "true")

	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "success" {
		t.Errorf("Expected 'success', got %s", w.Body.String())
	}
}

func TestBotDetector_MiddlewareWithBotDetectedCallback(t *testing.T) {
	detector := NewDetector()

	var callbackCalled bool
	var detectedResult *BotDetectionResult

	config := MiddlewareConfig{
		OnBotDetected: func(w http.ResponseWriter, r *http.Request, result *BotDetectionResult) {
			callbackCalled = true
			detectedResult = result
			w.WriteHeader(http.StatusTeapot) // Custom status for testing
			w.Write([]byte("custom bot response"))
		},
	}

	middleware := detector.MiddlewareWithConfig(config)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when bot is detected and callback is provided")
	})

	wrappedHandler := middleware(handler)

	// Test with bot request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "curl/7.68.0")

	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if !callbackCalled {
		t.Error("Expected OnBotDetected callback to be called")
	}
	if detectedResult == nil || !detectedResult.Bot {
		t.Error("Expected bot detection result in callback")
	}
	if w.Code != http.StatusTeapot {
		t.Errorf("Expected status %d, got %d", http.StatusTeapot, w.Code)
	}
	if w.Body.String() != "custom bot response" {
		t.Errorf("Expected custom response, got %s", w.Body.String())
	}
}

func TestBotDetector_MiddlewareWithErrorCallback(t *testing.T) {
	// Create a detector that will cause an error
	detector := NewDetectorWithCustomDetectors(map[string]DetectorFunc{
		"errorDetector": func(components *ComponentDict) *BotDetectionResult {
			// This won't cause an error in detection, but we can test error handling
			return &BotDetectionResult{Bot: false}
		},
	})

	var errorCallbackCalled bool

	config := MiddlewareConfig{
		OnError: func(w http.ResponseWriter, r *http.Request, err error) {
			errorCallbackCalled = true
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error occurred"))
		},
	}

	middleware := detector.MiddlewareWithConfig(config)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	wrappedHandler := middleware(handler)

	// Test with normal request (no error expected in this case)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	// Since we don't actually cause an error, the handler should succeed
	if errorCallbackCalled {
		t.Error("Error callback should not be called for successful detection")
	}
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestBotDetector_MiddlewareContextPassing(t *testing.T) {
	detector := NewDetector()
	middleware := detector.Middleware()

	var contextResult *BotDetectionResult
	var contextComponents *ComponentDict

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		contextResult, ok = GetResultFromContext(r.Context())
		if !ok {
			t.Error("Expected to find detection result in context")
		}

		contextComponents, ok = GetComponentsFromContext(r.Context())
		if !ok {
			t.Error("Expected to find components in context")
		}

		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "curl/7.68.0")

	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	if contextResult == nil {
		t.Fatal("Expected detection result in context")
	}
	if !contextResult.Bot {
		t.Error("Expected bot to be detected")
	}
	if contextComponents == nil {
		t.Error("Expected components in context")
	}
	if contextComponents.UserAgent.GetValue() != "curl/7.68.0" {
		t.Error("Expected user agent to be preserved in context")
	}
}

func TestBotDetector_HandlerFunc(t *testing.T) {
	detector := NewDetector()

	originalHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	wrappedHandler := detector.HandlerFunc(originalHandler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	w := httptest.NewRecorder()
	wrappedHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "success" {
		t.Errorf("Expected 'success', got %s", w.Body.String())
	}
}

func TestBotDetector_HandlerFuncWithConfig(t *testing.T) {
	detector := NewDetector()

	config := MiddlewareConfig{
		BlockBots: true,
	}

	originalHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	wrappedHandler := detector.HandlerFuncWithConfig(config, originalHandler)

	// Test with bot request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "curl/7.68.0")

	w := httptest.NewRecorder()
	wrappedHandler(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}
