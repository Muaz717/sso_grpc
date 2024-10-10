package app

import (
	"context"
	"log/slog"

	grpcapp "github.com/Muaz717/sso/internal/app/grpc"
	"github.com/Muaz717/sso/internal/config"
	"github.com/Muaz717/sso/internal/lib/logger/sl"
	"github.com/Muaz717/sso/internal/services/auth"
	"github.com/Muaz717/sso/internal/storage/postgres"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
) *App {

	storage, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		panic(err)
	}

	service := auth.New(log, storage, storage, storage, cfg.TokenTTL)

	grpcApp := grpcapp.New(log, service, cfg.GRPC.Port)

	return &App{
		GRPCSrv: grpcApp,
	}
}
