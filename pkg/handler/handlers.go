package handler

import (
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services service.Service
}

func NewHandler(services service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) SetupRoutes() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	router := gin.New()

	router.GET("/check", h.checkCart)

	return router
}
