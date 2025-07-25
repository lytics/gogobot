package gogobot

import (
	"net/http"
	"testing"
)

func TestIsGPTAgent(t *testing.T) {
	tests := []struct {
		name          string
		userAgent     string
		expectedIsGPT bool
		expectedKind  BotKind
	}{
		// OpenAI GPTBot
		{
			name:          "OpenAI GPTBot",
			userAgent:     "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.0; +https://openai.com/gptbot)",
			expectedIsGPT: true,
			expectedKind:  BotKindGPTBot,
		},
		{
			name:          "GPT-Bot variant",
			userAgent:     "GPT-Bot/1.0",
			expectedIsGPT: true,
			expectedKind:  BotKindGPTBot,
		},

		// ChatGPT variants
		{
			name:          "ChatGPT User",
			userAgent:     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) ChatGPT-User/1.0",
			expectedIsGPT: true,
			expectedKind:  BotKindChatGPT,
		},
		{
			name:          "ChatGPT Simple",
			userAgent:     "ChatGPT/1.0",
			expectedIsGPT: true,
			expectedKind:  BotKindChatGPT,
		},
		{
			name:          "OpenAI ChatGPT",
			userAgent:     "Mozilla/5.0 (compatible; openai-chatgpt)",
			expectedIsGPT: true,
			expectedKind:  BotKindChatGPT,
		},

		// OpenAI variants
		{
			name:          "OpenAI Bot",
			userAgent:     "OpenAI-Bot/1.0",
			expectedIsGPT: true,
			expectedKind:  BotKindOpenAI,
		},
		{
			name:          "OpenAI Crawler",
			userAgent:     "Mozilla/5.0 (compatible; openai-crawler/1.0; +https://openai.com/bot)",
			expectedIsGPT: true,
			expectedKind:  BotKindOpenAI,
		},
		{
			name:          "OpenAI Simple",
			userAgent:     "OpenAI/1.0",
			expectedIsGPT: true,
			expectedKind:  BotKindOpenAI,
		},

		// Claude/Anthropic variants
		{
			name:          "Claude Web",
			userAgent:     "Mozilla/5.0 (compatible; Claude-Web/1.0; +https://anthropic.com)",
			expectedIsGPT: true,
			expectedKind:  BotKindClaude,
		},
		{
			name:          "Claude Simple",
			userAgent:     "Claude/2.0",
			expectedIsGPT: true,
			expectedKind:  BotKindClaude,
		},
		{
			name:          "Anthropic Agent",
			userAgent:     "Anthropic-Agent/1.0",
			expectedIsGPT: true,
			expectedKind:  BotKindClaude,
		},

		// Generic AI Agents
		{
			name:          "AI Agent",
			userAgent:     "AI-Agent/1.0 (Language Model)",
			expectedIsGPT: true,
			expectedKind:  BotKindAIAgent,
		},
		{
			name:          "Language Model",
			userAgent:     "Mozilla/5.0 (compatible; Language Model Bot)",
			expectedIsGPT: true,
			expectedKind:  BotKindAIAgent,
		},
		{
			name:          "LLM Agent",
			userAgent:     "LLM-Agent/1.0",
			expectedIsGPT: true,
			expectedKind:  BotKindAIAgent,
		},
		{
			name:          "GPT-4 Agent",
			userAgent:     "GPT-4/Agent/1.0",
			expectedIsGPT: true,
			expectedKind:  BotKindAIAgent,
		},
		{
			name:          "Bard Agent",
			userAgent:     "Mozilla/5.0 (compatible; Bard/1.0)",
			expectedIsGPT: true,
			expectedKind:  BotKindAIAgent,
		},
		{
			name:          "Gemini Pro",
			userAgent:     "Gemini-Pro/1.0 (AI Assistant)",
			expectedIsGPT: true,
			expectedKind:  BotKindAIAgent,
		},

		// Non-GPT bots (should return false)
		{
			name:          "Regular Bot",
			userAgent:     "SomeBot/1.0",
			expectedIsGPT: false,
			expectedKind:  "",
		},
		{
			name:          "Googlebot",
			userAgent:     "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
			expectedIsGPT: false,
			expectedKind:  "",
		},
		{
			name:          "Curl",
			userAgent:     "curl/7.81.0",
			expectedIsGPT: false,
			expectedKind:  "",
		},

		// Regular browsers (should return false)
		{
			name:          "Chrome Browser",
			userAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			expectedIsGPT: false,
			expectedKind:  "",
		},
		{
			name:          "Firefox Browser",
			userAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
			expectedIsGPT: false,
			expectedKind:  "",
		},

		// Edge cases
		{
			name:          "Empty User Agent",
			userAgent:     "",
			expectedIsGPT: false,
			expectedKind:  "",
		},
		{
			name:          "GPT in context but not agent",
			userAgent:     "Mozilla/5.0 (compatible; SearchBot) visiting GPT website",
			expectedIsGPT: false,
			expectedKind:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isGPT, kind := IsGPTAgent(tt.userAgent)

			if isGPT != tt.expectedIsGPT {
				t.Errorf("IsGPTAgent(%s) isGPT = %v, expected %v", tt.userAgent, isGPT, tt.expectedIsGPT)
			}

			if kind != tt.expectedKind {
				t.Errorf("IsGPTAgent(%s) kind = %s, expected %s", tt.userAgent, kind, tt.expectedKind)
			}
		})
	}
}

func TestIsGPTRequest(t *testing.T) {
	tests := []struct {
		name          string
		userAgent     string
		expectedIsGPT bool
		expectedKind  BotKind
	}{
		{
			name:          "ChatGPT Request",
			userAgent:     "ChatGPT-User/1.0",
			expectedIsGPT: true,
			expectedKind:  BotKindChatGPT,
		},
		{
			name:          "Regular Browser Request",
			userAgent:     "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			expectedIsGPT: false,
			expectedKind:  "",
		},
		{
			name:          "No User Agent",
			userAgent:     "",
			expectedIsGPT: false,
			expectedKind:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.userAgent != "" {
				req.Header.Set("User-Agent", tt.userAgent)
			}

			isGPT, kind := IsGPTRequest(req)

			if isGPT != tt.expectedIsGPT {
				t.Errorf("IsGPTRequest(%s) isGPT = %v, expected %v", tt.userAgent, isGPT, tt.expectedIsGPT)
			}

			if kind != tt.expectedKind {
				t.Errorf("IsGPTRequest(%s) kind = %s, expected %s", tt.userAgent, kind, tt.expectedKind)
			}
		})
	}
}

func TestIsChatGPT(t *testing.T) {
	tests := []struct {
		name            string
		userAgent       string
		expectedChatGPT bool
	}{
		{
			name:            "ChatGPT User",
			userAgent:       "ChatGPT-User/1.0",
			expectedChatGPT: true,
		},
		{
			name:            "OpenAI Bot",
			userAgent:       "OpenAI-Bot/1.0",
			expectedChatGPT: true,
		},
		{
			name:            "ChatGPT Simple",
			userAgent:       "ChatGPT/1.0",
			expectedChatGPT: true,
		},
		{
			name:            "GPTBot (not ChatGPT)",
			userAgent:       "GPTBot/1.0",
			expectedChatGPT: false,
		},
		{
			name:            "Claude (not ChatGPT)",
			userAgent:       "Claude/1.0",
			expectedChatGPT: false,
		},
		{
			name:            "Regular Browser",
			userAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			expectedChatGPT: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsChatGPT(tt.userAgent)

			if result != tt.expectedChatGPT {
				t.Errorf("IsChatGPT(%s) = %v, expected %v", tt.userAgent, result, tt.expectedChatGPT)
			}
		})
	}
}

func TestIsOpenAIBot(t *testing.T) {
	tests := []struct {
		name           string
		userAgent      string
		expectedOpenAI bool
	}{
		{
			name:           "GPTBot",
			userAgent:      "GPTBot/1.0",
			expectedOpenAI: true,
		},
		{
			name:           "OpenAI Bot",
			userAgent:      "OpenAI-Bot/1.0",
			expectedOpenAI: true,
		},
		{
			name:           "ChatGPT",
			userAgent:      "ChatGPT-User/1.0",
			expectedOpenAI: true,
		},
		{
			name:           "OpenAI Crawler",
			userAgent:      "OpenAI-Crawler/1.0",
			expectedOpenAI: true,
		},
		{
			name:           "Claude (not OpenAI)",
			userAgent:      "Claude/1.0",
			expectedOpenAI: false,
		},
		{
			name:           "Generic AI Agent (not OpenAI)",
			userAgent:      "AI-Agent/1.0",
			expectedOpenAI: false,
		},
		{
			name:           "Regular Browser",
			userAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			expectedOpenAI: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsOpenAIBot(tt.userAgent)

			if result != tt.expectedOpenAI {
				t.Errorf("IsOpenAIBot(%s) = %v, expected %v", tt.userAgent, result, tt.expectedOpenAI)
			}
		})
	}
}

func TestGetAIAgentInfo(t *testing.T) {
	tests := []struct {
		name              string
		userAgent         string
		setupHeaders      func(*http.Request)
		expectedIsAI      bool
		expectedAgentType BotKind
		expectedBotResult bool
	}{
		{
			name:              "ChatGPT Request",
			userAgent:         "ChatGPT-User/1.0",
			expectedIsAI:      true,
			expectedAgentType: BotKindChatGPT,
			expectedBotResult: true,
		},
		{
			name:              "OpenAI GPTBot",
			userAgent:         "GPTBot/1.0",
			expectedIsAI:      true,
			expectedAgentType: BotKindGPTBot,
			expectedBotResult: true,
		},
		{
			name:              "Claude Agent",
			userAgent:         "Claude-Web/1.0",
			expectedIsAI:      true,
			expectedAgentType: BotKindClaude,
			expectedBotResult: true,
		},
		{
			name:              "Generic AI Agent",
			userAgent:         "AI-Agent/1.0 (Language Model)",
			expectedIsAI:      true,
			expectedAgentType: BotKindAIAgent,
			expectedBotResult: true,
		},
		{
			name:              "Regular Bot (not AI)",
			userAgent:         "SomeBot/1.0",
			expectedIsAI:      false,
			expectedAgentType: "",
			expectedBotResult: true,
		},
		{
			name:      "Regular Browser",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			setupHeaders: func(req *http.Request) {
				// Add standard browser headers to avoid false positives
				req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
				req.Header.Set("Accept-Language", "en-US,en;q=0.5")
				req.Header.Set("Accept-Encoding", "gzip, deflate")
				req.Header.Set("Connection", "keep-alive")
			},
			expectedIsAI:      false,
			expectedAgentType: "",
			expectedBotResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Set("User-Agent", tt.userAgent)

			if tt.setupHeaders != nil {
				tt.setupHeaders(req)
			}

			isAI, agentType, botResult, err := GetAIAgentInfo(req)

			if err != nil {
				t.Errorf("GetAIAgentInfo(%s) returned error: %v", tt.userAgent, err)
			}

			if isAI != tt.expectedIsAI {
				t.Errorf("GetAIAgentInfo(%s) isAI = %v, expected %v", tt.userAgent, isAI, tt.expectedIsAI)
			}

			if agentType != tt.expectedAgentType {
				t.Errorf("GetAIAgentInfo(%s) agentType = %s, expected %s", tt.userAgent, agentType, tt.expectedAgentType)
			}

			if botResult.Bot != tt.expectedBotResult {
				t.Errorf("GetAIAgentInfo(%s) botResult.Bot = %v, expected %v", tt.userAgent, botResult.Bot, tt.expectedBotResult)
			}
		})
	}
}

func TestGPTDetectionIntegration(t *testing.T) {
	// Test that GPT detection works with the existing bot detection system
	userAgent := "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.0; +https://openai.com/gptbot)"

	// Test with IsBotUserAgent
	isBot, botKind := IsBotUserAgent(userAgent)
	if !isBot {
		t.Error("Expected GPTBot to be detected as bot by IsBotUserAgent")
	}
	if botKind != BotKindGPTBot {
		t.Errorf("Expected BotKindGPTBot, got %s", botKind)
	}

	// Test with IsGPTAgent
	isGPT, gptKind := IsGPTAgent(userAgent)
	if !isGPT {
		t.Error("Expected GPTBot to be detected as GPT agent")
	}
	if gptKind != BotKindGPTBot {
		t.Errorf("Expected BotKindGPTBot, got %s", gptKind)
	}

	// Test with full detector
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", userAgent)

	detector := NewDetector()
	result, err := detector.DetectFromRequest(req)
	if err != nil {
		t.Errorf("Detector returned error: %v", err)
	}
	if !result.Bot {
		t.Error("Expected GPTBot to be detected as bot by full detector")
	}
	if result.BotKind != BotKindGPTBot {
		t.Errorf("Expected BotKindGPTBot, got %s", result.BotKind)
	}
}

func BenchmarkIsGPTAgent(b *testing.B) {
	userAgent := "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.0; +https://openai.com/gptbot)"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsGPTAgent(userAgent)
	}
}

func BenchmarkGetAIAgentInfo(b *testing.B) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("User-Agent", "ChatGPT-User/1.0")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetAIAgentInfo(req)
	}
}
