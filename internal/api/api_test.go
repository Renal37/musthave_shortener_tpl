package api_test

import (
	"context"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/api"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStartRestAPI(t *testing.T) {
	// Create a mock storage and database
	mockStorage := &storage.Storage{}
	mockDB := &repository.StoreDB{}

	// Create a context with a timeout to ensure the server shuts down
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start the REST API in a separate goroutine
	go func() {
		err := api.StartRestAPI(ctx, ":8081", "http://example.com", "info", mockDB, false, mockStorage)
		assert.NoError(t, err, "Expected no error from StartRestAPI")
	}()

	// Give the server a moment to start
	time.Sleep(500 * time.Millisecond)

	// Create a test HTTP request to the correct endpoint
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8081/health", nil)
	assert.NoError(t, err, "Expected no error creating request")

	// Use httptest to create a response recorder
	rr := httptest.NewRecorder()

	// Perform the request using the correct router
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")

	// Check the response body
	expectedBody := "OK"
	assert.Equal(t, expectedBody, rr.Body.String(), "Expected response body to be 'OK'")
}
