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
	//
	// user actions, external api
	mux.HandleFunc("GET /quiz/sessions", middleware.VerifyToken(handlers.NewGetUserActiveSessionsHandler(a.storage, a.logger).Handle, a.authClient))
	mux.HandleFunc("POST /quiz/sessions/new", middleware.VerifyToken(handlers.NewStartQuizHandler(a.storage, a.logger, a.statsClient).Handle, a.authClient))
	mux.HandleFunc("GET /quiz/sessions/{quizSessionId}/nextQuestion", middleware.VerifyToken(handlers.NewGetNextQuestionHandler(a.storage, a.logger).Handle, a.authClient))
	mux.HandleFunc("POST /quiz/sessions/{quizSessionId}/answer", middleware.VerifyToken(handlers.NewSubmitAnswerHandler(a.storage, a.logger, a.statsClient).Handle, a.authClient))
	mux.HandleFunc("POST /quiz/sessions/{quizSessionId}/finish", middleware.VerifyToken(handlers.NewFinishQuizHandler(a.storage, a.logger, a.statsClient).Handle, a.authClient))

	//// internal api
	apiKey := os.Getenv("INTERNAL_API_KEY")

	mux.HandleFunc("GET /quiz/summary", middleware.InternalAuth(handlers.NewSummaryHandler(a.storage, a.logger).Handle, a.logger, apiKey))
	//// Case routes
	//caseHandler := handlers.NewCaseHandler(a.storage, a.logger)
	//mux.HandleFunc("GET /quiz/cases", middleware.InternalAuth(caseHandler.GetAllCases, a.logger, apiKey))
	////todo: fix, idk why tf this is not working but it is not
	////mux.HandleFunc("GET /quiz/cases/{id}", middleware.VerifyToken(middleware.InternalAuth(caseHandler.GetCaseByID), a.authClient))
	//mux.HandleFunc("POST /quiz/cases", middleware.InternalAuth(caseHandler.CreateCase, a.logger, apiKey))
	//mux.HandleFunc("PUT /quiz/cases/{id}", middleware.InternalAuth(caseHandler.UpdateCase, a.logger, apiKey))
	//mux.HandleFunc("DELETE /quiz/cases/{id}", middleware.InternalAuth(caseHandler.DeleteCase, a.logger, apiKey))
	// Question routes
	questionHandler := handlers.NewQuestionHandler(a.storage, a.logger)
	mux.HandleFunc("GET /quiz/q/{id}", middleware.VerifyToken(questionHandler.GetQuestion, a.authClient))
	mux.HandleFunc("GET /quiz/questions/{id}", middleware.InternalAuth(questionHandler.GetQuestion, a.logger, apiKey))
	mux.HandleFunc("POST /quiz/questions", middleware.InternalAuth(questionHandler.CreateQuestion, a.logger, apiKey))
	mux.HandleFunc("PATCH /quiz/questions/{id}", middleware.InternalAuth(questionHandler.UpdateQuestion, a.logger, apiKey))
	mux.HandleFunc("DELETE /quiz/questions/{id}", middleware.InternalAuth(questionHandler.DeleteQuestion, a.logger, apiKey))
	mux.HandleFunc("GET /quiz/questions", middleware.InternalAuth(questionHandler.GetAllQuestions, a.logger, apiKey))

	// options routes
	optionsHandler := handlers.NewOptionsHandler(a.storage, a.logger)
	mux.HandleFunc("GET /quiz/options", middleware.InternalAuth(optionsHandler.GetAllOptions, a.logger, apiKey))
	mux.HandleFunc("POST /quiz/options", middleware.InternalAuth(optionsHandler.CreateOption, a.logger, apiKey))
	mux.HandleFunc("PATCH /quiz/options/{id}", middleware.InternalAuth(optionsHandler.UpdateOption, a.logger, apiKey))
	mux.HandleFunc("DELETE /quiz/options/{id}", middleware.InternalAuth(optionsHandler.DeleteOption, a.logger, apiKey))
	// Group routes
	groupHandler := handlers.NewGroupHandler(a.storage, a.logger)
	mux.HandleFunc("POST /quiz/groups", middleware.InternalAuth(groupHandler.CreateGroup, a.logger, apiKey))
	mux.HandleFunc("PUT /quiz/groups/{id}", middleware.InternalAuth(groupHandler.UpdateGroup, a.logger, apiKey))
	mux.HandleFunc("DELETE /quiz/groups/{id}", middleware.InternalAuth(groupHandler.DeleteGroup, a.logger, apiKey))
	mux.HandleFunc("GET /quiz/groups", middleware.InternalAuth(groupHandler.GetAllGroups, a.logger, apiKey))

	// Parameter routes
	parameterHandler := handlers.NewParameterHandler(a.storage, a.logger)
	mux.HandleFunc("GET /quiz/parameters", middleware.InternalAuth(parameterHandler.GetAllParameters, a.logger, apiKey))
	mux.HandleFunc("POST /quiz/parameters", middleware.InternalAuth(parameterHandler.CreateParameter, a.logger, apiKey))
	mux.HandleFunc("PATCH /quiz/parameters/{id}", middleware.InternalAuth(parameterHandler.UpdateParameter, a.logger, apiKey))
	mux.HandleFunc("DELETE /quiz/parameters/{id}", middleware.InternalAuth(parameterHandler.DeleteParameter, a.logger, apiKey))
}
