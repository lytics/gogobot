# gogobot

gogobot is a Golang port of the [BotD JavaScript library](https://github.com/fingerprintjs/BotD) for server-side bot detection.

## Features

- Server-side bot detection for web applications
- HTTP middleware for easy integration
- Supports detection of various automation tools and frameworks
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