package main

import (
	"avito-pvz/internal/config"
	productRepo "avito-pvz/internal/repository/product"
	pvzRepo "avito-pvz/internal/repository/pvz"
	receptionRepository "avito-pvz/internal/repository/reception"
	productService "avito-pvz/internal/service/product"
	pvzService "avito-pvz/internal/service/pvz"
	receptionService "avito-pvz/internal/service/reception"
	"avito-pvz/internal/transport/http/handlers/auth"
	productHandler "avito-pvz/internal/transport/http/handlers/product"
	pvzHandler "avito-pvz/internal/transport/http/handlers/pvz"
	receptionHandler "avito-pvz/internal/transport/http/handlers/reseption"
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

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Successful start!")

	repositoryPvz := pvzRepo.NewRepository(db)
	servicePvz := pvzService.NewService(repositoryPvz)
	handlerPvz := pvzHandler.New(servicePvz)

	repositoryReception := receptionRepository.NewRepository(db)
	serviceReception := receptionService.NewService(repositoryReception, repositoryPvz)
	handlerReception := receptionHandler.NewHandler(serviceReception)

	repositoryProduct := productRepo.NewRepository(db)
	serviceProduct := productService.New(repositoryProduct, repositoryPvz, repositoryReception)
	handlerProduct := productHandler.New(serviceProduct)

	authHandler := auth.NewHandler()

	r := chi.NewRouter()

	r.With(middlewares.AuthMiddleware, middlewares.RequestIDMiddleware).Post("/pvz", handlerPvz.CreatePVZ)

	r.With(middlewares.AuthMiddleware, middlewares.RequestIDMiddleware).Post("/receptions", handlerReception.CreateReception)

	r.With(middlewares.AuthMiddleware, middlewares.RequestIDMiddleware).Post("/products", handlerProduct.AddProduct)

	r.Post("/dummyLogin", authHandler.DummyLogin)
	http.ListenAndServe(":8080", r)
}
