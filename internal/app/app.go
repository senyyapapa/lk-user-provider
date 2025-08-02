package app

import (
	"context"
	"errors"
	"log/slog"
	"main/internal/config"
	"main/internal/database"
	"main/internal/database/sql"
	"main/internal/server"
	sl "main/libs/logger"
	"sync"
	"sync/atomic"
)

type App struct {
	log         *slog.Logger
	server      *server.Server
	cfg         *config.Config
	userStorage database.UserStorage
}

func New(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, log *slog.Logger, isShutDown *atomic.Bool) (*App, error) {
	dialector, dbConfig := sql.InitGorm(cfg, log)
	userStorage, err := sql.NewSQLStorage(dialector, dbConfig)
	if err != nil {
		log.Error("Failed connect to database or migration", sl.Err(err))
		return nil, err
	}

	srv := server.NewServer(ctx, log, isShutDown, userStorage)

	return &App{
		log:         log,
		server:      srv,
		cfg:         cfg,
		userStorage: userStorage,
	}, nil
}

func (a *App) Run() error {
	a.log.Info("Запуск HTTP сервера по адресу: '" + a.cfg.URL + ":" + a.cfg.Port + "'...")
	return a.server.Start(a.cfg.Env, a.cfg.URL+":"+a.cfg.Port)
}

func (a *App) ShutDown(shutDownCtx context.Context) error {
	if a == nil {
		return errors.New("app is nil")
	}

	err := errors.Join(
		a.server.ShutDown(shutDownCtx),
	)

	return err
}
