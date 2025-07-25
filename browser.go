package gogobot

import (
	"net/http"
	"regexp"
	"strings"
)

// ParseBrowserFromUserAgent extracts browser information from a user agent string
func ParseBrowserFromUserAgent(userAgent string) BrowserInfo {
	if userAgent == "" {
		return BrowserInfo{
			Name:    BrowserUnknown,
			Version: "",
			RawUA:   userAgent,
		}
	}

	ua := strings.TrimSpace(userAgent)
	browserInfo := BrowserInfo{
		RawUA: ua,
	}

	// First check if it's a bot
	isBot, botKind := IsBotUserAgent(ua)
	if isBot {
		browserInfo.IsBot = true
		browserInfo.BotKind = botKind
		browserInfo.Name = BrowserUnknown
		return browserInfo
	}

	// Parse browser name and version
	browserInfo.Name, browserInfo.Version = parseBrowserNameAndVersion(ua)
	return browserInfo
}

// parseBrowserNameAndVersion extracts browser name and version from user agent
func parseBrowserNameAndVersion(ua string) (BrowserName, string) {
	ua = strings.ToLower(ua)

	// Browser patterns ordered by specificity (most specific first)
	patterns := []struct {
		name    BrowserName
		pattern string
	}{
		// Microsoft Edge (must come before Chrome as it contains "chrome")
		{BrowserEdge, `edg(?:e|a|ios)?\/([0-9]+(?:\.[0-9]+)*)`},

		// Chrome-based browsers (must come before generic Chrome)
		{BrowserYandex, `yabrowser\/([0-9]+(?:\.[0-9]+)*)`},
		{BrowserVivaldi, `vivaldi\/([0-9]+(?:\.[0-9]+)*)`},
		{BrowserBrave, `brave\/([0-9]+(?:\.[0-9]+)*)`},
		{BrowserSamsung, `samsungbrowser\/([0-9]+(?:\.[0-9]+)*)`},
		{BrowserUCBrowser, `ucbrowser\/([0-9]+(?:\.[0-9]+)*)`},

		// Opera (various versions and formats)
		{BrowserOpera, `(?:opera|opr)\/([0-9]+(?:\.[0-9]+)*)`},
		{BrowserOpera, `version\/([0-9]+(?:\.[0-9]+)*).*opera`},

		// Chrome (must come after Chrome-based browsers)
		{BrowserChrome, `chrome\/([0-9]+(?:\.[0-9]+)*)`},
		{BrowserChrome, `chromium\/([0-9]+(?:\.[0-9]+)*)`},

		// Firefox
		{BrowserFirefox, `firefox\/([0-9]+(?:\.[0-9]+)*)`},
		{BrowserFirefox, `fxios\/([0-9]+(?:\.[0-9]+)*)`}, // Firefox iOS

		// Safari (must come after other webkit browsers)
		{BrowserSafari, `version\/([0-9]+(?:\.[0-9]+)*).*safari`},
		{BrowserSafari, `mobile\/[0-9a-z]+.*safari\/([0-9]+(?:\.[0-9]+)*)`}, // Mobile Safari

		// Internet Explorer
		{BrowserIE, `msie ([0-9]+(?:\.[0-9]+)*)`},
		{BrowserIE, `trident\/.*rv:([0-9]+(?:\.[0-9]+)*)`}, // IE 11+
	}

	for _, p := range patterns {
		re := regexp.MustCompile(p.pattern)
		matches := re.FindStringSubmatch(ua)
		if len(matches) >= 2 {
			version := matches[1]
			// For Safari, handle special cases
			if p.name == BrowserSafari {
				version = cleanSafariVersion(version, ua)
			}
			return p.name, version
		}
	}

	// Fallback patterns for edge cases
	fallbackPatterns := []struct {
		name    BrowserName
		pattern string
	}{
		{BrowserSafari, `safari`},
		{BrowserChrome, `chrome`},
		{BrowserFirefox, `firefox`},
		{BrowserOpera, `opera`},
		{BrowserEdge, `edge`},
	}

	for _, p := range fallbackPatterns {
		if strings.Contains(ua, string(p.name)) {
			return p.name, ""
		}
	}

	return BrowserUnknown, ""
}

// cleanSafariVersion handles Safari version parsing edge cases
func cleanSafariVersion(version, ua string) string {
	// Safari version mapping for iOS
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		// For mobile Safari, the version often corresponds to iOS version
		// This is a simplified mapping
		return version
	}

	// For desktop Safari, version is usually accurate
	return version
}

// ParseBrowserFromRequest extracts browser information from an HTTP request
func ParseBrowserFromRequest(req *http.Request) BrowserInfo {
	userAgent := req.Header.Get("User-Agent")
	return ParseBrowserFromUserAgent(userAgent)
}

// GetBrowserFamily returns the browser family (useful for grouping similar browsers)
func (bi BrowserInfo) GetBrowserFamily() string {
	switch bi.Name {
	case BrowserChrome, BrowserEdge, BrowserYandex, BrowserVivaldi, BrowserBrave, BrowserSamsung, BrowserUCBrowser:
		return "chromium"
	case BrowserFirefox:
		return "gecko"
	case BrowserSafari:
		return "webkit"
	case BrowserOpera:
		return "opera"
	case BrowserIE:
		return "trident"
	default:
		return "unknown"
	}
}

// IsMobile attempts to detect if the browser is on a mobile device
func (bi BrowserInfo) IsMobile() bool {
	ua := strings.ToLower(bi.RawUA)
	mobileIndicators := []string{
		"mobile", "android", "iphone", "ipad", "ipod",
		"blackberry", "windows phone", "palm", "symbian",
	}

	for _, indicator := range mobileIndicators {
		if strings.Contains(ua, indicator) {
			return true
		}
	}
	return false
}

// GetMajorVersion returns just the major version number
func (bi BrowserInfo) GetMajorVersion() string {
	if bi.Version == "" {
		return ""
	}

	parts := strings.Split(bi.Version, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return bi.Version
}

// IsSupported checks if the browser version meets minimum requirements
func (bi BrowserInfo) IsSupported(minVersions map[BrowserName]string) bool {
	if bi.IsBot {
		return false
	}

	minVersion, exists := minVersions[bi.Name]
	if !exists {
		return false
	}

	return compareVersions(bi.Version, minVersion) >= 0
}

// compareVersions compares two version strings
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	if v1 == v2 {
		return 0
	}

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var p1, p2 int

		if i < len(parts1) {
			p1 = parseVersionPart(parts1[i])
		}
		if i < len(parts2) {
			p2 = parseVersionPart(parts2[i])
		}

		if p1 < p2 {
			return -1
		}
		if p1 > p2 {
			return 1
		}
	}

	return 0
}

// parseVersionPart extracts numeric part from version component
func parseVersionPart(part string) int {
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(part)
	if match == "" {
		return 0
	}

	var result int
	for _, char := range match {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		}
	}
	return result
}
