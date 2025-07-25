# gogobot

gogobot is a Golang port of the [BotD JavaScript library](https://github.com/fingerprintjs/BotD) for server-side bot detection.

## Features

- Server-side bot detection for web applications
- **GPT and AI agent detection** (ChatGPT, GPTBot, Claude, etc.)
- Browser parsing and version extraction from user agents
- HTTP middleware for easy integration
- Supports detection of various automation tools and frameworks
- Mobile device detection
- Browser version compatibility checking
- Lightweight and fast detection algorithms
- Compatible with popular Go web frameworks

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/lytics/gogobot"
)

func main() {
    // Create a new bot detector
    detector := gogobot.NewDetector()
    
    // Analyze an HTTP request
    req, _ := http.NewRequest("GET", "/", nil)
    req.Header.Set("User-Agent", "Mozilla/5.0...")
    
    result, err := detector.DetectFromRequest(req)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    if result.Bot {
        fmt.Printf("Bot detected: %s\n", result.BotKind)
    } else {
        fmt.Println("Human traffic detected")
    }
}
```

### HTTP Middleware

```go
package main

import (
    "net/http"
    
    "github.com/lytics/gogobot"
    "github.com/gorilla/mux"
)

func main() {
    detector := gogobot.NewDetector()
    
    r := mux.NewRouter()
    
    // Use bot detection middleware
    r.Use(detector.Middleware())
    
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Access detection result from context
        if result, ok := gogobot.GetResultFromContext(r.Context()); ok {
            if result.Bot {
                http.Error(w, "Bot traffic not allowed", http.StatusForbidden)
                return
            }
        }
        w.Write([]byte("Welcome, human!"))
    })
    
    http.ListenAndServe(":8080", r)
}
```

### Browser Parsing

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/lytics/gogobot"
)

func main() {
    // Parse browser from user agent string
    userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
    
    browserInfo := gogobot.ParseBrowser(userAgent)
    fmt.Printf("Browser: %s %s\n", browserInfo.Name, browserInfo.Version)
    fmt.Printf("Browser Family: %s\n", browserInfo.GetBrowserFamily())
    fmt.Printf("Is Mobile: %v\n", browserInfo.IsMobile())
    
    // Parse browser from HTTP request
    req, _ := http.NewRequest("GET", "/", nil)
    req.Header.Set("User-Agent", userAgent)
    
    browserInfo = gogobot.ParseBrowserFromHTTPRequest(req)
    
    // Check browser version compatibility
    minVersions := map[gogobot.BrowserName]string{
        gogobot.BrowserChrome:  "100.0.0.0",
        gogobot.BrowserFirefox: "100.0.0.0",
        gogobot.BrowserSafari:  "15.0",
    }
    
    isSupported := browserInfo.IsSupported(minVersions)
    fmt.Printf("Browser supported: %v\n", isSupported)
    
    // Combined analysis - bot detection + browser parsing
    browserInfo, botResult, err := gogobot.GetBrowserInfo(req)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    if botResult.Bot {
        fmt.Printf("Bot detected: %s\n", botResult.BotKind)
    } else {
        fmt.Printf("Legitimate browser: %s %s\n", browserInfo.Name, browserInfo.Version)
    }
}
```

### GPT and AI Agent Detection

```go
package main

import (
    "fmt"
    "net/http"
    
    "github.com/lytics/gogobot"
)

func main() {
    // Detect GPT agents from user agent strings
    userAgent := "Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; GPTBot/1.0; +https://openai.com/gptbot)"
    
    isGPT, gptKind := gogobot.IsGPTAgent(userAgent)
    if isGPT {
        fmt.Printf("GPT Agent detected: %s\n", gptKind)
    }
    
    // Check for specific AI agents
    isChatGPT := gogobot.IsChatGPT("ChatGPT-User/1.0")
    isOpenAI := gogobot.IsOpenAIBot("GPTBot/1.0")
    
    fmt.Printf("Is ChatGPT: %v\n", isChatGPT)
    fmt.Printf("Is OpenAI Bot: %v\n", isOpenAI)
    
    // Analyze HTTP requests for AI agents
    req, _ := http.NewRequest("GET", "/", nil)
    req.Header.Set("User-Agent", userAgent)
    
    isAI, agentType, botResult, err := gogobot.GetAIAgentInfo(req)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    if isAI {
        fmt.Printf("AI Agent: %s\n", agentType)
        
        // Handle different AI agents appropriately
        switch agentType {
        case gogobot.BotKindGPTBot:
            fmt.Println("Action: Allow OpenAI training crawler")
        case gogobot.BotKindChatGPT:
            fmt.Println("Action: Allow ChatGPT browsing")
        case gogobot.BotKindClaude:
            fmt.Println("Action: Allow Claude AI assistant")
        default:
            fmt.Println("Action: Monitor unknown AI agent")
        }
    }
}
```

## Supported Detection Methods

This Go port focuses on server-side signals available from HTTP requests:

- **User Agent Analysis**: Detection of common automation tool signatures
- **Header Fingerprinting**: Analysis of HTTP header patterns
- **Request Timing**: Detection of unusually fast request patterns
- **IP Analysis**: Identification of datacenter and cloud provider IPs
- **Header Consistency**: Detection of inconsistent header combinations

## Architecture

The library follows the same architectural patterns as the original JavaScript version:

- **Sources**: Collect data from HTTP requests and server environment
- **Detectors**: Analyze collected data to identify bot patterns
- **Components**: Structured data with state management and error handling
- **Results**: Standardized detection results with confidence levels

## License

MIT

## Contributing

Contributions are welcome! Please read the contributing guidelines and submit pull requests. 