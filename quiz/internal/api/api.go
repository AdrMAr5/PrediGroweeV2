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
	a.logger.Info("run server")
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
	a.logger.Info("about to start the server")
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
	a.logger.Info("registering routes")

	// user actions, external api
	mux.HandleFunc("GET /quiz/sessions", middleware.VerifyToken(handlers.NewGetUserActiveSessionsHandler(a.storage, a.logger).Handle, a.authClient))
	mux.HandleFunc("POST /quiz/new", middleware.VerifyToken(handlers.NewStartQuizHandler(a.storage, a.logger, a.statsClient).Handle, a.authClient))
	mux.HandleFunc("GET /quiz/{quizSessionId}/nextQuestion", middleware.VerifyToken(handlers.NewGetNextQuestionHandler(a.storage, a.logger).Handle, a.authClient))
	mux.HandleFunc("POST /quiz/{quizSessionId}/answer", middleware.VerifyToken(handlers.NewSubmitAnswerHandler(a.storage, a.logger, a.statsClient).Handle, a.authClient))
	mux.HandleFunc("POST /quiz/{quizSessionId}/finish", middleware.VerifyToken(handlers.NewFinishQuizHandler(a.storage, a.logger, a.statsClient).Handle, a.authClient))

	// Case routes
	caseHandler := handlers.NewCaseHandler(a.storage, a.logger)
	mux.HandleFunc("GET /quiz/cases", middleware.VerifyToken(middleware.WithAdminRole(caseHandler.GetAllCases), a.authClient))
	//todo: fix, idk why tf this is not working but it is not
	//mux.HandleFunc("GET /quiz/cases/{id}", middleware.VerifyToken(middleware.WithAdminRole(caseHandler.GetCaseByID), a.authClient))
	mux.HandleFunc("POST /quiz/cases", middleware.VerifyToken(middleware.WithAdminRole(caseHandler.CreateCase), a.authClient))
	mux.HandleFunc("PUT /quiz/cases/{id}", middleware.VerifyToken(middleware.WithAdminRole(caseHandler.UpdateCase), a.authClient))
	mux.HandleFunc("DELETE /quiz/cases/{id}", middleware.VerifyToken(middleware.WithAdminRole(caseHandler.DeleteCase), a.authClient))
	// Question routes
	questionHandler := handlers.NewQuestionHandler(a.storage, a.logger)
	mux.HandleFunc("GET /quiz/{id}", middleware.VerifyToken(questionHandler.GetQuestion, a.authClient))
	mux.HandleFunc("POST /quiz/questions", middleware.VerifyToken(middleware.WithAdminRole(questionHandler.CreateQuestion), a.authClient))
	mux.HandleFunc("PUT /quiz/questions/{id}", middleware.VerifyToken(middleware.WithAdminRole(questionHandler.UpdateQuestion), a.authClient))
	mux.HandleFunc("DELETE /quiz/questions/{id}", middleware.VerifyToken(middleware.WithAdminRole(questionHandler.DeleteQuestion), a.authClient))
	mux.HandleFunc("GET /quiz/questions", middleware.VerifyToken(middleware.WithAdminRole(questionHandler.GetAllQuestions), a.authClient))
	// Group routes
	groupHandler := handlers.NewGroupHandler(a.storage, a.logger)
	mux.HandleFunc("POST /quiz/groups", middleware.VerifyToken(middleware.WithAdminRole(groupHandler.CreateGroup), a.authClient))
	mux.HandleFunc("PUT /quiz/groups/{id}", middleware.VerifyToken(middleware.WithAdminRole(groupHandler.UpdateGroup), a.authClient))
	mux.HandleFunc("DELETE /quiz/groups/{id}", middleware.VerifyToken(middleware.WithAdminRole(groupHandler.DeleteGroup), a.authClient))
	mux.HandleFunc("GET /quiz/groups", middleware.VerifyToken(middleware.WithAdminRole(groupHandler.GetAllGroups), a.authClient))

	// Parameter routes
	parameterHandler := handlers.NewParameterHandler(a.storage, a.logger)
	mux.HandleFunc("GET /quiz/parameters", middleware.VerifyToken(middleware.WithAdminRole(parameterHandler.GetAllParameters), a.authClient))
	mux.HandleFunc("POST /quiz/parameters", middleware.VerifyToken(middleware.WithAdminRole(parameterHandler.CreateParameter), a.authClient))
	mux.HandleFunc("PUT /quiz/parameters/{id}", middleware.VerifyToken(middleware.WithAdminRole(parameterHandler.UpdateParameter), a.authClient))
	mux.HandleFunc("DELETE /quiz/parameters/{id}", middleware.VerifyToken(middleware.WithAdminRole(parameterHandler.DeleteParameter), a.authClient))
}
