package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	COOKIE_NAME           = "token"
	COOKIE_IS_EMPTY_ERROR = "Поле есть, а токен пустой"
	KEY_C_SET             = "Valid"
)

func (h *Handler) authMiddleware(c *gin.Context) {
	if os.Getenv("TODO_PASSWORD") == "" {
		c.SetCookie("token", "nil", -1, "/", "localhost", false, true)
		c.Set("Valid", true)
		return
	}

	cookie, err := c.Request.Cookie(COOKIE_NAME)
	if err != nil {
		logrus.Println("Какая то ошибка c.Request.Cookie: " + err.Error())
		NewResponseError(c, 401, err.Error())
		return
	}
	if cookie.Value == "" {
		logrus.Println(COOKIE_IS_EMPTY_ERROR)
		NewResponseError(c, 401, COOKIE_IS_EMPTY_ERROR)
		return
	}

	isValid, err := h.service.Authorization.ParseToken(cookie.Value)
	if err != nil {
		logrus.Println("Какая то ошибка от h.service.Authorization.ParseToken: " + err.Error())
		NewResponseError(c, 401, err.Error())
		return
	}

	c.Set(KEY_C_SET, isValid)
}
