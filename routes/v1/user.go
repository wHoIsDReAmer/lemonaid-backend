package v1

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gofiber/fiber/v2"
	"lemonaid-backend/customutils"
	"lemonaid-backend/db"
	"strconv"
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

func Teachers(c *fiber.Ctx) error {
	//email := c.Locals("email")

	//var user db.User
	//db.DB.Select("plan").Where("email = ?", email).Find(&user)
	//
	//if user.Plan != db.RESUME && user.Plan != db.SPECIALIST {
	//	return c.JSON(fiber.Map{
	//		"status":  fiber.StatusForbidden,
	//		"message": "Permission denied",
	//	})
	//}

	var users []db.User
	db.DB.
		Select("id, first_name, last_name, email, phone_number, birthday, gender, nationality, image").
		Where("user_type = ?", "2").
		Find(&users)

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data":   users,
	})
}

type UserBody struct {
	Id uint `json:"id"`
}

func User(c *fiber.Ctx) error {
	var body UserBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Cannot parse body",
			})
	}

	var user db.User
	db.DB.Select("id, first_name, last_name, email, phone_number, birthday, gender, nationality, visa_code, occupation, image, plan, user_type").
		Where("id = ?", body.Id).Find(&user)

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data":   user,
	})
}

type UserEditBody struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
}

func UserEdit(c *fiber.Ctx) error {
	var body UserEditBody
	if err := c.BodyParser(&body); err != nil {

		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Cannot parse body",
			})
	}

	salt := strconv.Itoa(customutils.RandI(10000, 50000))
	hasher := sha256.New()
	hasher.Write([]byte(body.Password + salt))

	rst := db.DB.Model(&db.User{}).Where("id = ?", body.Id).
		Updates(db.User{
			Password: hex.EncodeToString(hasher.Sum(nil)),
			Salt:     salt,
		})

	if rst.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{
				"status":  fiber.StatusNotFound,
				"message": "Cannot find user",
			})
	}

	return c.JSON(fiber.Map{})
}
