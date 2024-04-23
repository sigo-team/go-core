package transport

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"sigo/internal/config"
	"sigo/internal/lib"
	"sigo/internal/services"
	"strconv"
	"time"
)

const CookieName = "SESSION_ID"
const UserIDKey = "user-id-key"

func AuthMiddleware(userService *services.UserService, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var (
			usr    services.User
			cookie string
		)

		defer func() {
			if err != nil {
				usr = userService.CreateUser("")
				cookie, err = lib.NewToken(jwt.StandardClaims{
					IssuedAt:  time.Now().Unix(),
					ExpiresAt: time.Now().Add(cfg.JWTMaxAge).Unix(),
					Subject:   strconv.FormatInt(usr.ID, 10),
				}, cfg.JWTSecret)
				if err != nil {
					return
				}
				c.Cookie(&fiber.Cookie{
					Name:     CookieName,
					Value:    cookie,
					SameSite: fiber.CookieSameSiteNoneMode,
					Secure:   true,
					HTTPOnly: true,
				})
			}
			c.SetUserContext(context.WithValue(c.UserContext(), UserIDKey, usr.ID))
			err = c.Next()
		}()

		cookie = c.Cookies(CookieName)
		session, err := lib.ParseToken(cookie, cfg.JWTSecret)
		if err != nil {
			return
		}
		userID, err := strconv.ParseInt(session.Subject, 10, 64)
		if err != nil {
			return
		}
		usr = userService.GetUser(userID)
		return
	}
}
