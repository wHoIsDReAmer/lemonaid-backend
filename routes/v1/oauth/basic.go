package oauth

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"lemonaid-backend/db"
	"os"
)

func GetOAuthProcessInfo(c *fiber.Ctx) error {
	token := c.Cookies("lsession", "")

	if token == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Not processing oauth",
			})
	}

	var session db.Session
	if rst := db.DB.Select("user_id, email").Where("uuid = ? and o_authing = 1", token).Find(&session); rst.RowsAffected == 0 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	return c.JSON(fiber.Map{
		"email": session.Email,
	})
}

func OAuthSetting() {
	googleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URI"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}
