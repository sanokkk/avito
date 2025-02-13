package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt"
	"github.com/sanokkk/avito-shop/pkg/tokens"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokensFromHeader := c.GetReqHeaders()["Authorization"]
	if tokensFromHeader == nil || len(tokensFromHeader) == 0 || tokensFromHeader[0] == "" {
		return c.SendStatus(401)
	}

	jwtVerified, err := tokens.VerifyToken(tokensFromHeader[0])
	if err != nil {
		return c.SendStatus(401)
	}

	claims := jwtVerified.Claims.(jwt.MapClaims)

	exp := claims["exp"].(float64)

	if float64(time.Now().Unix()) > exp {
		log.Error("Токен просрочен")
		return c.SendStatus(401)
	}

	c.Locals("uid", claims["uid"].(string))
	c.Locals("username", claims["username"].(string))

	return c.Next()
}
