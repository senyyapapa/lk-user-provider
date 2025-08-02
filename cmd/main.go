package main

import (
	"context"
	"log/slog"
	"main/internal/app"
	"main/internal/config"
	sl "main/libs/logger"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/lmittmann/tint"
)

var isShuttingDown atomic.Bool

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg := config.MustLoad()

	log := setupLogger(cfg)
	log.Info("Start application...")

	ongoingCtx, stopOngoingGraceful := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	a, err := app.New(ongoingCtx, wg, cfg, log, &isShuttingDown)
	if err != nil {
		log.Error("Не удалось инициализировать приложение.", sl.Err(err))
		stop()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-rootCtx.Done():
			return
		default:
			err = a.Run()
			if err != nil {
				log.Error("Ошибка при запуске приложения.", sl.Err(err))
				stop()
			}
		}
	}()

	<-rootCtx.Done()
	isShuttingDown.Store(true)
	log.Info("Получен сигнал отключения, выключение...")
	if cfg.Env == "prod" {
		log.Info("Распростронение информации о завершении работы...")
		time.Sleep(cfg.Readiness.DrainDelay)
	}

	shutDownCtx, cancel := context.WithTimeout(context.Background(), cfg.Shutdown.Period)
	defer cancel()

	err = a.ShutDown(shutDownCtx)
	stopOngoingGraceful()
	if err != nil {
		log.Error("Error shutting down", sl.Err(err))
		if cfg.Env == "prod" {
			time.Sleep(cfg.Shutdown.HardPeriod)
		}
	}
	wg.Wait()
	log.Info("Server shutdown gracefully.")

}

func setupLogger(cfg *config.Config) *slog.Logger {
	var log *slog.Logger

	// If logger.level varable is not set set [slog.Level] to DEBUG for "local" and "dev" and INFO for "prod"
	if cfg.Logger.Level == nil {
		var level slog.Level
		if cfg.Env != "prod" {
			level = slog.LevelDebug.Level()
		} else {
			level = slog.LevelInfo.Level()
		}
		cfg.Logger.Level = &level
	}

	switch cfg.Env {
	case "local":
		log = slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				AddSource: cfg.Logger.ShowPathCall,
				Level:     cfg.Logger.Level,
			}),
		)
	case "dev":
		log = slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				AddSource: cfg.Logger.ShowPathCall,
				Level:     cfg.Logger.Level,
			}),
		)
	case "prod":
		log = slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				AddSource: cfg.Logger.ShowPathCall,
				Level:     cfg.Logger.Level,
			}),
		)
	}

	return log
}
