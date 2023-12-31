package oauth

import (
	"crypto/rand"
	"encoding/base64"
	json2 "encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io"
	"lemonaid-backend/db"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func KakaoLogin(c *fiber.Ctx) error {
	cid := os.Getenv("KAKAO_OAUTH_CID")
	redirectUri := os.Getenv("KAKAO_OAUTH_REDIRECT_URI")

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

	url := fmt.Sprintf("https://kauth.kakao.com/oauth/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=account_email",
		cid, redirectUri, state)

	return c.Redirect(url)
}

type KakaoToken struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type KakaoOAuthInfo struct {
	KakaoAccount KakaoAccount `json:"kakao_account"`
}

type KakaoAccount struct {
	HasEmail     bool   `json:"has_email"`
	IsEmailValid bool   `json:"is_email_valid"`
	Email        string `json:"email"`
}

func KakaoCallback(c *fiber.Ctx) error {
	sess := store.Get(c)
	defer sess.Destroy()

	state := c.Query("state", "")

	if sess.Get("state") != state {
		return fiber.ErrUnauthorized
	}
	code := c.Query("code", "")

	// 받은 code로 token값 가져오기
	form := url.Values{}
	form.Add("grant_type", "authorization_code")
	form.Add("client_id", os.Getenv("KAKAO_OAUTH_CID"))
	form.Add("redirect_uri", os.Getenv("KAKAO_OAUTH_REDIRECT_URI"))
	form.Add("code", code)
	//form.Add() client_secret, but it's optional

	req, err := http.NewRequest("POST", "https://kauth.kakao.com/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return fiber.ErrInternalServerError
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	var tokenInfo KakaoToken
	err = json2.Unmarshal(body, &tokenInfo)

	if err != nil {
		return fiber.ErrInternalServerError
	}

	userReq, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	userReq.Header.Add("Authorization", tokenInfo.TokenType+" "+tokenInfo.AccessToken)

	res, err = client.Do(userReq)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	var oauthInfo KakaoOAuthInfo
	err = json2.Unmarshal(body, &oauthInfo)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	if !oauthInfo.KakaoAccount.HasEmail || !oauthInfo.KakaoAccount.IsEmailValid {
		return c.Status(fiber.StatusBadRequest).
			SendString("You don't have email to proceed or not valid email ")
	}

	var user db.User
	db.DB.
		Select("id, password").
		Where("email = ?", oauthInfo.KakaoAccount.Email).
		Find(&user)

	if user.Password == "oauth" {
		_uuid := uuid.New()
		CreateOAuthSession(_uuid.String(), oauthInfo.KakaoAccount.Email, 0, user.ID)

		return c.Redirect(os.Getenv("OAUTH_GLOBAL_LOGIN_REDIRECT_URI") + "?session=" + _uuid.String())
	}

	if user.Password != "oauth" && user.Password != "" {
		return c.Send([]byte("You already have account has same email"))
	}

	_uuid := uuid.New()
	CreateOAuthSession(_uuid.String(), oauthInfo.KakaoAccount.Email, 1, user.ID)

	return c.Redirect(os.Getenv("OAUTH_GLOBAL_REGISTER_REDIRECT_URI") + "?oauth=true&session=" + _uuid.String())
}
