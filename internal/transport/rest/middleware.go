package rest

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()

		endTime := time.Since(startTime).String()
		log.WithFields(log.Fields{
			"method":   c.Request.Method,
			"URL":      c.Request.URL,
			"duration": endTime,
		}).Info("middleware: loggerMiddleware")
	}
}

func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getTokenFromRequest(c.Request)
		if err != nil {
			log.WithField("middleware:", "authMiddleware").Error(err)
			http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
			c.Abort()
			return
		}

		userId, err := h.UserService.ParseToken(c.Request.Context(), token)
		if err != nil {
			log.WithField("middleware:", "authMiddleware").Error(err)
			http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}

func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
