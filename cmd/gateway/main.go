package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/lldlab-gateway/internal/config"
	"github.com/lldlab-gateway/internal/handler"
	"github.com/lldlab-gateway/internal/middleware"
	"github.com/lldlab-gateway/internal/proxy"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = fmt.Sprintf("config.%s.yaml", env)
	}

	log.Printf("Loading configuration for environment: %s from %s", env, configPath)
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	router := chi.NewRouter()
	proxyHandler := proxy.NewHandler()

	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.AuthMiddleware)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   cfg.CORS.AllowedMethods,
		AllowedHeaders:   cfg.CORS.AllowedHeaders,
		ExposedHeaders:   cfg.CORS.ExposedHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		MaxAge:           cfg.CORS.MaxAge,
	}))

	for _, route := range cfg.Routes {
		if route.Path == "/health" {
			for _, method := range route.Methods {
				router.Method(method, "/health", http.HandlerFunc(handler.HealthCheck))
			}
			continue
		}

		if route.Prefix {
			for _, method := range route.Methods {
				router.Method(method, route.Path+"*", proxyHandler.ProxyRequest(route.Target, route.Path))
			}
		} else {
			for _, method := range route.Methods {
				router.Method(method, route.Path, proxyHandler.ProxyRequest(route.Target, ""))
			}
		}
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Printf("API Gateway starting on %s:%d", cfg.Server.Host, cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server exited gracefully")
	}
}
