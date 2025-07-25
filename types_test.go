package gogobot

import (
	"context"
	"testing"
)

func TestState_Constants(t *testing.T) {
	// Test that state constants have expected values
	if StateSuccess != 0 {
		t.Errorf("Expected StateSuccess to be 0, got %d", StateSuccess)
	}
	if StateUndefined == StateSuccess {
		t.Error("StateUndefined should not equal StateSuccess")
	}
}

func TestBotKind_Constants(t *testing.T) {
	tests := []struct {
		kind     BotKind
		expected string
	}{
		{BotKindPhantomJS, "phantomjs"},
		{BotKindSelenium, "selenium"},
		{BotKindElectron, "electron"},
		{BotKindHeadlessChrome, "headless_chrome"},
		{BotKindPlaywright, "playwright"},
		{BotKindPuppeteer, "puppeteer"},
		{BotKindCurl, "curl"},
		{BotKindWget, "wget"},
		{BotKindBot, "bot"},
		{BotKindCrawler, "crawler"},
		{BotKindUnknown, "unknown"},
	}

	for _, test := range tests {
		if string(test.kind) != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, string(test.kind))
		}
	}
}

func TestBotDetectionResult(t *testing.T) {
	// Test bot detection result structure
	result := BotDetectionResult{
		Bot:     true,
		BotKind: BotKindSelenium,
	}

	if !result.Bot {
		t.Error("Expected Bot to be true")
	}
	if result.BotKind != BotKindSelenium {
		t.Errorf("Expected BotKind to be %s, got %s", BotKindSelenium, result.BotKind)
	}

	// Test non-bot result
	nonBotResult := BotDetectionResult{Bot: false}
	if nonBotResult.Bot {
		t.Error("Expected Bot to be false")
	}
}

func TestSuccessComponent(t *testing.T) {
	component := SuccessComponent[string]{
		State: StateSuccess,
		Value: "test-value",
	}

	if component.GetState() != StateSuccess {
		t.Errorf("Expected state %d, got %d", StateSuccess, component.GetState())
	}
	if component.GetValue() != "test-value" {
		t.Errorf("Expected value 'test-value', got '%s'", component.GetValue())
	}
	if component.GetError() != "" {
		t.Errorf("Expected empty error, got '%s'", component.GetError())
	}
}

func TestErrorComponent(t *testing.T) {
	component := ErrorComponent[string]{
		State: StateUndefined,
		Error: "test error",
	}

	if component.GetState() != StateUndefined {
		t.Errorf("Expected state %d, got %d", StateUndefined, component.GetState())
	}
	if component.GetValue() != "" {
		t.Errorf("Expected empty value, got '%s'", component.GetValue())
	}
	if component.GetError() != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", component.GetError())
	}
}

func TestBotdError(t *testing.T) {
	err := NewBotdError(StateUndefined, "test message")

	if err.State != StateUndefined {
		t.Errorf("Expected state %d, got %d", StateUndefined, err.State)
	}
	if err.Message != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", err.Message)
	}

	expectedError := "BotdError (state: 1): test message" // Updated to match actual StateUndefined value
	if err.Error() != expectedError {
		t.Errorf("Expected error string '%s', got '%s'", expectedError, err.Error())
	}
}

func TestGetResultFromContext(t *testing.T) {
	// Test with result in context
	result := &BotDetectionResult{Bot: true, BotKind: BotKindSelenium}
	ctx := context.WithValue(context.Background(), DetectionResultKey, result)

	retrieved, ok := GetResultFromContext(ctx)
	if !ok {
		t.Error("Expected to find result in context")
	}
	if retrieved == nil {
		t.Error("Expected non-nil result")
	}
	if retrieved.Bot != true || retrieved.BotKind != BotKindSelenium {
		t.Error("Retrieved result doesn't match expected values")
	}

	// Test with empty context
	emptyCtx := context.Background()
	_, ok = GetResultFromContext(emptyCtx)
	if ok {
		t.Error("Expected not to find result in empty context")
	}
}

func TestGetComponentsFromContext(t *testing.T) {
	// Test with components in context
	components := &ComponentDict{
		UserAgent: SuccessComponent[string]{
			State: StateSuccess,
			Value: "test-agent",
		},
	}
	ctx := context.WithValue(context.Background(), ComponentsKey, components)

	retrieved, ok := GetComponentsFromContext(ctx)
	if !ok {
		t.Error("Expected to find components in context")
	}
	if retrieved == nil {
		t.Error("Expected non-nil components")
	}
	if retrieved.UserAgent.GetValue() != "test-agent" {
		t.Error("Retrieved components don't match expected values")
	}

	// Test with empty context
	emptyCtx := context.Background()
	_, ok = GetComponentsFromContext(emptyCtx)
	if ok {
		t.Error("Expected not to find components in empty context")
	}
}

func TestComponentDict_Structure(t *testing.T) {
	dict := &ComponentDict{
		UserAgent:   SuccessComponent[string]{State: StateSuccess, Value: "test"},
		Headers:     SuccessComponent[map[string][]string]{State: StateSuccess, Value: make(map[string][]string)},
		HeaderCount: SuccessComponent[int]{State: StateSuccess, Value: 5},
	}

	// Test that all expected fields exist and can be accessed
	if dict.UserAgent.GetState() != StateSuccess {
		t.Error("UserAgent field not accessible")
	}
	if dict.Headers.GetState() != StateSuccess {
		t.Error("Headers field not accessible")
	}
	if dict.HeaderCount.GetState() != StateSuccess {
		t.Error("HeaderCount field not accessible")
	}
}

func TestDetectionDict_Structure(t *testing.T) {
	dict := &DetectionDict{
		UserAgent:   BotDetectionResult{Bot: true, BotKind: BotKindSelenium},
		Headers:     BotDetectionResult{Bot: false},
		HeaderCount: BotDetectionResult{Bot: true, BotKind: BotKindUnknown},
	}

	// Test that all expected fields exist and can be accessed
	if !dict.UserAgent.Bot {
		t.Error("UserAgent detection result not accessible")
	}
	if dict.Headers.Bot {
		t.Error("Headers detection result should be false")
	}
	if !dict.HeaderCount.Bot {
		t.Error("HeaderCount detection result should be true")
	}
}
