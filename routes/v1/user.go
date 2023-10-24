package v1

import (
	"github.com/gofiber/fiber/v2"
	"lemonaid-backend/db"
)

func Me(c *fiber.Ctx) error {
	email := c.Locals("email")

	var user db.User
	db.DB.Select("first_name, last_name, email, phone_number, birthday, gender, nationality, visa_code, occupation, image, plan, user_type").Where("email = ?", email).Find(&user)

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data":   user,
	})
}

//func Teachers(c *fiber.Ctx) error {
//	var users []db.User
//	db.DB.
//		Select("first_name, last_name, email, phone_number, birthday, gender, nationality, visa_code, occupation, image, plan, admin").
//		Where("email = ?", email).
//		Find(&users)
//
//	return c.JSON(fiber.Map{
//		"status": fiber.StatusOK,
//		"data":   user,
//	})
//}

func User(c *fiber.Ctx) error {

	// TODO

	return c.SendString("TODO")
}
