package service

import (
	"context"
	"crm-project/internal/config"
	"crm-project/internal/repository/postgres"
	"log/slog"
)

type TempService struct {
	userRepo *postgres.UserRepo
	cfg      *config.Config
	logger   *slog.Logger
}

func NewTempService(userRepo *postgres.UserRepo, cfg *config.Config, logger *slog.Logger) *TempService {
	return &TempService{
		userRepo: userRepo,
		cfg:      cfg,
		logger:   logger,
	}
}

func (s *TempService) GrantReceptionRole(ctx context.Context, username string) error {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return err
	}

	user.RoleID = s.cfg.Roles.ReceptionID

	return s.userRepo.Update(ctx, *user)
}
