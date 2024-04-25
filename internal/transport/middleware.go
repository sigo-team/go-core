package transport

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt"
	"sigo/internal/config"
	"sigo/internal/lib"
	"sigo/internal/models"
	"sigo/internal/services"
	"strconv"
	"time"
)

const CookieName = "SESSION_ID"
const UserIDKey = "user-id-key"

func AuthMiddleware(userService *services.UserService, cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var (
			user   *models.User
			cookie string
		)
		defer func() {
			if err != nil {
				user, err = userService.CreateUser()
				if err != nil {
					return
				}
				cookie, err = lib.NewToken(jwt.StandardClaims{
					IssuedAt:  time.Now().Unix(),
					ExpiresAt: time.Now().Add(cfg.JWTMaxAge).Unix(),
					Subject:   strconv.FormatInt(user.Id(), 10),
				}, cfg.JWTSecret)
				if err != nil {
					return
				}
				c.Cookie(&fiber.Cookie{
					Name:  CookieName,
					Value: cookie,
					//SameSite: fiber.CookieSameSiteNoneMode,
					//Secure:   true,
					//HTTPOnly: true,
				})
			}
			c.Locals(UserIDKey, user)
			err = c.Next()
		}()

		cookie = c.Cookies(CookieName)
		session, err := lib.ParseToken(cookie, cfg.JWTSecret)
		if err != nil {
			return
		}
		userId, err := strconv.ParseInt(session.Subject, 10, 64)
		if err != nil {
			return
		}
		user, err = userService.ReadUser(userId)
		return
	}
}

func UpgradeMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
