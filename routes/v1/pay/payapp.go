package pay

import (
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	PAYAPP_API_DOMAIN = "api.payapp.kr"
	PAYAPP_API_URL    = "/oapi/apiLoad.html"
)

func PayApp(c *fiber.Ctx) error {
	// Parse the form data
	if err := c.Request().ParseForm(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form data")
	}

	postData := make(map[string]string)
	for key, values := range c.Request().PostForm {
		postData[key] = values[0]
	}

	response, err := payappOapiPost(postData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString(response)
}

func payappOapiPost(postData string) (string, error) {
	formData := url.Values{}
	for _, v := range strings.Split(postData, "&") {
		pair := strings.SplitN(v, "=", 2)
		if len(pair) == 2 {
			formData.Set(pair[0], pair[1])
		}
	}

	resp, err := http.PostForm("http://"+PAYAPP_API_DOMAIN+PAYAPP_API_URL, formData)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
