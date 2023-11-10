package oauth

import (
	"crypto/rand"
	"encoding/base64"
	json2 "encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io/ioutil"
	"lemonaid-backend/db"
	"net/http"
	"os"
	"time"
)

func NaverLogin(c *fiber.Ctx) error {
	clientId := os.Getenv("NAVER_OAUTH_CID")
	redirectURI := os.Getenv("NAVER_OAUTH_REDIRECT_URI")

	sess := store.Get(c)

	// 상태 문자열 생성
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	state := base64.URLEncoding.EncodeToString(b)

	// 세션에 상태 저장
	sess.Set("state", state)
	err = sess.Save()
	if err != nil {
		return err
	}

	// 로그인 페이지로 리디렉션
	return c.Redirect("https://nid.naver.com/oauth2.0/authorize?response_type=code&client_id=" +
		clientId + "&redirect_uri=" + redirectURI + "&state=" + state)
}

type NaverToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string `json:"expires_in"`
}

type NaverOAuthInfo struct {
	ResultCode string        `json:"resultcode"`
	Message    string        `json:"message"`
	Response   InnerResponse `json:"response"` // 내부 구조체를 참조합니다.
}

type InnerResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func NaverCallback(c *fiber.Ctx) error {
	sess := store.Get(c)

	state := c.Query("state")

	if sess.Get("state") != state {
		return fiber.ErrUnauthorized
	}

	code := c.Query("code")

	clientId := os.Getenv("NAVER_OAUTH_CID")
	clientSecret := os.Getenv("NAVER_OAUTH_SECRET")

	reqeust, err := http.NewRequest("GET",
		"https://nid.naver.com/oauth2.0/token?grant_type=authorization_code&client_id="+
			clientId+"&client_secret="+clientSecret+"&code="+code, nil)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Error occurs while logining oauth..",
		})
	}

	client := http.Client{}
	resp, err := client.Do(reqeust)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Error occurs while logining oauth..",
		})
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Error occurs while logining oauth..",
		})
	}

	var data NaverToken
	err = json2.Unmarshal(body, &data)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return NaverAuthProcessing(c, data)
}

func NaverAuthProcessing(c *fiber.Ctx, data NaverToken) error {
	request, err := http.NewRequest("GET", "https://openapi.naver.com/v1/nid/me", nil)
	if err != nil {
		return c.Send([]byte("Maybe error occurs while logining.."))
	}

	request.Header.Set("Authorization", data.TokenType+" "+data.AccessToken)

	client := http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		return fiber.ErrInternalServerError
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return fiber.ErrInternalServerError
	}

	var oauthInfo NaverOAuthInfo
	err = json2.Unmarshal(body, &oauthInfo)

	if err != nil {
		return fiber.ErrInternalServerError
	}

	var user db.User
	db.DB.
		Select("password").
		Where("email = ?", oauthInfo.Response.Email).
		Find(&user)

	if user.Password == "oauth" {
		_uuid := uuid.New()

		// add session
		sess := new(db.Session)
		db.DB.Where("email = ?", oauthInfo.Response.Email).FirstOrInit(sess)

		sess.Uuid = _uuid.String()
		sess.OAuthing = 0
		sess.UserID = user.ID
		fmt.Println(sess.UserID)
		sess.Email = oauthInfo.Response.Email
		sess.Expires = time.Now().Add(time.Duration(6) * time.Hour)

		db.DB.Save(sess)

		return c.Redirect(os.Getenv("OAUTH_GLOBAL_LOGIN_REDIRECT_URI") + "?session=" + _uuid.String())
	}

	if user.Password != "oauth" && user.Password != "" {
		return c.Send([]byte("You already have account has same email"))
	}

	sess := new(db.Session)
	db.DB.Where("email = ?", oauthInfo.Response.Email).FirstOrInit(sess)
	_uuid := uuid.New()

	sess.Uuid = _uuid.String()
	sess.OAuthing = 1
	sess.Email = oauthInfo.Response.Email
	sess.Expires = time.Now().Add(time.Duration(6) * time.Hour)

	db.DB.Save(sess)

	return c.Redirect(os.Getenv("OAUTH_GLOBAL_REGISTER_REDIRECT_URI") + "?oauth=true&session=" + _uuid.String())
}
