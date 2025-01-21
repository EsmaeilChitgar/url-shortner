package http

import (
	"url-shortner-module/internal/http/handlers"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CustomValidator wraps the validator instance
type CustomValidator struct {
	validator *validator.Validate
}

// Validate implements the Validator interface for Echo
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type Server struct {
	*echo.Echo
}

func NewServer() *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register the custom validator
	e.Validator = &CustomValidator{validator: validator.New()}

	return &Server{e}
}

func (server Server) Serve() {
	v1 := server.Group("/v1")
	auth := v1.Group("/url")

	authHandler := handlers.NewShortnerHandler()

	auth.POST("/short", authHandler.DoShort())
}
