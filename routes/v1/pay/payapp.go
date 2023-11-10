package pay

import (
	"fmt"
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

func PayAppFeedback(c *fiber.Ctx) error {
	// Body에서 URL 인코딩된 폼 데이터 파싱
	body := string(c.Body())
	formValues, err := url.ParseQuery(body)

	fmt.Println(formValues.Encode())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form data")
	}

	response, err := sendPayRequest(formValues)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if response["state"][0] == "1" {
		// 결제 요청 성공 처리
	} else {
		// 결제 요청 실패 처리
	}

	return c.SendString("Success")
}

func sendPayRequest(formData url.Values) (map[string][]string, error) {
	// HTTP 클라이언트 초기화
	client := &http.Client{}

	// 요청 URL 구성
	requestURL := "http://" + PAYAPP_API_DOMAIN + PAYAPP_API_URL

	// POST 요청 생성
	req, err := http.NewRequest("POST", requestURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	// 필요한 헤더 설정
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 요청 보내기
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 응답 데이터 읽기
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 응답 데이터 파싱 (URL 인코딩된 형태)
	responseValues, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	return responseValues, nil
}
