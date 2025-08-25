package middleware

import (
	"errors"
	"net/http"
	"strings"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"forum/internal/configs"
	"forum/pkg/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
    secretKey := configs.Get().Service.SecretJWT
    return func(c *gin.Context) {
        header := c.Request.Header.Get("Authorization")

        header = strings.TrimSpace(header)
        if header == "" {
            c.AbortWithError(http.StatusUnauthorized, errors.New("missing token"))
            return
        }

        parts := strings.Split(header, " ")
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            c.AbortWithError(http.StatusUnauthorized, errors.New("invalid authorization header format"))
            return
        }

        token := parts[1]
		log.Info().Msgf("Received token: %s", token)

        // Validate the token
        userID, username, err := jwt.ValidateToken(token, secretKey)
        if err != nil {
            c.AbortWithError(http.StatusUnauthorized, err)
            return
        }

        c.Set("userID", userID)
        c.Set("username", username)
        c.Next()
    }
}

func AuthRefreshMiddleware() gin.HandlerFunc {
	secretKey := configs.Get().Service.SecretJWT
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")

		header = strings.TrimSpace(header)
		if header == "" {
			c.AbortWithError(http.StatusUnauthorized, errors.New("missing token"))
			return
		}

		userID, username, err := jwt.ValidateTokenWithoutExpiry(header, secretKey)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		c.Set("userID", userID)
		c.Set("username", username)
		c.Next()
	}
}
