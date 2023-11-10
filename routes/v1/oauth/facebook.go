package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
)

func FacebookLogin(c *fiber.Ctx) error {
	clientId := os.Getenv("FACEBOOK_OAUTH_CID")
	redirectURI := os.Getenv("FACEBOOK_OAUTH_REDIRECT_URI")

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
	return c.Redirect("https://www.facebook.com/v18.0/dialog/oauth?client_id=" +
		clientId + "&redirect_uri=" + redirectURI + "&state=" + state +
		"&scope=email" +
		"&response_type=code" +
		"&auth_type=rerequest")
}

func FacebookCallback(c *fiber.Ctx) error {
	sess := store.Get(c)

	state := c.Query("state")

	if sess.Get("state") != state {
		return fiber.ErrUnauthorized
	}

	code := c.Query("code")

	//clientId := os.Getenv("FACEBOOK_OAUTH_CID")
	//clientSecret := os.Getenv("FACEBOOK_OAUTH_SECRET")

	fmt.Println(code)

	return fiber.ErrBadGateway
	//reqeust, err := http.NewRequest("GET",
	//	"https://nid.naver.com/oauth2.0/token?grant_type=authorization_code&client_id="+
	//		clientId+"&client_secret="+clientSecret+"&code="+code, nil)
	//
	//if err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	//		"status":  fiber.StatusBadRequest,
	//		"message": "Error occurs while logining oauth..",
	//	})
	//}
	//
	//client := http.Client{}
	//resp, err := client.Do(reqeust)
	//
	//if err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	//		"status":  fiber.StatusBadRequest,
	//		"message": "Error occurs while logining oauth..",
	//	})
	//}
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//
	//if err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	//		"status":  fiber.StatusBadRequest,
	//		"message": "Error occurs while logining oauth..",
	//	})
	//}
	//
	//var data NaverToken
	//err = json2.Unmarshal(body, &data)
	//if err != nil {
	//	return fiber.ErrInternalServerError
	//}

	//return NaverAuthProcessing(c, data)
}

//func NaverAuthProcessing(c *fiber.Ctx, data NaverToken) error {
//	request, err := http.NewRequest("GET", "https://openapi.naver.com/v1/nid/me", nil)
//	if err != nil {
//		return c.Send([]byte("Maybe error occurs while logining.."))
//	}
//
//	request.Header.Set("Authorization", data.TokenType+" "+data.AccessToken)
//
//	client := http.Client{}
//	resp, err := client.Do(request)
//
//	if err != nil {
//		return fiber.ErrInternalServerError
//	}
//
//	body, err := ioutil.ReadAll(resp.Body)
//
//	if err != nil {
//		return fiber.ErrInternalServerError
//	}
//
//	var oauthInfo NaverOAuthInfo
//	err = json2.Unmarshal(body, &oauthInfo)
//
//	if err != nil {
//		return fiber.ErrInternalServerError
//	}
//
//	var user db.User
//	db.DB.
//		Select("id, password").
//		Where("email = ?", oauthInfo.Response.Email).
//		Find(&user)
//
//	if user.Password == "oauth" {
//		_uuid := uuid.New()
//
//		// add session
//		sess := new(db.Session)
//		db.DB.Where("email = ?", oauthInfo.Response.Email).FirstOrInit(sess)
//
//		sess.Uuid = _uuid.String()
//		sess.OAuthing = 0
//		sess.UserID = user.ID
//		sess.Email = oauthInfo.Response.Email
//		sess.Expires = time.Now().Add(time.Duration(6) * time.Hour)
//
//		db.DB.Save(sess)
//
//		return c.Redirect(os.Getenv("OAUTH_GLOBAL_LOGIN_REDIRECT_URI") + "?session=" + _uuid.String())
//	}
//
//	if user.Password != "oauth" && user.Password != "" {
//		return c.Send([]byte("You already have account has same email"))
//	}
//
//	sess := new(db.Session)
//	db.DB.Where("email = ?", oauthInfo.Response.Email).FirstOrInit(sess)
//	_uuid := uuid.New()
//
//	sess.Uuid = _uuid.String()
//	sess.OAuthing = 1
//	sess.Email = oauthInfo.Response.Email
//	sess.Expires = time.Now().Add(time.Duration(6) * time.Hour)
//
//	db.DB.Save(sess)
//
//	return c.Redirect(os.Getenv("OAUTH_GLOBAL_REGISTER_REDIRECT_URI") + "?oauth=true&session=" + _uuid.String())
//}
