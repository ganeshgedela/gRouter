package web

import (
	"time"

	"grouter/pkg/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware creates CORS middleware from config
func CORSMiddleware(cfg config.CORSConfig) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowAllOrigins:  false,
		AllowCredentials: cfg.AllowCredentials,
	}

	// Set allowed origins
	if len(cfg.AllowedOrigins) > 0 {
		corsConfig.AllowOrigins = cfg.AllowedOrigins
	} else {
		corsConfig.AllowAllOrigins = true
	}

	// Set allowed methods
	if len(cfg.AllowedMethods) > 0 {
		corsConfig.AllowMethods = cfg.AllowedMethods
	} else {
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}

	// Set allowed headers
	if len(cfg.AllowedHeaders) > 0 {
		corsConfig.AllowHeaders = cfg.AllowedHeaders
	} else {
		corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	}

	// Set exposed headers
	if len(cfg.ExposedHeaders) > 0 {
		corsConfig.ExposeHeaders = cfg.ExposedHeaders
	}

	// Set max age
	if cfg.MaxAge > 0 {
		corsConfig.MaxAge = time.Duration(cfg.MaxAge) * time.Second
	} else {
		corsConfig.MaxAge = 12 * time.Hour // Default
	}

	return cors.New(corsConfig)
}
