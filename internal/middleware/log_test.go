package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggerMiddleware(t *testing.T) {
	var logBuffer bytes.Buffer
	writer := zapcore.AddSync(&logBuffer)
	encoderCfg := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), writer, zap.InfoLevel)
	logger := zap.New(core).Sugar()

	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(LoggerMiddleware(logger))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Logger test!")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, `"method":"GET"`)
	assert.Contains(t, logOutput, `"path":"/test"`)
	assert.Contains(t, logOutput, `"statusCode":200`)
	assert.Contains(t, logOutput, `"duration"`)
	assert.Contains(t, logOutput, `"contentLength"`)
}
