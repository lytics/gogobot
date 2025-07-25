package main

import (
	"fmt"
	"net/http"

	"gogobot"
)

func main() {
	fmt.Println("=== GPT and AI Agent Detection Examples ===\n")

	// Example 1: Various GPT and AI Agent User Agents
	testUserAgents := []struct {
		name      string
		userAgent string
	}{
		{
			name:      "OpenAI GPTBot",
			userAgent: "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.0; +https://openai.com/gptbot)",
		},
		{
			name:      "ChatGPT User Agent",
			userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) ChatGPT-User/1.0",
		},
		{
			name:      "OpenAI Bot",
			userAgent: "OpenAI-Bot/1.0 (+https://openai.com/bot)",
		},
		{
			name:      "Claude Web Agent",
			userAgent: "Mozilla/5.0 (compatible; Claude-Web/1.0; +https://anthropic.com)",
		},
		{
			name:      "Generic AI Agent",
			userAgent: "AI-Agent/1.0 (Language Model Assistant)",
		},
		{
			name:      "GPT-4 Agent",
			userAgent: "GPT-4/Agent/1.0 (Conversational AI)",
		},
		{
			name:      "Bard/Gemini Agent",
			userAgent: "Mozilla/5.0 (compatible; Bard/1.0; Google AI)",
		},
		{
			name:      "Language Model Bot",
			userAgent: "LLM-Bot/1.0 (Large Language Model)",
		},
		{
			name:      "Regular Search Bot (Googlebot)",
			userAgent: "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		},
		{
			name:      "Regular Browser (Chrome)",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		},
		{
			name:      "Command Line Tool (curl)",
			userAgent: "curl/7.81.0",
		},
	}

	for i, test := range testUserAgents {
		fmt.Printf("Example %d - %s:\n", i+1, test.name)
		fmt.Printf("User Agent: %s\n", test.userAgent)

		// Check if it's a GPT/AI agent
		isGPT, gptKind := gogobot.IsGPTAgent(test.userAgent)
		fmt.Printf("Is GPT/AI Agent: %v", isGPT)
		if isGPT {
			fmt.Printf(" (Type: %s)", gptKind)
		}
		fmt.Println()

		// Check specific OpenAI detection
		isOpenAI := gogobot.IsOpenAIBot(test.userAgent)
		fmt.Printf("Is OpenAI Bot: %v\n", isOpenAI)

		// Check specific ChatGPT detection
		isChatGPT := gogobot.IsChatGPT(test.userAgent)
		fmt.Printf("Is ChatGPT: %v\n", isChatGPT)

		// General bot detection
		isBot, botKind := gogobot.IsBotUserAgent(test.userAgent)
		fmt.Printf("Is Bot (General): %v", isBot)
		if isBot {
			fmt.Printf(" (Type: %s)", botKind)
		}
		fmt.Println()

		fmt.Println()
	}

	// Example 2: HTTP Request Analysis
	fmt.Println("=== HTTP Request Analysis ===\n")

	// Simulate various HTTP requests
	testRequests := []struct {
		name      string
		userAgent string
	}{
		{
			name:      "ChatGPT Browsing Request",
			userAgent: "ChatGPT-User/1.0",
		},
		{
			name:      "OpenAI GPTBot Crawler",
			userAgent: "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.0; +https://openai.com/gptbot)",
		},
		{
			name:      "Claude AI Agent",
			userAgent: "Claude-Web/1.0 (Anthropic AI Assistant)",
		},
		{
			name:      "Regular Browser Request",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		},
	}

	for i, test := range testRequests {
		fmt.Printf("Request %d - %s:\n", i+1, test.name)

		req, _ := http.NewRequest("GET", "/api/data", nil)
		req.Header.Set("User-Agent", test.userAgent)

		// Add realistic headers for browser requests
		if test.name == "Regular Browser Request" {
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
			req.Header.Set("Accept-Language", "en-US,en;q=0.5")
			req.Header.Set("Accept-Encoding", "gzip, deflate")
			req.Header.Set("Connection", "keep-alive")
		}

		// Check if request is from GPT/AI agent
		isGPT, gptKind := gogobot.IsGPTRequest(req)
		fmt.Printf("  Is GPT Request: %v", isGPT)
		if isGPT {
			fmt.Printf(" (Type: %s)", gptKind)
		}
		fmt.Println()

		// Comprehensive AI agent analysis
		isAI, agentType, botResult, err := gogobot.GetAIAgentInfo(req)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		} else {
			fmt.Printf("  Is AI Agent: %v", isAI)
			if isAI {
				fmt.Printf(" (Agent Type: %s)", agentType)
			}
			fmt.Println()
			fmt.Printf("  Bot Detection Result: %v", botResult.Bot)
			if botResult.Bot {
				fmt.Printf(" (Bot Kind: %s)", botResult.BotKind)
			}
			fmt.Println()
		}

		fmt.Println()
	}

	// Example 3: Integration with Browser Parsing
	fmt.Println("=== Combined Browser Parsing and AI Detection ===\n")

	combinedTests := []struct {
		name      string
		userAgent string
	}{
		{
			name:      "ChatGPT with Browser-like UA",
			userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) ChatGPT-User/1.0",
		},
		{
			name:      "Regular Chrome Browser",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		},
		{
			name:      "AI Agent with minimal UA",
			userAgent: "AI-Agent/1.0",
		},
	}

	for i, test := range combinedTests {
		fmt.Printf("Test %d - %s:\n", i+1, test.name)

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("User-Agent", test.userAgent)

		// Add browser headers for legitimate browser
		if test.name == "Regular Chrome Browser" {
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
			req.Header.Set("Accept-Language", "en-US,en;q=0.5")
			req.Header.Set("Accept-Encoding", "gzip, deflate")
			req.Header.Set("Connection", "keep-alive")
		}

		// Parse browser information
		browserInfo := gogobot.ParseBrowserFromHTTPRequest(req)
		fmt.Printf("  Browser: %s %s\n", browserInfo.Name, browserInfo.Version)
		fmt.Printf("  Browser Family: %s\n", browserInfo.GetBrowserFamily())
		fmt.Printf("  Is Bot (Browser Analysis): %v\n", browserInfo.IsBot)

		// Check for AI agents
		isGPT, gptKind := gogobot.IsGPTRequest(req)
		fmt.Printf("  Is GPT/AI Agent: %v", isGPT)
		if isGPT {
			fmt.Printf(" (Type: %s)", gptKind)
		}
		fmt.Println()

		// Full bot detection
		detector := gogobot.NewDetector()
		botResult, err := detector.DetectFromRequest(req)
		if err != nil {
			fmt.Printf("  Error in bot detection: %v\n", err)
		} else {
			fmt.Printf("  Full Bot Detection: %v", botResult.Bot)
			if botResult.Bot {
				fmt.Printf(" (Kind: %s)", botResult.BotKind)
			}
			fmt.Println()
		}

		fmt.Println()
	}

	// Example 4: Practical Usage for Web Applications
	fmt.Println("=== Practical Web Application Usage ===\n")

	// Simulate incoming requests to a web application
	incomingRequests := []string{
		"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.0; +https://openai.com/gptbot)",
		"ChatGPT-User/1.0",
		"Mozilla/5.0 (compatible; Claude-Web/1.0)",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"curl/7.81.0",
	}

	for i, userAgent := range incomingRequests {
		fmt.Printf("Incoming Request %d:\n", i+1)
		fmt.Printf("User-Agent: %s\n", userAgent)

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Set("User-Agent", userAgent)

		// Decision logic for handling different types of requests
		isAI, agentType, _, err := gogobot.GetAIAgentInfo(req)
		if err != nil {
			fmt.Printf("Action: Allow (detection error: %v)\n", err)
		} else if isAI {
			switch agentType {
			case gogobot.BotKindGPTBot, gogobot.BotKindOpenAI:
				fmt.Printf("Action: Allow AI crawler (OpenAI GPTBot - training data)\n")
			case gogobot.BotKindChatGPT:
				fmt.Printf("Action: Allow ChatGPT browsing (user-initiated)\n")
			case gogobot.BotKindClaude:
				fmt.Printf("Action: Allow Claude agent (Anthropic AI)\n")
			case gogobot.BotKindAIAgent:
				fmt.Printf("Action: Rate limit generic AI agent\n")
			default:
				fmt.Printf("Action: Allow unknown AI agent with monitoring\n")
			}
		} else {
			// Check if it's a regular bot or browser
			isBot, botKind := gogobot.IsBotUserAgent(userAgent)
			if isBot {
				fmt.Printf("Action: Apply bot rate limiting (Bot type: %s)\n", botKind)
			} else {
				fmt.Printf("Action: Allow normal user traffic\n")
			}
		}

		fmt.Println()
	}
}
