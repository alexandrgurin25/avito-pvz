package main

import (
	"avito-pvz/internal/config"
	authRepository "avito-pvz/internal/repository/auth"
	productRepo "avito-pvz/internal/repository/product"
	pvzRepo "avito-pvz/internal/repository/pvz"
	receptionRepository "avito-pvz/internal/repository/reception"
	authService "avito-pvz/internal/service/auth"
	productService "avito-pvz/internal/service/product"
	pvzService "avito-pvz/internal/service/pvz"
	receptionService "avito-pvz/internal/service/reception"
	authHandler "avito-pvz/internal/transport/http/handlers/auth"
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

	repositoryAuth := authRepository.NewRepository(db)
	serviceAuth := authService.NewService(repositoryAuth)
	handlerAuth := authHandler.NewHandler(serviceAuth)

	r := chi.NewRouter()

	r.With(middlewares.AuthMiddleware, middlewares.RequestIDMiddleware).Post("/pvz", handlerPvz.CreatePVZ)

	r.With(middlewares.AuthMiddleware, middlewares.RequestIDMiddleware).Post("/receptions", handlerReception.CreateReception)

	r.Route("/pvz", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware, middlewares.RequestIDMiddleware)
		r.Get("/", handlerPvz.GetPVZs)
		r.Post("/", handlerPvz.CreatePVZ)
		r.Route("/{pvzId}", func(r chi.Router) {
			r.Post("/close_last_reception", handlerReception.CloseLastReception)
			r.Post("/delete_last_product", handlerProduct.DeleteLastProduct)
		})
	})

	r.With(middlewares.AuthMiddleware, middlewares.RequestIDMiddleware).Post("/products", handlerProduct.AddProduct)

	r.Post("/dummyLogin", handlerAuth.DummyLogin)
	r.Post("/login", handlerAuth.Login)
	r.Post("/register", handlerAuth.Register)
	http.ListenAndServe(":8080", r)
}
