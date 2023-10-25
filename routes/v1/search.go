package v1

import (
	"github.com/gofiber/fiber/v2"
	"lemonaid-backend/db"
)

/*
	response like
	response

	{
		"status": 200,
		"teachers": []
		"posts": []
	}

*/

func SearchPostAndTeachers(c *fiber.Ctx) error {
	email := c.Locals("email")

	var user db.User
	db.DB.Select("plan").Where("email = ?", email).Find(&user)

	if user.Plan != db.RESUME && user.Plan != db.SPECIALIST {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Permission denied",
		})
	}

	var teachers []db.User
	var posts []db.JobPost

	db.DB.
		Select("id, first_name, last_name, email, phone_number, birthday, gender, nationality, image").
		Where("user_type = ?", "2").
		Find(&teachers)

	db.DB.Find(&posts)

	return c.JSON(fiber.Map{
		"status":   200,
		"teachers": teachers,
		"posts":    posts,
	})
}
