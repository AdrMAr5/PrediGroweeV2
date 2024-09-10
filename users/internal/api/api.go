package api

import (
	"PrediGroweeV2/users/internal/handlers"
	"PrediGroweeV2/users/internal/middleware"
	"PrediGroweeV2/users/internal/storage"
	"context"
	"encoding/json"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ApiServer struct {
	addr    string
	storage storage.Store
	logger  *zap.Logger
}

func NewApiServer(addr string, store storage.Store, logger *zap.Logger) *ApiServer {
	return &ApiServer{
		addr:    addr,
		storage: store,
		logger:  logger,
	}
}

func (a *ApiServer) Run() {
	mux := http.NewServeMux()
	a.registerRoutes(mux)
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow requests from this origin
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},  // Add the methods you need
		AllowedHeaders:   []string{"Authorization", "Content-Type"}, // Add the headers you need
	})
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      corsMiddleware.Handler(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	// Start server
	go func() {
		a.logger.Info("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	a.logger.Info("Server exiting")
}

func (a *ApiServer) registerRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /health", a.HealthCheckHandler)
	router.HandleFunc("POST /register", handlers.NewRegisterHandler(a.storage, a.logger).Handle)
	router.HandleFunc("POST /login", handlers.NewLoginHandler(a.storage, a.logger).Handle)
	router.HandleFunc("GET /users/{id}", middleware.ValidateSession(middleware.ValidateAccessToken(handlers.NewGetUserHandler(a.storage, a.logger).Handle, a.storage), a.storage))
	router.HandleFunc("POST /verify", middleware.ValidateSession(middleware.ValidateAccessToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).ServeHTTP, a.storage), a.storage))
	router.HandleFunc("POST /refresh", middleware.ValidateSession(middleware.ValidateAccessToken(handlers.NewRefreshTokenHandler(a.storage, a.logger).Handle, a.storage), a.storage))
	router.HandleFunc("POST /logout", middleware.ValidateSession(handlers.NewLogOutHandler(a.storage, a.logger).Handle, a.storage))
}

func (a *ApiServer) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := a.storage.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
