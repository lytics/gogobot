package gogobot

import (
	"context"
	"fmt"
	"net/http"
)

// State represents the source collection state
type State int

const (
	StateSuccess State = iota
	StateUndefined
	StateNotFunction
	StateUnexpectedBehaviour
	StateNull
)

// BotKind represents different types of detected bots
type BotKind string

const (
	BotKindAwesomium      BotKind = "awesomium"
	BotKindCef            BotKind = "cef"
	BotKindCefSharp       BotKind = "cefsharp"
	BotKindCoachJS        BotKind = "coachjs"
	BotKindElectron       BotKind = "electron"
	BotKindFMiner         BotKind = "fminer"
	BotKindGeb            BotKind = "geb"
	BotKindNightmareJS    BotKind = "nightmarejs"
	BotKindPhantomas      BotKind = "phantomas"
	BotKindPhantomJS      BotKind = "phantomjs"
	BotKindRhino          BotKind = "rhino"
	BotKindSelenium       BotKind = "selenium"
	BotKindSequentum      BotKind = "sequentum"
	BotKindSlimerJS       BotKind = "slimerjs"
	BotKindWebDriverIO    BotKind = "webdriverio"
	BotKindWebDriver      BotKind = "webdriver"
	BotKindHeadlessChrome BotKind = "headless_chrome"
	BotKindPlaywright     BotKind = "playwright"
	BotKindPuppeteer      BotKind = "puppeteer"
	BotKindCurl           BotKind = "curl"
	BotKindWget           BotKind = "wget"
	BotKindBot            BotKind = "bot"
	BotKindCrawler        BotKind = "crawler"
	BotKindSpider         BotKind = "spider"
	BotKindUnknown        BotKind = "unknown"
)

// BotDetectionResult represents the result of bot detection
type BotDetectionResult struct {
	Bot     bool    `json:"bot"`
	BotKind BotKind `json:"botKind,omitempty"`
}

// Component represents a data component with state and value
type Component[T any] interface {
	GetState() State
	GetValue() T
	GetError() string
}

// SuccessComponent represents a successful component
type SuccessComponent[T any] struct {
	State State
	Value T
}

func (c SuccessComponent[T]) GetState() State  { return c.State }
func (c SuccessComponent[T]) GetValue() T      { return c.Value }
func (c SuccessComponent[T]) GetError() string { return "" }

// ErrorComponent represents a failed component
type ErrorComponent[T any] struct {
	State State
	Error string
}

func (c ErrorComponent[T]) GetState() State  { return c.State }
func (c ErrorComponent[T]) GetValue() T      { var zero T; return zero }
func (c ErrorComponent[T]) GetError() string { return c.Error }

// ComponentDict holds all collected components
type ComponentDict struct {
	UserAgent            Component[string]
	XForwardedFor        Component[string]
	XRealIP              Component[string]
	AcceptLanguage       Component[string]
	AcceptEncoding       Component[string]
	AcceptCharset        Component[string]
	Accept               Component[string]
	Connection           Component[string]
	CacheControl         Component[string]
	UpgradeInsecure      Component[bool]
	DNT                  Component[string]
	Headers              Component[map[string][]string]
	ContentLength        Component[int64]
	RequestMethod        Component[string]
	RequestPath          Component[string]
	RequestQuery         Component[string]
	RemoteAddr           Component[string]
	HeaderOrder          Component[[]string]
	HeaderCount          Component[int]
	MissingCommonHeaders Component[[]string]
}

// DetectionDict holds detection results for each detector
type DetectionDict struct {
	UserAgent      BotDetectionResult
	Headers        BotDetectionResult
	HeaderOrder    BotDetectionResult
	HeaderCount    BotDetectionResult
	MissingHeaders BotDetectionResult
	AcceptHeaders  BotDetectionResult
	Connection     BotDetectionResult
	ContentLength  BotDetectionResult
}

// BotDetectorInterface defines the interface for bot detectors
type BotDetectorInterface interface {
	Detect() BotDetectionResult
	Collect(*http.Request) (*ComponentDict, error)
	GetComponents() *ComponentDict
	GetDetections() *DetectionDict
}

// BotdError represents errors during bot detection
type BotdError struct {
	State   State
	Message string
}

func (e BotdError) Error() string {
	return fmt.Sprintf("BotdError (state: %d): %s", e.State, e.Message)
}

// NewBotdError creates a new BotdError
func NewBotdError(state State, message string) *BotdError {
	return &BotdError{
		State:   state,
		Message: message,
	}
}

// DetectorFunc is a function that performs bot detection on components
type DetectorFunc func(*ComponentDict) *BotDetectionResult

// SourceFunc is a function that collects data from an HTTP request
type SourceFunc[T any] func(*http.Request) Component[T]

// Context keys for storing detection results
type contextKey string

const (
	DetectionResultKey contextKey = "gogobot_detection_result"
	ComponentsKey      contextKey = "gogobot_components"
)

// GetResultFromContext retrieves the detection result from request context
func GetResultFromContext(ctx context.Context) (*BotDetectionResult, bool) {
	result, ok := ctx.Value(DetectionResultKey).(*BotDetectionResult)
	return result, ok
}

// GetComponentsFromContext retrieves the components from request context
func GetComponentsFromContext(ctx context.Context) (*ComponentDict, bool) {
	components, ok := ctx.Value(ComponentsKey).(*ComponentDict)
	return components, ok
}
