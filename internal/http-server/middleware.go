package http_server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os/user"
	"strconv"
)

//type AuthService interface {
//  NewSession(userID int64) (string, error)
//  GetSession(token string) (*auth.Claims, error)
//}
//
//type UserService interface {
//  CreateUser(name string) (int64, error)
//  GetUser(id int64) (*user.User, error)
//
//}

const CookieName = "SESSION_ID"
const UserIDKey = "user-id-key"

func AuthMiddleware(authService *auth.AuthService, userService *user.UserService) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		var (
			usr     user.User
			cookie  string
			session *auth.Claims
			userID  int
		)

		defer func() {
			fmt.Println(err)
			if err != nil {
				usr, err = userService.CreateUser("")
				if err != nil {
					return
				}
				cookie, err = authService.NewToken(authService.NewSession(usr.ID))
				if err != nil {
					return
				}
				c.Cookie(&fiber.Cookie{Name: CookieName, Value: cookie, SameSite: fiber.CookieSameSiteNoneMode, Secure: true, HTTPOnly: true})
				fmt.Println("Set cookie")
			}
			c.SetUserContext(context.WithValue(c.UserContext(), UserIDKey, usr.ID))
			err = c.Next()
		}()

		cookie = c.Cookies(CookieName)
		session, err = authService.ParseToken(cookie)
		if err != nil {
			return
		}
		userID, err = strconv.Atoi(session.Subject)
		if err != nil {
			return
		}
		usr, err = userService.GetUser(userID)
		return
	}
}
