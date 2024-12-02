package controller

import (
	"context"
	"log/slog"
	"main/internal/models"
	"main/internal/service"
)

type controller struct {
	logger  *slog.Logger
	service service.Service
}

type Controller interface {
	Login(ctx context.Context, user models.User) (models.Response, error)
	Refresh(ctx context.Context, req models.RefreshRequest) (models.Response, error)
}

func New(service service.Service, logger *slog.Logger) Controller {
	return &controller{
		logger:  logger,
		service: service,
	}
}

func (c *controller) Login(ctx context.Context, user models.User) (models.Response, error) {
	response, err := c.service.Login(ctx, user)
	if err != nil {
		return models.Response{}, err
	}

	return response, nil
}

func (c *controller) Refresh(ctx context.Context, req models.RefreshRequest) (models.Response, error) {
	response, err := c.service.Refresh(ctx, req)
	if err != nil {
		return models.Response{}, err
	}

	return response, nil
}
