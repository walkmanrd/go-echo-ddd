//go:build unit

package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type HealthCheck struct {
	Message string
	Status  bool
}

func TestHealthCheckControllerIndex(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	var resultForTest, resultActual HealthCheck

	healthCheckController := &HealthCheckController{}
	json.Unmarshal([]byte(`{"message":"success","status":true}`), &resultForTest)

	if assert.NoError(t, healthCheckController.Index(c)) {
		json.Unmarshal([]byte(string(rec.Body.String())), &resultActual)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, resultForTest, resultActual)
	}
}
