package repository

import (
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
)

const (
	TOKEN_TTL                 = 8 * time.Hour
	RESPONSE                  = "Получили объект User со следующими данными: login: %s, password: %s"
	EMPTY_USER_PASSWORD       = "Пользователь передал пустое поле пароля"
	EMPTY_USER_PASSWORD_ERROR = "А где пароль-то?!"
	ENV_PASSWORD              = "TODO_PASSWORD"
	NO_PASSWORD_IN_ENV        = "Пароль не задан. Проверь указал ли ты TODO_PASSWORD"
	NO_PASSWORD_IN_ENV_LOGRUS = "Пароль не задан. Проверь указал ли ты TODO_PASSWORD в окружении на сверере. Пускаю без пароля. Твой токен выше =)"
	WRONG_PASSWORD            = "Неправильный пароль"
	SEARCH_COOKIE             = "Искали старые куки с токеном -- "
	TOKEN_COOKEI_DONE         = "Токен выдали, куки записали"
	SUCCESS_LOGIN             = "Похоже, что всё верно - выдаю токен =)"
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
	if os.Getenv("TODO_PASSWORD") == "" {
		c.SetCookie("token", "nil", -1, "/", "localhost", false, true)
		c.JSON(200, gin.H{})
		return
	}

	var u User
	if c.ShouldBindJSON(&u) == nil {
		logrus.Println(fmt.Sprintf(RESPONSE, u.Login, u.Password))
	}

	if u.Password == "" {
		logrus.Error(EMPTY_USER_PASSWORD)
		c.JSON(401, gin.H{"error": EMPTY_USER_PASSWORD_ERROR})
		return
	}

	passwordENV := os.Getenv(ENV_PASSWORD)

	if len(passwordENV) == 0 {
		token, err := GenerateJWT(u.Login)
		if err != nil {
			logrus.Error(err)
			c.JSON(500, gin.H{"error": err.Error()})
		}

		logrus.Error(NO_PASSWORD_IN_ENV)

		c.JSON(200, gin.H{"warning": NO_PASSWORD_IN_ENV_LOGRUS, "token": token})
		return
	}

	hashPassENV := generatePasswordHash(passwordENV)

	if u.Password != "" {
		hashPass := generatePasswordHash(u.Password)

		if hashPass != hashPassENV {
			logrus.Error(WRONG_PASSWORD)
			c.JSON(401, gin.H{"error": WRONG_PASSWORD})
		}
	}

	if generatePasswordHash(u.Password) == hashPassENV {
		token, err := GenerateJWT(u.Login)
		if err != nil {
			logrus.Error(err)
			c.JSON(500, gin.H{"error": err.Error()})
		}

		oldCookie, err := c.Cookie("token")
		if err == nil {
			deleteCookie := &http.Cookie{
				Name:   "token",
				MaxAge: -1,
			}
			http.SetCookie(c.Writer, deleteCookie)
		}

		logrus.Println(SEARCH_COOKIE + oldCookie)
		//c.SetCookie("token", token, 3600, "/", "localhost", false, true)
		c.JSON(200, gin.H{"token": token})
		logrus.Println(TOKEN_COOKEI_DONE)
	}
}

func GenerateJWT(username string) (string, error) {
	if username == "" {
		username = "default" // функционал на будующее, если будут использоваться пользователи. А пока в токене будем возвращать дефолтное имя
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &myClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TOKEN_TTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		username,
	})
	logrus.Println(SUCCESS_LOGIN)

	return token.SignedString([]byte(viper.Get("SIGN_KEY").(string)))
}

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprint("%x", hash.Sum(nil))
}
