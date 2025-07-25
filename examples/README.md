# BotD-Go Examples

This directory contains examples demonstrating how to use the BotD-Go library for server-side bot detection.

## Examples

### 1. Basic Usage (`basic/main.go`)

Demonstrates the fundamental usage patterns:

```bash
cd examples/basic
go run main.go
```

This example shows:
- Creating a bot detector
- Analyzing different types of requests (browser vs bot)
- Quick user agent checking
- Header analysis

### 2. HTTP Middleware (`middleware/main.go`)

Shows how to integrate BotD-Go with HTTP servers:

```bash
cd examples/middleware
go run main.go
```

Then test with curl:
```bash
# This should be blocked (curl user agent)
curl http://localhost:8080/protected

# This returns detection details
curl http://localhost:8080/api/detect

# This allows bots but logs them
curl http://localhost:8080/public

# Test with a search engine crawler (should be allowed)
curl -H 'User-Agent: Googlebot/2.1' http://localhost:8080/smart-protection

# Skip detection with special header
curl -H 'X-Skip-Bot-Detection: true' http://localhost:8080/smart-protection
```

### 3. Custom Detectors (`custom-detectors/main.go`)

Demonstrates how to create and use custom detection logic:

```bash
cd examples/custom-detectors
go run main.go
```

This example shows:
- Creating custom detector functions
- Adding detectors at runtime
- IP-based detection
- Header pattern analysis
- Referer validation

## Building and Running

Make sure you have Go 1.21 or later installed:

```bash
# Clone or download the BotD-Go library
cd BotD-Go

# Run any example
cd examples/basic
go run main.go

# Or build for distribution
cd examples/middleware
go build -o server main.go
./server
```

## Integration Patterns

### Standard HTTP Handler

```go
detector := botd.Load()
http.HandleFunc("/api", detector.HandlerFunc(yourHandler))
```

### Middleware Chain

```go
detector := botd.Load()
middleware := detector.Middleware()
handler := middleware(http.HandlerFunc(yourHandler))
```

### Custom Configuration

```go
config := botd.MiddlewareConfig{
    BlockBots: true,
    OnBotDetected: func(w http.ResponseWriter, r *http.Request, result *botd.BotDetectionResult) {
        log.Printf("Bot blocked: %s", result.BotKind)
    },
}
middleware := detector.MiddlewareWithConfig(config)
```

### One-off Detection

```go
result, err := botd.Detect(request)
if err == nil && result.Bot {
    // Handle bot traffic
}
```

## Common Use Cases

1. **API Protection**: Block automated API access while allowing legitimate clients
2. **Rate Limiting**: Apply different rate limits for bots vs humans
3. **Analytics**: Track bot vs human traffic separately
4. **SEO**: Allow search engine crawlers while blocking malicious bots
5. **Security**: Prevent scraping and automated attacks

## Performance Considerations

- Use `QuickCheck()` for high-traffic scenarios (faster, fewer detectors)
- Implement caching if running detection on every request
- Consider running detection asynchronously for non-critical paths
- Use custom detectors for domain-specific bot patterns 