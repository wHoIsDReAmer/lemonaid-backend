package v1

import (
	"lemonaid-backend/db"
	"strings"

	"github.com/gofiber/fiber/v2"
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
	//email := c.Locals("email")
	name := c.Query("value")

	//var user db.User
	//db.DB.Select("plan").Where("email = ?", email).Find(&user)
	//
	//if user.Plan != db.RESUME && user.Plan != db.SPECIALIST {
	//	return c.JSON(fiber.Map{
	//		"status":  fiber.StatusForbidden,
	//		"message": "Permission denied",
	//	})
	//}

	if len(strings.Replace(name, " ", "", -1)) == 0 {
		return c.JSON(fiber.Map{
			"status":   200,
			"teachers": make([]interface{}, 0),
			"posts":    make([]interface{}, 0),
		})
	}

	var teachers []db.User
	var posts []db.JobPost

	db.DB.
		Select("id, first_name, last_name, email, phone_number, birthday, gender, nationality, image").
		Where("user_type = ? and (first_name like ? or last_name like ?)", "2", "%"+name+"%", "%"+name+"%").
		Find(&teachers)

	db.DB.
		Where("academy like ? or campus like ?", "%"+name+"%", "%"+name+"%").
		Find(&posts)

	return c.JSON(fiber.Map{
		"status":   200,
		"teachers": teachers,
		"posts":    posts,
	})
}
