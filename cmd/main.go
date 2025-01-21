package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortner-module/internal/http"
)

func main() {
	server := http.NewServer()

	server.Serve()

	go func() {
		if err := server.Start(":8080"); err != nil {
			log.Printf("server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Error during server shutdown:", err)
	}

	log.Println("Server gracefully stopped")
}
