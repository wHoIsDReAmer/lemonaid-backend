package v1

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"lemonaid-backend/db"
)

func GetJobPosts(c *fiber.Ctx) error {
	var posts []db.JobPost
	db.DB.Find(&posts)

	return c.JSON(fiber.Map{
		"status": 200,
		"arrays": posts,
	})
}

func WriteJobPost(c *fiber.Ctx) error {
	email := c.Locals("email").(string)
	var user db.User
	if rst := db.DB.Select("admin").Where("email = ?", email).Find(&user); rst.RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": fiber.StatusBadRequest,
		})
	}

	if user.Admin == 0 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Permission denied",
		})
	}

	var body db.JobPost
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data is incorrect",
		})
	}

	if body.PostName == "" || body.PostOwn == "" || body.Position == "" || body.StudentLevel == "" || body.Severance == "" || body.Insurance == "" || body.Housing == "" || body.HousingAllowance == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Missing required field",
		})
	}

	db.DB.Create(&body)

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Created post",
	})
}

type RemoveBody struct {
	Id *int `json:"id"`
}

func RemoveJobPost(c *fiber.Ctx) error {
	email := c.Locals("email")

	var user db.User
	db.DB.Select("admin").Where("email = ?", email).Find(&user)

	if user.Admin != 1 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Permission denied",
		})
	}

	var body RemoveBody
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data is not correct",
		})
	}

	if body.Id == nil {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Missing required fields",
		})
	}

	var jobPost db.JobPost
	if result := db.DB.Unscoped().Where("id = ?", body.Id).Delete(&jobPost); result.RowsAffected > 0 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Successfully deleted",
		})
	}

	return c.JSON(fiber.Map{
		"status": fiber.StatusBadRequest,
	})
}

func GetTours(c *fiber.Ctx) error {
	var tours []db.Tour
	db.DB.Find(&tours)

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"arrays": tours,
	})
}

func WriteTour(c *fiber.Ctx) error {
	email := c.Locals("email")

	var user db.User
	db.DB.Select("admin").Where("email = ?", email).Find(&user)

	if user.Admin != 1 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Permission denied",
		})
	}

	var body db.Tour
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data is not correct",
		})
	}

	if body.TourName == "" || body.Description == "" || body.PostOwn == "" || body.Company == "" || body.Theme == "" || body.Location == "" || body.Date == "" || body.Itinerary == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Missing required field",
		})
	}

	go db.DB.Create(&body)

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Successfully created",
	})
}

func RemoveTour(c *fiber.Ctx) error {
	email := c.Locals("email")

	var user db.User
	db.DB.Select("admin").Where("email = ?", email).Find(&user)

	if user.Admin != 1 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Permission denied",
		})
	}

	var body RemoveBody
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data is not correct",
		})
	}

	if body.Id == nil {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Missing required fields",
		})
	}

	var tour db.Tour
	if result := db.DB.Unscoped().Where("id = ?", body.Id).Delete(&tour); result.RowsAffected > 0 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Successfully deleted",
		})
	}

	return c.JSON(fiber.Map{
		"status": fiber.StatusBadRequest,
	})
}

func GetPartyAndEvents(c *fiber.Ctx) error {
	var partyAndEvents []db.PartyAndEvents
	db.DB.Find(&partyAndEvents)

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"arrays": partyAndEvents,
	})
}

func WritePartyAndEvents(c *fiber.Ctx) error {
	email := c.Locals("email")

	var user db.User
	db.DB.Select("admin").Where("email = ?", email).Find(&user)

	if user.Admin != 1 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Permission denied",
		})
	}

	var body db.PartyAndEvents
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data is not correct",
		})
	}

	if body.PartyName == "" || body.Description == "" || body.PostOwn == "" || body.Company == "" || body.Theme == "" || body.Location == "" || body.Date == "" || body.Itinerary == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Missing required field",
		})
	}

	go db.DB.Create(&body)

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Successfully created",
	})
}

func RemovePartyAndEvents(c *fiber.Ctx) error {
	email := c.Locals("email")

	var user db.User
	db.DB.Select("admin").Where("email = ?", email).Find(&user)

	if user.Admin != 1 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Permission denied",
		})
	}

	var body RemoveBody
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data is not correct",
		})
	}

	if body.Id == nil {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Missing required fields",
		})
	}

	var partyAndEvents db.PartyAndEvents
	if result := db.DB.Unscoped().Where("id = ?", body.Id).Delete(&partyAndEvents); result.RowsAffected > 0 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Successfully deleted",
		})
	}

	return c.JSON(fiber.Map{
		"status": fiber.StatusBadRequest,
	})
}

type ApplyJobPostBody struct {
	PostID uint `json:"post_id"`
}

func ApplyJobPost(c *fiber.Ctx) error {
	var body ApplyJobPostBody
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data is incorrect",
		})
	}

	email := c.Locals("email")

	var user db.User
	db.DB.Select("plan, id").Where("email = ?", email).Find(&user)

	if user.Plan != db.STANDARD || user.Plan != db.PREMIUM {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Permission denied",
		})
	}

	var jobPost db.JobPost
	var column db.ApplyJobPost
	db.DB.Select("id").Where("id = ?", body.PostID).Select(&jobPost)

	if jobPost.Model.ID == 0 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Invalid post id",
		})
	}

	column.JobPost = jobPost
	column.User = db.User{Model: gorm.Model{ID: user.ID}}
	db.DB.Create(&column)

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Success",
	})
}
