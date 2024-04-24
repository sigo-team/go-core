package transport

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"sigo/internal/config"
	"sigo/internal/lib"
	"strconv"
	"time"
)

const CookieName = "SESSION_ID"
const UserIDKey = "user-id-key"

func tryToParseCookie(cookie string, cfg *config.Config) (int64, error) {
	session, err := lib.ParseToken(cookie, cfg.JWTSecret)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(session.Subject, 10, 64)
}

func AuthMiddleware(idManager *lib.IdentifierManager, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			userID int64
			cookie string
			err    error
		)
		cookie = c.Cookies(CookieName)
		userID, err = tryToParseCookie(cookie, cfg)
		if err != nil {
			userID = idManager.NewID()
			cookie, err = lib.NewToken(jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(cfg.JWTMaxAge).Unix(),
				Subject:   strconv.FormatInt(userID, 10),
			}, cfg.JWTSecret)
			if err != nil {
				return err
			}
			c.Cookie(&fiber.Cookie{
				Name:     CookieName,
				Value:    cookie,
				SameSite: fiber.CookieSameSiteNoneMode,
				Secure:   true,
				HTTPOnly: true,
			})
		}
		c.Locals(UserIDKey, userID)
		return c.Next()
	}
}
