package rest

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func loggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()

		endTime := time.Since(startTime).String()
		log.WithFields(log.Fields{
			"method":   ctx.Request.Method,
			"URL":      ctx.Request.URL,
			"duration": endTime,
		}).Info("Middleware: loggerMiddleware")
	}
}
