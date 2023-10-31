package v1

import (
	"github.com/gofiber/fiber/v2"
	"lemonaid-backend/db"
)

func authMiddleWare(c *fiber.Ctx) error {
	token := c.Get(fiber.HeaderAuthorization, "")

	var session db.Session
	if rst := db.DB.Select("email").Where("uuid = ?", token).Find(&session); rst.RowsAffected == 0 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	c.Locals("email", session.Email)
	return c.Next()
}

func adminMiddleWare(c *fiber.Ctx) error {
	token := c.Get(fiber.HeaderAuthorization, "")

	var session db.Session
	if rst := db.DB.Select("email").Where("uuid = ?", token).Find(&session); rst.RowsAffected == 0 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	c.Locals("email", session.Email)
	email := c.Locals("email")

	var user db.User
	db.DB.Select("user_type").Where("email = ?", email).Find(&user)

	if user.UserType != 3 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Permission denied",
		})
	}

	return c.Next()
}
