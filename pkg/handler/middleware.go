package handler

import (
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	AUTH_HEADER           = "Authorization"
	HEADER_IS_EMPTY_ERROR = "В хедере нет заголовка Authorization"
	HEADER_FORMAT_ERROR   = "Неверный формат хедера"
)

func (h *Handler) authMiddleware(c *gin.Context) {
	header := c.GetHeader(AUTH_HEADER)
	if header == "" {
		NewResponseError(c, 401, HEADER_IS_EMPTY_ERROR)
		return
	}

	headerStr := strings.Split(header, " ")
	if len(headerStr) != 2 {
		NewResponseError(c, 401, HEADER_FORMAT_ERROR)
		return
	}

	if headerStr[0] != "Bearer" {
		NewResponseError(c, 401, HEADER_FORMAT_ERROR)
		return
	}

	isValid, err := h.service.Authorization.ParseToken(headerStr[1])

	if err != nil {
		NewResponseError(c, 401, err.Error())
		return
	}

	c.Set("Valid", isValid)
}
