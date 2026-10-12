package main

import (
	_ "net/http/pprof"
	"os"
	"log/slog"
	"log"
	"net/http"
	"time"
	"context"

	"github.com/labstack/echo/v5"

	"github.com/redis/go-redis/v9"

	"github.com/nekto-sns/nekto/app/route"
	"github.com/nekto-sns/nekto/app/handler"
	"github.com/nekto-sns/nekto/app/middleware"
	"github.com/nekto-sns/nekto/app/repository"
	"github.com/nekto-sns/nekto/app/service"

	"github.com/nekto-sns/nekto/app/client/scratchauth"

	"github.com/nekto-sns/nekto/app/sessionmanager"

	"github.com/nekto-sns/nekto/app/shared/config"
	"github.com/nekto-sns/nekto/app/shared/database"
	"github.com/nekto-sns/nekto/app/shared/errorhandler"
)

func main() {
	cfg := config.Load()

	logOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: false,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, logOpts))
	slog.SetDefault(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	dbPool, err := database.NewPool(ctx, cfg.DBUrl)
	defer dbPool.Close()
	if err != nil {
		slog.Error("Server startup failed", "error", err)
		os.Exit(1)
	}

	err = database.RunMigrations(ctx, dbPool)
	if err != nil {
		slog.Error("DB migration failed", "error", err)
		os.Exit(1)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
		PoolSize: 512,
	})

	userRepo    := repository.NewUser(dbPool)
	userSvc     := service.NewUser(userRepo)
	userHandler := handler.NewUser(userSvc)

	sessionRepo := repository.NewSession(rdb)
	sessionMN   := sessionmanager.New(24 * 7 * time.Hour)
	sessionSvc  := service.NewSession(sessionRepo, sessionMN)

	saClient    := scratchauth.New(&http.Client{}, cfg.ScratchAuthURL, []string{cfg.LoginCallbackURL})
	authRepo    := repository.NewScratchAuth(dbPool)
	authSvc     := service.NewScratchAuth(sessionRepo, authRepo, saClient, sessionMN)
	authHandler := handler.NewScratchAuth(authSvc, sessionMN, cfg.ScratchAuthURL, cfg.LoginCallbackURL)
	authMW      := middleware.NewAuth(sessionMN, sessionSvc)

	e := echo.New()
	e.HTTPErrorHandler = errorhandler.ErrorHandler

	route.Setup(e, userHandler, authHandler, authMW)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	e.Start(cfg.Port)
}
