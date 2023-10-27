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
	Id *uint `json:"id"`
}

func RemoveJobPost(c *fiber.Ctx) error {
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

type UpdateJobPostBody struct {
	Id *uint `json:"id"`
	db.JobPost
}

func UpdateJobPost(c *fiber.Ctx) error {
	var body UpdateJobPostBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data is incorrect",
		})
	}

	if body.Id == nil || body.PostName == "" || body.PostOwn == "" || body.Position == "" || body.StudentLevel == "" || body.Severance == "" || body.Insurance == "" || body.Housing == "" || body.HousingAllowance == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Missing required field",
		})
	}

	var post db.JobPost

	res := db.DB.
		Where("id = ?", body.Id).
		Find(&post)

	post = body.JobPost
	post.ID = *body.Id

	if res.Error == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Id has not found",
		})
	}

	db.DB.Save(&post)

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Post updated",
	})
}

type UpdateTourBody struct {
	Id *uint `json:"id"`
	db.Tour
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

func UpdateTour(c *fiber.Ctx) error {
	var body UpdateTourBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data is incorrect",
		})
	}
	if body.Id == nil || body.TourName == "" || body.Description == "" || body.PostOwn == "" || body.Company == "" || body.Theme == "" || body.Location == "" || body.Date == "" || body.Itinerary == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Missing required field",
		})
	}

	var tour db.Tour

	res := db.DB.
		Where("id = ?", body.Id).
		Find(&tour)

	tour = body.Tour
	tour.ID = *body.Id

	if res.Error == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Id has not found",
		})
	}

	db.DB.Save(&tour)

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Post updated",
	})
}

func RemoveTour(c *fiber.Ctx) error {
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

type UpdatePartyAndEventsBody struct {
	Id *uint `json:"id"`
	db.PartyAndEvents
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

func UpdatePartyAndEvents(c *fiber.Ctx) error {
	var body UpdatePartyAndEventsBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Data is incorrect",
		})
	}

	if body.PartyName == "" || body.Description == "" || body.PostOwn == "" || body.Company == "" || body.Theme == "" || body.Location == "" || body.Date == "" || body.Itinerary == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Missing required field",
		})
	}

	var partyAndEvents db.PartyAndEvents

	res := db.DB.
		Where("id = ?", body.Id).
		Find(&partyAndEvents)

	partyAndEvents = body.PartyAndEvents
	partyAndEvents.ID = *body.Id

	if res.Error == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Id has not found",
		})
	}

	db.DB.Save(&partyAndEvents)

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Post updated",
	})
}

func RemovePartyAndEvents(c *fiber.Ctx) error {
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
