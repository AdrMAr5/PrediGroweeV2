package api

import (
	"context"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"images/internal/clients"
	"images/internal/handlers"
	"images/internal/middleware"
	"images/internal/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ApiServer struct {
	addr       string
	storage    storage.Store
	logger     *zap.Logger
	authClient *clients.AuthClient
}

func NewApiServer(addr string, store storage.Store, logger *zap.Logger, authClient *clients.AuthClient) *ApiServer {
	return &ApiServer{
		addr:       addr,
		storage:    store,
		logger:     logger,
		authClient: authClient,
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
		Addr:         a.addr,
		Handler:      middleware.VerifyToken(corsMiddleware.Handler(mux).ServeHTTP, a.authClient),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	// Start server
	go func() {
		a.logger.Info("Starting server on " + a.addr)
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
func (a *ApiServer) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /images/{id}", handlers.NewGetImageHandler(a.storage, a.logger).Handle)
	mux.HandleFunc("POST /images", handlers.NewCreateImageHandler(a.storage, a.logger).Handle)
	mux.HandleFunc("PATCH /images/{id}", handlers.NewUpdateImageHandler(a.storage, a.logger).Handle)
}
