package handlers

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ShortnerHandler struct {
}

func NewShortnerHandler() *ShortnerHandler {
	return &ShortnerHandler{}
}

type URLRequest struct {
	URL string `json:"url" validate:"required"`
}

type URLResponse struct {
	URL      string `json:"url" validate:"required"`
	ShortURL string `json:"shortUrl" validate:"required"`
}

var urls = make(map[string]string)

func (handler *ShortnerHandler) DoShort() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(URLRequest)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
		}

		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		response := &URLResponse{
			URL: req.URL,
		}

		shortUrl, found := urls[req.URL]
		if found {
			response.ShortURL = shortUrl
		} else {
			response.ShortURL = "https://my.us/" + generateShortKey()
			urls[req.URL] = response.ShortURL
		}

		return c.JSON(http.StatusOK, response)
	}
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[random.Intn(len(charset))]
	}
	return string(shortKey)
}
