package pay

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"lemonaid-backend/db"
	"net/url"
	"os"
)

var (
	planMapping = map[string]string{
		"1": "69000",
		"2": "199000",
		"3": "49000",
		"4": "1000000",
	}
)

func PayAppFeedback(c *fiber.Ctx) error {
	// Body에서 URL 인코딩된 폼 데이터 파싱
	body := string(c.Body())
	formValues, err := url.ParseQuery(body)

	fmt.Println(formValues.Encode())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form data")
	}

	if formValues.Get("pay_state") != "4" ||
		formValues.Get("userid") != os.Getenv("PAYAPP_VERIFY_ID") || formValues.Get("linkval") != os.Getenv("PAYAPP_LINK_VALUE") {
		return fiber.ErrNotAcceptable
	}

	value, ok := planMapping[formValues.Get("var1")]
	if !ok || value != formValues.Get("price") {
		return fiber.ErrBadRequest
	}

	phone, ok := planMapping[formValues.Get("recvphone")]
	if !ok {
		return fiber.ErrBadRequest
	}

	db.DB.Model(&db.User{}).Where("phone_number = ?", phone).
		Update("plan = ?", formValues.Get("var1"))

	return c.SendString("Success")
}
