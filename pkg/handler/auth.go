package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) login(c *gin.Context) {
	h.service.Authorization.CheckAuth(c)
}
