package main

import (
	"avito-pvz/internal/config"
	pvzRepo "avito-pvz/internal/repository/pvz"
	pvzService "avito-pvz/internal/service/pvz"
	"avito-pvz/internal/transport/http/handlers/auth"
	pvzHandler "avito-pvz/internal/transport/http/handlers/pvz"
	middlewares "avito-pvz/internal/transport/middleware"
	"avito-pvz/pkg/logger"
	"avito-pvz/pkg/postgres"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	log := logger.GetLoggerFromCtx(ctx)

	cfg, err := config.New()
	if err != nil {
		log.Fatal(ctx, "unable to load config", zap.Error(err))
		return
	}

	db, err := postgres.New(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "unable to connect db", zap.Error(err))
		return
	}

	log.Info(ctx, "Successful start!")

	repository := pvzRepo.NewRepository(db)
	service := pvzService.NewService(repository)
	handler := pvzHandler.New(service)

	authHandler := auth.NewHandler()

	r := chi.NewRouter()

	r.With(middlewares.AuthMiddleware).Post("/pvz", handler.CreatePVZ)

	r.Post("/dummyLogin", authHandler.DummyLogin)
	http.ListenAndServe(":8080", r)
}
