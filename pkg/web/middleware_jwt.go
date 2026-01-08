package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Enabled   bool
	Secret    string
	PublicKey string // For RS256
	Algorithm string // HS256, RS256, etc.
	Issuer    string
	Audience  string
	Logger    *zap.Logger // For logging auth failures
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// JWTAuthMiddleware validates JWT tokens
func JWTAuthMiddleware(cfg JWTConfig) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "MISSING_TOKEN",
				Message: "authorization header is required",
				Code:    "UNAUTHORIZED",
			})
			return
		}

		// Parse Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "INVALID_TOKEN_FORMAT",
				Message: "authorization header must be Bearer {token}",
				Code:    "UNAUTHORIZED",
			})
			return
		}

		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Verify signing algorithm
			if token.Method.Alg() != cfg.Algorithm {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return key based on algorithm
			switch cfg.Algorithm {
			case "HS256", "HS384", "HS512":
				return []byte(cfg.Secret), nil
			case "RS256", "RS384", "RS512":
				// TODO: Load public key from cfg.PublicKey
				return nil, fmt.Errorf("RS algorithms not yet supported")
			default:
				return nil, fmt.Errorf("unsupported algorithm: %s", cfg.Algorithm)
			}
		})

		if err != nil {
			// Log detailed error but return generic message to client
			if cfg.Logger != nil {
				cfg.Logger.Error("JWT token validation failed",
					zap.Error(err),
					zap.String("request_id", c.GetString("request_id")),
				)
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "INVALID_TOKEN",
				Message: "invalid or expired token",
				Code:    "UNAUTHORIZED",
			})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "INVALID_TOKEN",
				Message: "token is not valid",
				Code:    "UNAUTHORIZED",
			})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "INVALID_CLAIMS",
				Message: "failed to parse token claims",
				Code:    "UNAUTHORIZED",
			})
			return
		}

		// Verify issuer if configured
		if cfg.Issuer != "" && claims.Issuer != cfg.Issuer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "INVALID_ISSUER",
				Message: "token issuer does not match",
				Code:    "UNAUTHORIZED",
			})
			return
		}

		// Verify audience if configured
		if cfg.Audience != "" {
			found := false
			for _, aud := range claims.Audience {
				if aud == cfg.Audience {
					found = true
					break
				}
			}
			if !found {
				c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
					Error:   "INVALID_AUDIENCE",
					Message: "token audience does not match",
					Code:    "UNAUTHORIZED",
				})
				return
			}
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_roles", claims.Roles)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// GenerateJWT generates a JWT token for testing/development
func GenerateJWT(userID, email string, roles []string, secret string, expiresIn time.Duration) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
