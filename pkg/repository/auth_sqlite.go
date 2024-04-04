package repository

import (
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	TOKEN_TTL = 8 * time.Hour
)

type AuthSqlite struct {
	db *sqlx.DB
}

type User struct {
	Login    string `json:"-"` //не используем. но, может когда-то пригодиться? =) не использовать же ЦЕЛУЮ структуру из-за одного только поля?!
	Password string `json:"password"`
}

type myClaims struct {
	jwt.StandardClaims
	Login string `json:"login"`
}

func NewAuthSqlite(db *sqlx.DB) *AuthSqlite {
	return &AuthSqlite{db: db}
}

func (a *AuthSqlite) CheckAuth(c *gin.Context) {
	var u User
	if c.ShouldBindJSON(&u) == nil {
		logrus.Println(fmt.Sprintf(
			"Получили объект User со следующими данными: login: %s, password: %s",
			u.Login, u.Password))
	}

	if u.Password == "" {
		logrus.Error("Пользователь передал пустое поле пароля")
		c.JSON(401, gin.H{"error": "А где пароль-то?!"})
		return
	}

	passwordENV := os.Getenv("TODO_PASSWORD")

	if len(passwordENV) == 0 {
		token, err := GenerateJWT(u.Login)
		if err != nil {
			logrus.Error(err)
			c.JSON(500, gin.H{"error": err.Error()})
		}

		logrus.Error("Пароль не задан. Проверь указал ли ты TODO_PASSWORD")
		c.JSON(200, gin.H{"warning": "Пароль не задан. Проверь указал ли ты TODO_PASSWORD в окружении на сверере. Пускаю без пароля. Твой токен выше =)", "token": token})
		return
	}

	hashPassENV := generatePasswordHash(passwordENV)

	if u.Password != "" {
		hashPass := generatePasswordHash(u.Password)

		if hashPass != hashPassENV {
			logrus.Error("Неверный пароль")
			c.JSON(401, gin.H{"error": "Неверный пароль"})
		}
	}

	if generatePasswordHash(u.Password) == hashPassENV {
		token, err := GenerateJWT(u.Login)
		if err != nil {
			logrus.Error(err)
			c.JSON(500, gin.H{"error": err.Error()})
		}
		c.JSON(200, gin.H{"token": token})
	}
}

func GenerateJWT(username string) (string, error) {
	if username == "" {
		username = "default" // функционал на будующее, если будут исопльзоваться пользователи. А пока в токене будем возвращать дефолтное имя
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &myClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TOKEN_TTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		username,
	})

	return token.SignedString([]byte(viper.Get("SIGN_KEY").(string)))
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprint("%x", hash.Sum(nil))
}
