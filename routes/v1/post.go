package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"lemonaid-backend/customutils"
	"lemonaid-backend/db"
	"path/filepath"
	"strings"
)

const (
	JobPost        = "JOB_POST"
	Tour           = "TOUR"
	PartyAndEvents = "PARTY_AND_EVENTS"
)

func GetJobPosts(c *fiber.Ctx) error {
	var posts []db.JobPost
	db.DB.Find(&posts)

	return c.JSON(fiber.Map{
		"status": 200,
		"data":   posts,
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

	if body.Academy == "" || body.Campus == "" || body.Category == "" || body.Position == "" || body.StudentLevel == "" || body.Severance == "" || body.Insurance == "" || body.Housing == "" || body.HousingAllowance == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Missing required field",
		})
	}

	_body := db.PendingJobPost{}
	_body.JobPost = body

	db.DB.Create(&_body)

	return c.JSON(fiber.Map{
		"id": _body.ID,
		"status":  fiber.StatusOK,
		"message": "A pending post has been created. Please wait for administrator to confirm.",
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
			"message": "Cannot parse body",
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
			"message": "Cannot parse body",
		})
	}

	if body.Id == nil || body.Academy == "" || body.Campus == "" || body.Category == "" || body.Position == "" || body.StudentLevel == "" || body.Severance == "" || body.Insurance == "" || body.Housing == "" || body.HousingAllowance == "" {
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

	go db.DB.Save(&post)

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Post updated",
	})
}

func GetPendingJobPosts(c *fiber.Ctx) error {
	var queue []db.PendingJobPost

	db.DB.Find(&queue)

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data":   queue,
	})
}

type PendingJobPostBody struct {
	Id []uint `json:"id"`
}

func AcceptPendingJobPost(c *fiber.Ctx) error {
	var body PendingJobPostBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Cannot parse body",
		})
	}

	var columns []db.PendingJobPost
	db.DB.Where("id in (?)", body.Id).
		Find(&columns)

	for _, value := range columns {
		db.DB.Create(&value.JobPost)
	}

	go db.DB.Unscoped().Delete(&columns)

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Successfully accept job post",
	})
}

func DenyPendingJobPost(c *fiber.Ctx) error {
	var body PendingJobPostBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Cannot parse body",
		})
	}

	db.DB.Unscoped().Where("id in (?)", body.Id).
		Delete(&[]db.PendingJobPost{})

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Successfully deny job post",
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
		"data":   tours,
	})
}

func WriteTour(c *fiber.Ctx) error {
	var body db.Tour
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Cannot parse body",
		})
	}

	if body.TourName == "" || body.Description == "" || body.PostOwn == "" || body.Company == "" || body.Theme == "" || body.Location == "" || body.Date == "" || body.Itinerary == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Missing required field",
		})
	}

	db.DB.Create(&body)

	return c.JSON(fiber.Map{
		"id": body.ID,
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
			"message": "Cannot parse body",
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
		"data":   partyAndEvents,
	})
}

func WritePartyAndEvents(c *fiber.Ctx) error {
	var body db.PartyAndEvents
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Cannot parse body",
		})
	}

	if body.PartyName == "" || body.Description == "" || body.PostOwn == "" || body.Company == "" || body.Theme == "" || body.Location == "" || body.Date == "" || body.Itinerary == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Missing required field",
		})
	}

	db.DB.Create(&body)

	return c.JSON(fiber.Map{
		"id": body.ID,
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
			"message": "Cannot parse body",
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

func UploadImageToJobPost(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest,
			"It's incorrect request. "+
				"multipart form must be provided")
	}

	id := form.Value["id"]
	images := form.File["images"]

	if id == nil || len(images) > 4 || len(images) == 0 {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Data is incorrect",
			})
	}

	var data db.JobPost

	fileNames := make([]string, 4)

	for _, value := range images {
		//os.MkdirAll("./public/contents", 0777)
		fileName := uuid.New().String() + filepath.Ext(value.Filename)
		fileNames = append(fileNames, "./contents/"+fileName)

		go func() {
			file, _ := value.Open()
			defer file.Close()

			buffer, err := io.ReadAll(file)

			if err != nil {
				fmt.Println("Error occurs while image writing..")
				return
			}

			customutils.ImageProcessing(buffer, 70, fileName)
		}()
	}

	result := db.DB.Model(&data).
		Where("id = ?", id).
		Update("images", strings.Join(fileNames, ","))

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{
				"status":  fiber.StatusNotAcceptable,
				"message": "Cannot find post id",
			})
	}

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Successfully upload images",
	})
}

func UploadImageToPost(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest,
			"It's incorrect request. "+
				"multipart form must be provided")
	}

	checkArr := []string{
		JobPost,
		Tour,
		PartyAndEvents,
	}

	id := form.Value["id"]
	postType := form.Value["post_type"]

	images := form.File["images"]

	if id == nil || postType == nil || len(images) > 4 || len(images) == 0 {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Data is incorrect",
			})
	}

	flag := false
	for _, value := range checkArr {
		if value == postType[0] {
			flag = true
		}
	}

	if !flag {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Data is incorrect",
			})
	}

	var data interface{}
	if postType[0] == JobPost {
		data = db.JobPost{}
	} else if postType[0] == Tour {
		data = db.Tour{}
	} else if postType[0] == PartyAndEvents {
		data = db.PartyAndEvents{}
	}

	fileNames := make([]string, 4)

	for _, value := range images {
		//os.MkdirAll("./public/contents", 0777)
		fileName := uuid.New().String() + filepath.Ext(value.Filename)
		fileNames = append(fileNames, "./contents/"+fileName)

		go func() {
			file, _ := value.Open()
			defer file.Close()

			buffer, err := io.ReadAll(file)

			if err != nil {
				fmt.Println("Error occurs while image writing..")
				return
			}

			customutils.ImageProcessing(buffer, 70, fileName)
		}()
	}

	result := db.DB.Model(&data).
		Where("id = ?", id).
		Update("images", strings.Join(fileNames, ","))

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{
				"status":  fiber.StatusNotAcceptable,
				"message": "Cannot find post id",
			})
	}

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Successfully upload images",
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

func ApplyJobPostApprovalQueue(c *fiber.Ctx) error {
	var queues []db.ApplyJobPost

	db.DB.Find(&queues)

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data":   queues,
	})
}

type ApplyJobPostAcceptDenyBody struct {
	Id []uint `json:"id"`
}

func AcceptApplyJobPost(c *fiber.Ctx) error {
	var body ApplyJobPostAcceptDenyBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Cannot parse body",
			})
	}

	var queue []db.ApplyJobPost

	result := db.DB.
		Where("id in (?)", body.Id).
		Find(&queue)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{
				"status":  fiber.StatusNotAcceptable,
				"message": "Queue item not found",
			})
	}

	for _, value := range queue {
		db.DB.Create(&value.JobPost)
	}

	go db.DB.Unscoped().Delete(&queue)

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Successfully approval queue",
	})
}

func DenyApplyJobPost(c *fiber.Ctx) error {
	var body ApplyJobPostAcceptDenyBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Cannot parse body",
			})
	}

	result := db.DB.
		Unscoped().
		Where("id in (?)", body.Id).
		Delete(&db.ApplyJobPost{})

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{
				"status":  fiber.StatusNotAcceptable,
				"message": "Queue item not found",
			})
	}

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Successfully deny user",
	})
}
