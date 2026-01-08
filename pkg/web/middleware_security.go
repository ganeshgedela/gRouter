package web

import (
	"grouter/pkg/config"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware(cfg config.SecurityConfig) gin.HandlerFunc {
	if !cfg.Enabled {
		// Return no-op middleware if disabled
		return func(c *gin.Context) {
			c.Next()
		}
	}

	secureConfig := secure.Config{
		// Security Headers
		BrowserXssFilter:   cfg.XSSProtection != "",
		ContentTypeNosniff: cfg.ContentTypeNosniff != "",
		FrameDeny:          cfg.XFrameOptions == "DENY",

		// HSTS
		STSSeconds:           int64(cfg.HSTSMaxAge),
		STSIncludeSubdomains: !cfg.HSTSExcludeSubdomains,

		// Custom headers
		ContentSecurityPolicy: cfg.ContentSecurityPolicy,
		ReferrerPolicy:        cfg.ReferrerPolicy,

		IsDevelopment: false,
	}

	// Override frame options if specified
	if cfg.XFrameOptions != "" && cfg.XFrameOptions != "DENY" {
		secureConfig.FrameDeny = false
		secureConfig.CustomFrameOptionsValue = cfg.XFrameOptions
	}

	return secure.New(secureConfig)
}
