package pay

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"os"
)

const (
	PAYAPP_API_DOMAIN = "api.payapp.kr"
	PAYAPP_API_URL    = "/oapi/apiLoad.html"
)

func PayAppFeedback(c *fiber.Ctx) error {
	// Body에서 URL 인코딩된 폼 데이터 파싱
	body := string(c.Body())
	formValues, err := url.ParseQuery(body)

	fmt.Println(formValues.Encode())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form data")
	}

	if formValues.Get("linkval") != os.Getenv("PAYAPP_LINK_VALUE") {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form data")
	}

	return c.SendString("Success")
}
