package v1

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io"
	"lemonaid-backend/customutils"
	"lemonaid-backend/db"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var body LoginBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data is incorrect"})
	}

	if body.Email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing required fields"})
	}

	var user db.User
	result := db.DB.Select("email, user_accepted", "password", "salt").Where("email = ?", body.Email).Find(&user)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Email has not found",
		})
	}

	if user.UserAccepted == 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  fiber.StatusForbidden,
			"message": "Cannot access until admin accept you",
		})
	}

	salt := user.Salt
	hasher := sha256.New()
	hasher.Write([]byte(body.Password + salt))

	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	if hashedPassword != user.Password {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Password is incorrect",
		})
	}

	_uuid := uuid.New()

	// add session
	sess := new(db.Session)
	db.DB.Where("email = ?", body.Email).FirstOrInit(sess)

	sess.Uuid = _uuid.String()
	sess.Email = body.Email
	sess.Expires = time.Now().Add(time.Duration(6) * time.Hour)

	db.DB.Save(sess)

	return c.JSON(fiber.Map{
		"status":  200,
		"session": _uuid.String(),
	})
}

func Register(c *fiber.Ctx) error {
	// http.Request
	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "It's incorrect request. multipart form must be provided")
	}
	defer form.RemoveAll()

	firstName := form.Value["first_name"] // required
	lastName := form.Value["last_name"]   // required
	email := form.Value["email"]          // required
	password := form.Value["password"]    // required

	phoneNumber := form.Value["phone_number"] // required
	birthday := form.Value["birthday"]        // required
	gender := form.Value["gender"]

	nationality := form.Value["nationality"]
	visacode := form.Value["visa_code"]
	occupation := form.Value["occupation"]
	videoMessenger := form.Value["video_messenger"]
	videoMessengerId := form.Value["video_messenger_id"]

	userType := form.Value["user_type"]

	if firstName == nil || firstName[0] == "" || lastName == nil || lastName[0] == "" || email == nil || email[0] == "" || password == nil || password[0] == "" || phoneNumber == nil || phoneNumber[0] == "" || birthday == nil || birthday[0] == "" || userType == nil || userType[0] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Missing required fields",
		})
	}

	if v, err := strconv.Atoi(userType[0]); err != nil || (v != db.TEACHER && v != db.ACADEMY) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "User type is incorrect.",
		})
	}

	var checkUser db.User
	result := db.DB.Select("email", "phone_number").Where("email = ? or phone_number = ?", email[0], phoneNumber[0]).Find(&checkUser)

	if result.RowsAffected > 0 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Already exists email or phone number",
		})
	}

	if !emailValidation(email[0]) {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Email must be email form",
		})
	}

	var imagePath *string
	profile := form.File["profile_image"]
	if profile != nil {
		h := profile[0]

		filename := hex.EncodeToString([]byte(email[0]+"profile")) + filepath.Ext(h.Filename)
		path := "./contents/" + filename
		imagePath = &path

		go func() {
			if h.Size > (1024*1024)*5 {
				return
			}

			file, _ := h.Open()
			defer file.Close()

			buffer, err := io.ReadAll(file)

			if err != nil {
				fmt.Println("Error while image writing..")
			}

			customutils.ImageProcessing(buffer, 70, filename)

			//os.MkdirAll("./public/contents", 0777)
			//
			//dst, _ := os.Create("./public/contents/" + hex.EncodeToString([]byte(email[0]+"profile")) + filepath.Ext(h.Filename))
			//defer dst.Close()
		}()
	}

	var resumeBlob *[]byte
	resume := form.File["resume"]
	if resume != nil {
		h := resume[0]
		file, err := h.Open()
		defer file.Close()

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot open resume file"})
		}

		bytes := make([]byte, h.Size)
		file.Read(bytes)
		resumeBlob = &bytes
	}

	user := new(db.User)
	user.FirstName = firstName[0]
	user.LastName = lastName[0]
	user.Email = email[0]

	user.Salt = string(rune(customutils.RandI(10000, 50000)))

	hasher := sha256.New()
	hasher.Write([]byte(password[0] + user.Salt))
	user.Password = hex.EncodeToString(hasher.Sum(nil))

	user.PhoneNumber = phoneNumber[0]
	user.Image = imagePath
	user.Resume = resumeBlob

	value, err := time.Parse("2006-01-02", birthday[0])
	if err != nil {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Incorrect form of date",
		})
	}
	user.Birthday = value

	// ...Optional values
	if gender == nil {
		user.Gender = nil
	} else {
		user.Gender = &gender[0]
	}

	if nationality == nil {
		user.Nationality = nil
	} else {
		user.Nationality = &nationality[0]
	}

	if visacode == nil {
		user.VisaCode = nil
	} else {
		user.VisaCode = &visacode[0]
	}

	if occupation == nil {
		user.Occupation = nil
	} else {
		user.Occupation = &occupation[0]
	}

	if videoMessenger == nil {
		user.VideoMessenger = nil
	} else {
		user.VideoMessenger = &videoMessenger[0]
	}

	if videoMessengerId == nil {
		user.VideoMessengerID = nil
	} else {
		user.VideoMessengerID = &videoMessengerId[0]
	}

	go db.DB.Create(&user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  fiber.StatusCreated,
		"message": "Successfully account created",
	})
}

func Logout(c *fiber.Ctx) error {
	token := c.Get(fiber.HeaderAuthorization, "")

	if token == "" {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Invalid Token",
		})
	}

	sess := new(db.Session)
	db.DB.Where("uuid = ?", token).FirstOrInit(sess)
	if sess == nil {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Invalid Token",
		})
	}

	go db.DB.Unscoped().Delete(sess)

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Successfully logout",
	})
}

func UserApprovalQueue(c *fiber.Ctx) error {
	var queues []db.User

	db.DB.Where("user_accepted = 0").Find(&queues)

	return c.JSON(fiber.Map{
		"status": fiber.StatusOK,
		"data":   queues,
	})
}

type UserAcceptDenyBody struct {
	Email []string `json:"email"`
}

func AcceptUser(c *fiber.Ctx) error {
	var body UserAcceptDenyBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Cannot parse body",
			})
	}

	var user []db.User

	result := db.DB.Select("id, user_accepted").
		Where("email in (?)", body.Email).
		Find(&user)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotAcceptable).
			JSON(fiber.Map{
				"status":  fiber.StatusNotAcceptable,
				"message": "Queue item not found",
			})
	}

	go db.DB.Model(&user).
		Update("user_accepted", 1)

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Successfully approval user",
	})
}

func DenyUser(c *fiber.Ctx) error {
	var body UserAcceptDenyBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Cannot parse body",
			})
	}

	result := db.DB.Unscoped().Model(&db.User{}).Where("email in (?)").
		Delete(&db.User{})

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

func emailValidation(email string) bool {
	const emailPattern = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailPattern)

	return re.MatchString(email)
}
