package gogobot

import (
	"net/http"
	"testing"
)

func TestParseBrowserFromUserAgent(t *testing.T) {
	tests := []struct {
		name            string
		userAgent       string
		expectedBrowser BrowserName
		expectedVersion string
		expectedIsBot   bool
		expectedFamily  string
		expectedMobile  bool
	}{
		{
			name:            "Chrome Desktop",
			userAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			expectedBrowser: BrowserChrome,
			expectedVersion: "120.0.0.0",
			expectedIsBot:   false,
			expectedFamily:  "chromium",
			expectedMobile:  false,
		},
		{
			name:            "Firefox Desktop",
			userAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
			expectedBrowser: BrowserFirefox,
			expectedVersion: "121.0",
			expectedIsBot:   false,
			expectedFamily:  "gecko",
			expectedMobile:  false,
		},
		{
			name:            "Safari Desktop",
			userAgent:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
			expectedBrowser: BrowserSafari,
			expectedVersion: "17.1",
			expectedIsBot:   false,
			expectedFamily:  "webkit",
			expectedMobile:  false,
		},
		{
			name:            "Edge Desktop",
			userAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
			expectedBrowser: BrowserEdge,
			expectedVersion: "120.0.0.0",
			expectedIsBot:   false,
			expectedFamily:  "chromium",
			expectedMobile:  false,
		},
		{
			name:            "Opera Desktop",
			userAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/106.0.0.0",
			expectedBrowser: BrowserOpera,
			expectedVersion: "106.0.0.0",
			expectedIsBot:   false,
			expectedFamily:  "opera",
			expectedMobile:  false,
		},
		{
			name:            "Safari Mobile (iPhone)",
			userAgent:       "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1",
			expectedBrowser: BrowserSafari,
			expectedVersion: "17.1",
			expectedIsBot:   false,
			expectedFamily:  "webkit",
			expectedMobile:  true,
		},
		{
			name:            "Chrome Mobile (Android)",
			userAgent:       "Mozilla/5.0 (Linux; Android 14; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
			expectedBrowser: BrowserChrome,
			expectedVersion: "120.0.0.0",
			expectedIsBot:   false,
			expectedFamily:  "chromium",
			expectedMobile:  true,
		},
		{
			name:            "Samsung Internet",
			userAgent:       "Mozilla/5.0 (Linux; Android 14; SM-G998B) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/23.0 Chrome/115.0.0.0 Mobile Safari/537.36",
			expectedBrowser: BrowserSamsung,
			expectedVersion: "23.0",
			expectedIsBot:   false,
			expectedFamily:  "chromium",
			expectedMobile:  true,
		},
		{
			name:            "Internet Explorer 11",
			userAgent:       "Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
			expectedBrowser: BrowserIE,
			expectedVersion: "11.0",
			expectedIsBot:   false,
			expectedFamily:  "trident",
			expectedMobile:  false,
		},
		{
			name:            "Vivaldi Browser",
			userAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Vivaldi/6.5.3206.39",
			expectedBrowser: BrowserVivaldi,
			expectedVersion: "6.5.3206.39",
			expectedIsBot:   false,
			expectedFamily:  "chromium",
			expectedMobile:  false,
		},
		{
			name:            "Yandex Browser",
			userAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 YaBrowser/24.1.0.0 Safari/537.36",
			expectedBrowser: BrowserYandex,
			expectedVersion: "24.1.0.0",
			expectedIsBot:   false,
			expectedFamily:  "chromium",
			expectedMobile:  false,
		},
		{
			name:            "Bot - Googlebot",
			userAgent:       "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			expectedBrowser: BrowserUnknown,
			expectedVersion: "",
			expectedIsBot:   true,
			expectedFamily:  "unknown",
			expectedMobile:  false,
		},
		{
			name:            "Bot - PhantomJS",
			userAgent:       "Mozilla/5.0 (Unknown; Linux x86_64) AppleWebKit/538.1 (KHTML, like Gecko) PhantomJS/2.1.1 Safari/538.1",
			expectedBrowser: BrowserUnknown,
			expectedVersion: "",
			expectedIsBot:   true,
			expectedFamily:  "unknown",
			expectedMobile:  false,
		},
		{
			name:            "Empty User Agent",
			userAgent:       "",
			expectedBrowser: BrowserUnknown,
			expectedVersion: "",
			expectedIsBot:   false,
			expectedFamily:  "unknown",
			expectedMobile:  false,
		},
		{
			name:            "Unknown Browser",
			userAgent:       "SomeCustomBrowser/1.0",
			expectedBrowser: BrowserUnknown,
			expectedVersion: "",
			expectedIsBot:   false,
			expectedFamily:  "unknown",
			expectedMobile:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseBrowserFromUserAgent(tt.userAgent)

			if result.Name != tt.expectedBrowser {
				t.Errorf("Expected browser %s, got %s", tt.expectedBrowser, result.Name)
			}

			if result.Version != tt.expectedVersion {
				t.Errorf("Expected version %s, got %s", tt.expectedVersion, result.Version)
			}

			if result.IsBot != tt.expectedIsBot {
				t.Errorf("Expected IsBot %v, got %v", tt.expectedIsBot, result.IsBot)
			}

			if result.GetBrowserFamily() != tt.expectedFamily {
				t.Errorf("Expected family %s, got %s", tt.expectedFamily, result.GetBrowserFamily())
			}

			if result.IsMobile() != tt.expectedMobile {
				t.Errorf("Expected IsMobile %v, got %v", tt.expectedMobile, result.IsMobile())
			}

			if result.RawUA != tt.userAgent {
				t.Errorf("Expected RawUA to be preserved")
			}
		})
	}
}

func TestParseBrowserFromRequest(t *testing.T) {
	tests := []struct {
		name            string
		userAgent       string
		expectedBrowser BrowserName
		expectedVersion string
	}{
		{
			name:            "Chrome from HTTP request",
			userAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			expectedBrowser: BrowserChrome,
			expectedVersion: "120.0.0.0",
		},
		{
			name:            "No User-Agent header",
			userAgent:       "",
			expectedBrowser: BrowserUnknown,
			expectedVersion: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.userAgent != "" {
				req.Header.Set("User-Agent", tt.userAgent)
			}

			result := ParseBrowserFromRequest(req)

			if result.Name != tt.expectedBrowser {
				t.Errorf("Expected browser %s, got %s", tt.expectedBrowser, result.Name)
			}

			if result.Version != tt.expectedVersion {
				t.Errorf("Expected version %s, got %s", tt.expectedVersion, result.Version)
			}
		})
	}
}

func TestBrowserInfoMethods(t *testing.T) {
	browserInfo := BrowserInfo{
		Name:    BrowserChrome,
		Version: "120.5.1.2",
		IsBot:   false,
	}

	// Test GetMajorVersion
	if browserInfo.GetMajorVersion() != "120" {
		t.Errorf("Expected major version 120, got %s", browserInfo.GetMajorVersion())
	}

	// Test IsSupported
	minVersions := map[BrowserName]string{
		BrowserChrome:  "119.0.0.0",
		BrowserFirefox: "115.0.0.0",
	}

	if !browserInfo.IsSupported(minVersions) {
		t.Error("Expected browser to be supported")
	}

	// Test with unsupported version
	browserInfo.Version = "118.0.0.0"
	if browserInfo.IsSupported(minVersions) {
		t.Error("Expected browser to be unsupported")
	}

	// Test with bot
	browserInfo.IsBot = true
	if browserInfo.IsSupported(minVersions) {
		t.Error("Expected bot to be unsupported")
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.1", "1.0.0", 1},
		{"1.0.0", "1.0.1", -1},
		{"2.0.0", "1.9.9", 1},
		{"1.10.0", "1.9.0", 1},
		{"120.0.0.0", "119.0.0.0", 1},
		{"119.0.0.0", "120.0.0.0", -1},
		{"", "1.0.0", -1},
		{"1.0.0", "", 1},
		{"", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.v1+"_vs_"+tt.v2, func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("compareVersions(%s, %s) = %d, expected %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestAPIFunctions(t *testing.T) {
	// Test ParseBrowser
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	result := ParseBrowser(userAgent)
	if result.Name != BrowserChrome {
		t.Errorf("ParseBrowser failed, expected Chrome, got %s", result.Name)
	}

	// Test GetBrowserFamily
	family := GetBrowserFamily(userAgent)
	if family != "chromium" {
		t.Errorf("GetBrowserFamily failed, expected chromium, got %s", family)
	}

	// Test with HTTP request
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", userAgent)
	// Add standard browser headers to avoid false positives
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")

	// Test ParseBrowserFromHTTPRequest
	result2 := ParseBrowserFromHTTPRequest(req)
	if result2.Name != BrowserChrome {
		t.Errorf("ParseBrowserFromHTTPRequest failed, expected Chrome, got %s", result2.Name)
	}

	// Test IsMobileBrowser
	if IsMobileBrowser(req) {
		t.Error("Expected desktop browser, got mobile")
	}

	// Test with mobile user agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1")
	if !IsMobileBrowser(req) {
		t.Error("Expected mobile browser, got desktop")
	}

	// Test IsSupportedBrowser
	req.Header.Set("User-Agent", userAgent)
	minVersions := map[BrowserName]string{
		BrowserChrome: "119.0.0.0",
	}
	if !IsSupportedBrowser(req, minVersions) {
		t.Error("Expected browser to be supported")
	}

	// Test GetBrowserInfo
	browserInfo, botResult, err := GetBrowserInfo(req)
	if err != nil {
		t.Errorf("GetBrowserInfo returned error: %v", err)
	}
	if browserInfo.Name != BrowserChrome {
		t.Errorf("GetBrowserInfo browser parsing failed")
	}
	if botResult.Bot {
		t.Error("GetBrowserInfo incorrectly identified legitimate browser as bot")
	}
}

func BenchmarkParseBrowserFromUserAgent(b *testing.B) {
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseBrowserFromUserAgent(userAgent)
	}
}
