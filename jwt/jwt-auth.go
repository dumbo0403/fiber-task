package jwt

import (
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("trishna")

type CustomClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	ID       uint64 `json:"id"`
}

func GenerateJWT(username string, id uint64) (string, error) {

	claims := CustomClaims{
		Username: username,
		ID:       id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken() func(*fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(jwtSecret)},
	})
}

func GetClaims(c *fiber.Ctx) jwt5.MapClaims {
	user := c.Locals("user").(*jwt5.Token)
	claims := user.Claims.(jwt5.MapClaims)
	return claims
}
