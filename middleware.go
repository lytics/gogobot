package gogobot

import (
	"context"
	"net/http"
)

// MiddlewareConfig holds configuration for the middleware
type MiddlewareConfig struct {
	// SkipFunc allows skipping detection for specific requests
	SkipFunc func(*http.Request) bool
	// OnBotDetected is called when a bot is detected
	OnBotDetected func(http.ResponseWriter, *http.Request, *BotDetectionResult)
	// OnError is called when an error occurs during detection
	OnError func(http.ResponseWriter, *http.Request, error)
	// BlockBots determines if detected bots should be blocked
	BlockBots bool
	// BlockedStatusCode is the HTTP status code to return for blocked bots
	BlockedStatusCode int
	// BlockedMessage is the message to return for blocked bots
	BlockedMessage string
}

// DefaultMiddlewareConfig returns a default middleware configuration
func DefaultMiddlewareConfig() MiddlewareConfig {
	return MiddlewareConfig{
		SkipFunc:          nil,
		OnBotDetected:     nil,
		OnError:           nil,
		BlockBots:         false,
		BlockedStatusCode: http.StatusForbidden,
		BlockedMessage:    "Bot traffic is not allowed",
	}
}

// Middleware returns an HTTP middleware function
func (d *BotDetector) Middleware() func(http.Handler) http.Handler {
	return d.MiddlewareWithConfig(DefaultMiddlewareConfig())
}

// MiddlewareWithConfig returns an HTTP middleware function with custom configuration
func (d *BotDetector) MiddlewareWithConfig(config MiddlewareConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip detection if configured
			if config.SkipFunc != nil && config.SkipFunc(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Perform bot detection
			result, err := d.DetectFromRequest(r)
			if err != nil {
				if config.OnError != nil {
					config.OnError(w, r, err)
					return
				}
				// Continue processing if no error handler is configured
				next.ServeHTTP(w, r)
				return
			}

			// Store result in context
			ctx := context.WithValue(r.Context(), DetectionResultKey, &result)
			ctx = context.WithValue(ctx, ComponentsKey, d.GetComponents())
			r = r.WithContext(ctx)

			// Handle bot detection
			if result.Bot {
				if config.OnBotDetected != nil {
					config.OnBotDetected(w, r, &result)
					return
				}

				if config.BlockBots {
					// Ensure we have a valid status code
					statusCode := config.BlockedStatusCode
					if statusCode == 0 {
						statusCode = http.StatusForbidden
					}
					message := config.BlockedMessage
					if message == "" {
						message = "Bot traffic is not allowed"
					}
					http.Error(w, message, statusCode)
					return
				}
			}

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// HandlerFunc is a convenience function that wraps a http.HandlerFunc with bot detection
func (d *BotDetector) HandlerFunc(handler http.HandlerFunc) http.HandlerFunc {
	middleware := d.Middleware()
	wrappedHandler := middleware(handler)
	return wrappedHandler.ServeHTTP
}

// HandlerFuncWithConfig is a convenience function that wraps a http.HandlerFunc with bot detection and custom config
func (d *BotDetector) HandlerFuncWithConfig(config MiddlewareConfig, handler http.HandlerFunc) http.HandlerFunc {
	middleware := d.MiddlewareWithConfig(config)
	wrappedHandler := middleware(handler)
	return wrappedHandler.ServeHTTP
}
