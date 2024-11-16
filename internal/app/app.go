package app

import (
	"auth/internal/pkg/handler"
	"auth/internal/pkg/middleware"
	"auth/internal/pkg/postgres"
	"auth/internal/pkg/repository"
	"auth/internal/pkg/service"
	"net/http"

	"github.com/gorilla/sessions"
)

const Secret_Key = "Secret_Key"

type App struct {
	s *service.Service
	repo *repository.Repository
	h *handler.Handler
}

func New() *App {
	cfg := postgres.Config{
		DBName: "postgres",
		Port: "5433",
		Name: "postgres",
		Password: "qwerty",
		Sslmode: "disable",
	}
	db, err := postgres.New(cfg)
	if err != nil {
		panic(err.Error())
	}
	store := sessions.NewCookieStore([]byte(Secret_Key))
	middleware := middleware.New(store)
	repository := repository.New(db)
	service := service.New(repository)
	handler := handler.New(*middleware, service)
	return &App{
		s: service,
		repo: repository,
		h: handler,
	}
	
}

func(a *App) Run(port string) error {
	return http.ListenAndServe(":" + port, a.h.InitRoutes())

}