package api

import (
	"context"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"stats/internal/clients"
	"stats/internal/handlers"
	"stats/internal/middleware"
	"stats/internal/storage"
	"syscall"
	"time"
)

type ApiServer struct {
	addr       string
	authClient *clients.AuthClient
	storage    storage.Storage
	logger     *zap.Logger
}

func NewApiServer(addr string, storage storage.Storage, logger *zap.Logger, authClient *clients.AuthClient) *ApiServer {
	return &ApiServer{
		addr:       addr,
		authClient: authClient,
		storage:    storage,
		logger:     logger,
	}
}

func (a *ApiServer) Run() {
	mux := http.NewServeMux()
	a.registerRoutes(mux)
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "https://predigrowee.agh.edu.pl",
			"https://www.predigrowee.agh.edu.pl"}, // Allow requests from this origin
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},  // Add the methods you need
		AllowedHeaders:   []string{"Authorization", "Content-Type"}, // Add the headers you need
	})
	srv := &http.Server{
		Addr:         a.addr,
		Handler:      corsMiddleware.Handler(mux),
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
	mux.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
	internalApiKey := os.Getenv("INTERNAL_API_KEY")
	//internal
	mux.HandleFunc("POST /stats/saveSession", middleware.InternalAuth(handlers.NewQuizStatsHandler(a.storage, a.logger).SaveSession, a.logger, internalApiKey))
	mux.HandleFunc("POST /stats/{quizSessionId}/saveResponse", middleware.InternalAuth(handlers.NewQuizStatsHandler(a.storage, a.logger).SaveResponse, a.logger, internalApiKey))
	mux.HandleFunc("POST /stats/{quizSessionId}/finish", middleware.InternalAuth(handlers.NewQuizStatsHandler(a.storage, a.logger).FinishSession, a.logger, internalApiKey))

	//external
	mux.HandleFunc("GET /stats/userStats", middleware.VerifyToken(handlers.NewGetUserStatsHandler(a.storage, a.logger).Handle, a.authClient))
	mux.HandleFunc("GET /stats/{quizSessionId}", middleware.VerifyToken(handlers.NewQuizStatsHandler(a.storage, a.logger).GetStats, a.authClient))
	mux.HandleFunc("POST /stats/survey", middleware.VerifyToken(handlers.NewSaveSurveyHandler(a.storage, a.logger).Handle, a.authClient))
}
