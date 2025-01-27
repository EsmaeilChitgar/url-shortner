package httpq

import (
	"url-shortner-module/internal/httpq/handlers"
	"url-shortner-module/internal/httpq/validators"

	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/julienschmidt/httprouter"
)

// ------------------------------------------------------------------------------

func checkCache(f httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		res, ok := handlers.Cache.Read(id)
		if ok {
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
			return
		}
		log.Println("From Controller")
		f(w, r, p)
	}
}

// ------------------------------------------------------------------------------

type Server struct {
	*echo.Echo
}

func NewServer() *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register the custom validator
	e.Validator = &validators.CustomValidator{Validator: validator.New()}

	return &Server{e}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func (server Server) Serve() {
	v1 := server.Group("/v1")
	auth := v1.Group("/url")

	authHandler := handlers.NewShortnerHandler()

	auth.POST("/short", authHandler.DoShort())
}

func (server Server) Serve2() {
	router := httprouter.New()
	router.GET("/index", Index)
	router.GET("/product/:id", checkCache(handlers.GetProduct))

	err := http.ListenAndServe(":8081", router)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server running on :8081")
}
