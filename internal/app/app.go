package app

import (
	"main/internal/config"
	"main/internal/database"
	"main/internal/server"
)

type App struct {
	server *server.Server
	config *config.Config
}

func New(config *config.Config, db_url string) *App {
	db := database.GetDB(db_url)
	srv := server.NewServer(db)
	return &App{
		server: srv,
		config: config,
	}
}
func (a *App) Run() error {
	a.server.Start(a.config.URL + ":" + a.config.Port)
	return nil
}
