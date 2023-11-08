package oauth

import (
	"github.com/gofiber/fiber/v2"
	"lemonaid-backend/db"
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
