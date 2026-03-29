package config

import (
	"time"

	"github.com/gin-contrib/cors"
)

func (c Config) CORS() cors.Config {
	// Create a CORS config with the configured allowed origins
	return cors.Config{
		AllowOrigins:           c.AllowedOrigins,
		AllowMethods:           []string{"GET", "HEAD", "OPTIONS"},
		AllowHeaders:           []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:          []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials:       false,
		AllowWildcard:          false,
		AllowBrowserExtensions: false,
		AllowWebSockets:        false,
		AllowPrivateNetwork:    false,
		MaxAge:                 48 * time.Hour,
	}
}
