package gogobot

import (
	"net/http"
	"regexp"
	"strings"
)

// BotDetector is the main struct for bot detection
type BotDetector struct {
	components    *ComponentDict
	detections    *DetectionDict
	detectorFuncs map[string]DetectorFunc
}

// NewDetector creates a new BotDetector instance
func NewDetector() *BotDetector {
	return &BotDetector{
		detectorFuncs: getDefaultDetectors(),
	}
}

// NewDetectorWithCustomDetectors creates a new BotDetector with custom detectors
func NewDetectorWithCustomDetectors(customDetectors map[string]DetectorFunc) *BotDetector {
	allDetectors := getDefaultDetectors()

	// Merge custom detectors with default ones
	for name, detector := range customDetectors {
		allDetectors[name] = detector
	}

	return &BotDetector{
		detectorFuncs: allDetectors,
	}
}

// Collect gathers data from the HTTP request
func (d *BotDetector) Collect(req *http.Request) (*ComponentDict, error) {
	d.components = collectAllSources(req)
	return d.components, nil
}

// Detect performs bot detection on the collected components
func (d *BotDetector) Detect() BotDetectionResult {
	if d.components == nil {
		panic("BotDetector.Detect() called before Collect()")
	}

	detections := &DetectionDict{}
	finalResult := BotDetectionResult{Bot: false}
	var bestResult BotDetectionResult

	// Run all detectors
	for name, detectorFunc := range d.detectorFuncs {
		result := detectorFunc(d.components)
		if result == nil {
			result = &BotDetectionResult{Bot: false}
		}

		// Store individual detection results
		switch name {
		case "userAgent":
			detections.UserAgent = *result
		case "headers":
			detections.Headers = *result
		case "headerOrder":
			detections.HeaderOrder = *result
		case "headerCount":
			detections.HeaderCount = *result
		case "missingHeaders":
			detections.MissingHeaders = *result
		case "acceptHeaders":
			detections.AcceptHeaders = *result
		case "connection":
			detections.Connection = *result
		case "contentLength":
			detections.ContentLength = *result
		}

		// If any detector finds a bot, consider it for final result
		if result.Bot {
			// Prioritize specific bot kinds over unknown
			if !bestResult.Bot ||
				(result.BotKind != BotKindUnknown && bestResult.BotKind == BotKindUnknown) ||
				(name == "userAgent" && result.BotKind != BotKindUnknown) { // Prioritize user agent detection for specific types
				bestResult = *result
			}
			finalResult.Bot = true // At least one detector found a bot
		}
	}

	// Use the best (most specific) result
	if bestResult.Bot {
		finalResult = bestResult
	}

	d.detections = detections
	return finalResult
}

// DetectFromRequest is a convenience method that collects and detects in one call
func (d *BotDetector) DetectFromRequest(req *http.Request) (BotDetectionResult, error) {
	_, err := d.Collect(req)
	if err != nil {
		return BotDetectionResult{Bot: false}, err
	}

	result := d.Detect()
	return result, nil
}

// GetComponents returns the collected components
func (d *BotDetector) GetComponents() *ComponentDict {
	return d.components
}

// GetDetections returns the detection results for each detector
func (d *BotDetector) GetDetections() *DetectionDict {
	return d.detections
}

// AddDetector adds a custom detector to the detector
func (d *BotDetector) AddDetector(name string, detector DetectorFunc) {
	if d.detectorFuncs == nil {
		d.detectorFuncs = make(map[string]DetectorFunc)
	}
	d.detectorFuncs[name] = detector
}

// RemoveDetector removes a detector by name
func (d *BotDetector) RemoveDetector(name string) {
	if d.detectorFuncs != nil {
		delete(d.detectorFuncs, name)
	}
}

// GetDetectorNames returns the names of all active detectors
func (d *BotDetector) GetDetectorNames() []string {
	names := make([]string, 0, len(d.detectorFuncs))
	for name := range d.detectorFuncs {
		names = append(names, name)
	}
	return names
}

// collectAllSources collects all data sources from the HTTP request
func collectAllSources(req *http.Request) *ComponentDict {
	return &ComponentDict{
		UserAgent:            getUserAgent(req),
		XForwardedFor:        getXForwardedFor(req),
		XRealIP:              getXRealIP(req),
		AcceptLanguage:       getAcceptLanguage(req),
		AcceptEncoding:       getAcceptEncoding(req),
		AcceptCharset:        getAcceptCharset(req),
		Accept:               getAccept(req),
		Connection:           getConnection(req),
		CacheControl:         getCacheControl(req),
		UpgradeInsecure:      getUpgradeInsecureRequests(req),
		DNT:                  getDNT(req),
		Headers:              getHeaders(req),
		ContentLength:        getContentLength(req),
		RequestMethod:        getRequestMethod(req),
		RequestPath:          getRequestPath(req),
		RequestQuery:         getRequestQuery(req),
		RemoteAddr:           getRemoteAddr(req),
		HeaderOrder:          getHeaderOrder(req),
		HeaderCount:          getHeaderCount(req),
		MissingCommonHeaders: getMissingCommonHeaders(req),
	}
}

// Source collection functions
func getUserAgent(req *http.Request) Component[string] {
	userAgent := req.Header.Get("User-Agent")
	if userAgent == "" {
		return ErrorComponent[string]{
			State: StateUndefined,
			Error: "User-Agent header is missing",
		}
	}
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: userAgent,
	}
}

func getXForwardedFor(req *http.Request) Component[string] {
	xff := req.Header.Get("X-Forwarded-For")
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: xff,
	}
}

func getXRealIP(req *http.Request) Component[string] {
	xrip := req.Header.Get("X-Real-IP")
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: xrip,
	}
}

func getAcceptLanguage(req *http.Request) Component[string] {
	acceptLang := req.Header.Get("Accept-Language")
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: acceptLang,
	}
}

func getAcceptEncoding(req *http.Request) Component[string] {
	acceptEnc := req.Header.Get("Accept-Encoding")
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: acceptEnc,
	}
}

func getAcceptCharset(req *http.Request) Component[string] {
	acceptCharset := req.Header.Get("Accept-Charset")
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: acceptCharset,
	}
}

func getAccept(req *http.Request) Component[string] {
	accept := req.Header.Get("Accept")
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: accept,
	}
}

func getConnection(req *http.Request) Component[string] {
	connection := req.Header.Get("Connection")
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: connection,
	}
}

func getCacheControl(req *http.Request) Component[string] {
	cacheControl := req.Header.Get("Cache-Control")
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: cacheControl,
	}
}

func getUpgradeInsecureRequests(req *http.Request) Component[bool] {
	upgrade := req.Header.Get("Upgrade-Insecure-Requests")
	value := upgrade == "1"
	return SuccessComponent[bool]{
		State: StateSuccess,
		Value: value,
	}
}

func getDNT(req *http.Request) Component[string] {
	dnt := req.Header.Get("DNT")
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: dnt,
	}
}

func getHeaders(req *http.Request) Component[map[string][]string] {
	headers := make(map[string][]string)
	for k, v := range req.Header {
		headers[k] = v
	}
	return SuccessComponent[map[string][]string]{
		State: StateSuccess,
		Value: headers,
	}
}

func getContentLength(req *http.Request) Component[int64] {
	return SuccessComponent[int64]{
		State: StateSuccess,
		Value: req.ContentLength,
	}
}

func getRequestMethod(req *http.Request) Component[string] {
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: req.Method,
	}
}

func getRequestPath(req *http.Request) Component[string] {
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: req.URL.Path,
	}
}

func getRequestQuery(req *http.Request) Component[string] {
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: req.URL.RawQuery,
	}
}

func getRemoteAddr(req *http.Request) Component[string] {
	return SuccessComponent[string]{
		State: StateSuccess,
		Value: req.RemoteAddr,
	}
}

func getHeaderOrder(req *http.Request) Component[[]string] {
	var order []string
	for k := range req.Header {
		order = append(order, k)
	}
	return SuccessComponent[[]string]{
		State: StateSuccess,
		Value: order,
	}
}

func getHeaderCount(req *http.Request) Component[int] {
	count := len(req.Header)
	return SuccessComponent[int]{
		State: StateSuccess,
		Value: count,
	}
}

func getMissingCommonHeaders(req *http.Request) Component[[]string] {
	commonHeaders := []string{
		"Accept",
		"Accept-Language",
		"Accept-Encoding",
		"Connection",
		"User-Agent",
	}

	var missing []string
	for _, header := range commonHeaders {
		if req.Header.Get(header) == "" {
			missing = append(missing, header)
		}
	}

	return SuccessComponent[[]string]{
		State: StateSuccess,
		Value: missing,
	}
}

// Detector functions
func detectUserAgent(components *ComponentDict) *BotDetectionResult {
	if components.UserAgent.GetState() != StateSuccess {
		return &BotDetectionResult{Bot: false}
	}

	userAgent := strings.ToLower(components.UserAgent.GetValue())

	// Check for specific bot types in order of specificity (most specific first)
	specificBots := []struct {
		kind     BotKind
		patterns []string
	}{
		// AI Agents (check first as they're highly specific)
		{BotKindGPTBot, []string{"gptbot", "gpt-bot"}},
		{BotKindChatGPT, []string{"chatgpt-user", "chatgpt", "openai-chatgpt"}},
		{BotKindOpenAI, []string{"openai", "openai-bot", "openai-crawler"}},
		{BotKindClaude, []string{"claude-web", "claude", "anthropic"}},
		{BotKindAIAgent, []string{"ai-agent", "aiagent", "ai_agent", "artificial intelligence", "language model", "llm", "gpt-", "claude-", "bard", "gemini-pro"}},

		// Automation Tools
		{BotKindPhantomJS, []string{"phantomjs"}},
		{BotKindSelenium, []string{"selenium", "webdriver"}},
		{BotKindElectron, []string{"electron"}},
		{BotKindHeadlessChrome, []string{"headlesschrome", "headless"}},
		{BotKindPlaywright, []string{"playwright"}},
		{BotKindPuppeteer, []string{"puppeteer"}},

		// Command Line Tools
		{BotKindCurl, []string{"curl/"}},
		{BotKindWget, []string{"wget/"}},

		// Search Engine Crawlers
		{BotKindCrawler, []string{"googlebot", "bingbot", "slurp", "duckduckbot", "baiduspider"}}, // Check crawlers before generic "bot"

		// Generic Bots (last to avoid false positives)
		{BotKindBot, []string{"bot", "crawler", "spider", "scraper"}},
	}

	for _, botType := range specificBots {
		for _, pattern := range botType.patterns {
			if strings.Contains(userAgent, pattern) {
				return &BotDetectionResult{
					Bot:     true,
					BotKind: botType.kind,
				}
			}
		}
	}

	// Check for suspicious user agent patterns
	suspiciousPatterns := []string{
		"^$",
		"^\\s*$",
		"mozilla/5.0$",
		"mozilla/4.0$",
		"python",
		"java",
		"go-http-client",
		"libwww-perl",
		"httpclient",
		"okhttp",
		"requests",
		"urllib",
	}

	for _, pattern := range suspiciousPatterns {
		if matched, _ := regexp.MatchString(pattern, userAgent); matched {
			return &BotDetectionResult{
				Bot:     true,
				BotKind: BotKindUnknown,
			}
		}
	}

	return &BotDetectionResult{Bot: false}
}

func detectHeaders(components *ComponentDict) *BotDetectionResult {
	if components.Headers.GetState() != StateSuccess {
		return &BotDetectionResult{Bot: false}
	}

	headers := components.Headers.GetValue()

	// Check for automation-specific headers
	automationHeaders := []string{
		"X-Requested-With",
		"X-DevTools-Emulate-Network-Conditions-Client-Id",
		"Chrome-Proxy",
		"Purpose",
	}

	for _, header := range automationHeaders {
		if _, exists := headers[header]; exists {
			return &BotDetectionResult{
				Bot:     true,
				BotKind: BotKindUnknown,
			}
		}
	}

	return &BotDetectionResult{Bot: false}
}

func detectHeaderCount(components *ComponentDict) *BotDetectionResult {
	if components.HeaderCount.GetState() != StateSuccess {
		return &BotDetectionResult{Bot: false}
	}

	count := components.HeaderCount.GetValue()

	if count < 3 { // Reduced from 4 to be less aggressive
		return &BotDetectionResult{
			Bot:     true,
			BotKind: BotKindUnknown,
		}
	}

	if count > 30 {
		return &BotDetectionResult{
			Bot:     true,
			BotKind: BotKindUnknown,
		}
	}

	return &BotDetectionResult{Bot: false}
}

func detectMissingHeaders(components *ComponentDict) *BotDetectionResult {
	if components.MissingCommonHeaders.GetState() != StateSuccess {
		return &BotDetectionResult{Bot: false}
	}

	missing := components.MissingCommonHeaders.GetValue()

	// Missing User-Agent is highly suspicious
	for _, header := range missing {
		if header == "User-Agent" {
			return &BotDetectionResult{
				Bot:     true,
				BotKind: BotKindUnknown,
			}
		}
	}

	// Only flag as bot if missing many headers (increased threshold)
	if len(missing) >= 4 { // Increased from 3 to be less aggressive
		return &BotDetectionResult{
			Bot:     true,
			BotKind: BotKindUnknown,
		}
	}

	return &BotDetectionResult{Bot: false}
}

func detectAcceptHeaders(components *ComponentDict) *BotDetectionResult {
	accept := components.Accept.GetValue()
	acceptLang := components.AcceptLanguage.GetValue()
	acceptEnc := components.AcceptEncoding.GetValue()

	// Only flag if ALL accept headers are missing (more conservative)
	if accept == "" && acceptLang == "" && acceptEnc == "" {
		return &BotDetectionResult{
			Bot:     true,
			BotKind: BotKindUnknown,
		}
	}

	// This pattern is too strict for modern browsers
	// Commenting out as it causes false positives
	// if accept == "*/*" && acceptLang == "" {
	// 	return &BotDetectionResult{
	// 		Bot:     true,
	// 		BotKind: BotKindUnknown,
	// 	}
	// }

	return &BotDetectionResult{Bot: false}
}

func detectConnection(components *ComponentDict) *BotDetectionResult {
	if components.Connection.GetState() != StateSuccess {
		return &BotDetectionResult{Bot: false}
	}

	connection := strings.ToLower(components.Connection.GetValue())
	suspiciousConnections := []string{"upgrade", "te"}

	for _, suspicious := range suspiciousConnections {
		if strings.Contains(connection, suspicious) {
			return &BotDetectionResult{
				Bot:     true,
				BotKind: BotKindUnknown,
			}
		}
	}

	return &BotDetectionResult{Bot: false}
}

func detectContentLength(components *ComponentDict) *BotDetectionResult {
	if components.ContentLength.GetState() != StateSuccess {
		return &BotDetectionResult{Bot: false}
	}

	contentLength := components.ContentLength.GetValue()
	method := components.RequestMethod.GetValue()

	if method == "GET" && contentLength > 0 {
		return &BotDetectionResult{
			Bot:     true,
			BotKind: BotKindUnknown,
		}
	}

	return &BotDetectionResult{Bot: false}
}

func detectHeaderOrder(components *ComponentDict) *BotDetectionResult {
	if components.HeaderOrder.GetState() != StateSuccess {
		return &BotDetectionResult{Bot: false}
	}

	order := components.HeaderOrder.GetValue()

	// Very conservative - only flag if extremely few headers
	if len(order) < 2 { // Reduced from 3
		return &BotDetectionResult{
			Bot:     true,
			BotKind: BotKindUnknown,
		}
	}

	return &BotDetectionResult{Bot: false}
}

// getDefaultDetectors returns the default set of detectors
func getDefaultDetectors() map[string]DetectorFunc {
	return map[string]DetectorFunc{
		"userAgent":      detectUserAgent,
		"headers":        detectHeaders,
		"headerOrder":    detectHeaderOrder,
		"headerCount":    detectHeaderCount,
		"missingHeaders": detectMissingHeaders,
		"acceptHeaders":  detectAcceptHeaders,
		"connection":     detectConnection,
		"contentLength":  detectContentLength,
	}
}
