package api

import (
	"context"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"quiz/internal/clients"
	"quiz/internal/handlers"
	"quiz/internal/middleware"
	"quiz/internal/storage"
	"syscall"
	"time"
)

type ApiServer struct {
	addr        string
	storage     storage.Store
	logger      *zap.Logger
	authClient  *clients.AuthClient
	statsClient *clients.StatsClient
}

func NewApiServer(addr string, store storage.Store, logger *zap.Logger, authClient *clients.AuthClient, statsClient *clients.StatsClient) *ApiServer {
	return &ApiServer{
		addr:        addr,
		storage:     store,
		logger:      logger,
		authClient:  authClient,
		statsClient: statsClient,
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
	mux.HandleFunc("GET /quiz/{quizSessionId}/nextQuestion", middleware.VerifyToken(handlers.NewGetNextQuestionHandler(a.storage, a.logger).Handle, a.authClient))
	mux.HandleFunc("GET /quiz/{quizSessionId}/question/{id}", middleware.VerifyToken(handlers.NewGetQuestionHandler(a.storage, a.logger).Handle, a.authClient))
	mux.HandleFunc("POST /quiz/new", middleware.VerifyToken(handlers.NewStartQuizHandler(a.storage, a.logger).Handle, a.authClient))
	mux.HandleFunc("POST /quiz/{quizSessionId}/answer", middleware.VerifyToken(handlers.NewSubmitAnswerHandler(a.storage, a.logger, a.statsClient).Handle, a.authClient))
	mux.HandleFunc("POST /quiz/{quizSessionId}/finish", middleware.VerifyToken(handlers.NewFinishQuizHandler(a.storage, a.logger).Handle, a.authClient))
	mux.HandleFunc("GET /quiz/sessions", middleware.VerifyToken(handlers.NewGetUserSessionsHandler(a.storage, a.logger).Handle, a.authClient))
}
