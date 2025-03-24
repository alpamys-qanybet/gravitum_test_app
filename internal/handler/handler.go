package handler

import (
	"gravitum-test-app/config"
	"gravitum-test-app/internal/handler/user"
	"gravitum-test-app/internal/service"
	"gravitum-test-app/pkg/logger"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetList(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
}

type Handler struct {
	User UserHandler
}

func NewHandler(
	cfg *config.Config,
	services *service.Service,
	log *logger.Logger,
) *Handler {
	return &Handler{
		User: user.NewHandler(cfg, services.User, log),
	}
}

var _ UserHandler = (*user.UserHandler)(nil)
