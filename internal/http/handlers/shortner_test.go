package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortner-module/internal/http/handlers"
	"url-shortner-module/internal/http/validators"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDoShort(t *testing.T) {
	e := echo.New()
	e.Validator = &validators.CustomValidator{Validator: validator.New()}
	handler := handlers.NewShortnerHandler()

	req := httptest.NewRequest(http.MethodPost, "/v1/url/short", strings.NewReader(`{"url": "http://example.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.DoShort()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDoShort_onSameUrls(t *testing.T) {
	e := echo.New()
	e.Validator = &validators.CustomValidator{Validator: validator.New()}
	handler := handlers.NewShortnerHandler()

	t.Run("Same URL generates same short URL", func(t *testing.T) {
		req1 := httptest.NewRequest(http.MethodPost, "/v1/url/short", strings.NewReader(`{"url": "http://example.com"}`))
		req1.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec1 := httptest.NewRecorder()
		c1 := e.NewContext(req1, rec1)

		assert.NoError(t, handler.DoShort()(c1))
		assert.Equal(t, http.StatusOK, rec1.Code)

		req2 := httptest.NewRequest(http.MethodPost, "/v1/url/short", strings.NewReader(`{"url": "http://example.com"}`))
		req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec2 := httptest.NewRecorder()
		c2 := e.NewContext(req2, rec2)

		assert.NoError(t, handler.DoShort()(c2))
		assert.Equal(t, http.StatusOK, rec2.Code)

		// Verify both responses have the same short URL
		assert.Equal(t, rec1.Body.String(), rec2.Body.String())
	})
}

func TestDoShort_onTableDeriven(t *testing.T) {
	e := echo.New()
	e.Validator = &validators.CustomValidator{Validator: validator.New()}
	handler := handlers.NewShortnerHandler()

	t.Run("Table-driven tests for multiple URLs", func(t *testing.T) {
		tests := []struct {
			name       string
			url        string
			statusCode int
		}{
			{"Valid URL 1", "http://example1.com", http.StatusOK},
			{"Valid URL 2", "http://example2.com", http.StatusOK},
			{"Valid URL 3", "http://example3.com", http.StatusOK},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/v1/url/short", strings.NewReader(`{"url": "`+tt.url+`"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				if assert.NoError(t, handler.DoShort()(c)) {
					assert.Equal(t, tt.statusCode, rec.Code)
				}
			})
		}
	})
}
