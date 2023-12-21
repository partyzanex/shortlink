package logger

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Middleware(l *Logger) gin.HandlerFunc {
	return logger.SetLogger(logger.WithLogger(func(_ *gin.Context, _ zerolog.Logger) zerolog.Logger {
		return l.zl.With().Logger()
	}))
}
