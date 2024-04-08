package handler

import (
	"github.com/gin-gonic/gin"
)

type Error struct {
	Message string `json:"error"`
}

func NewResponseError(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, Error{Message: message})
}
