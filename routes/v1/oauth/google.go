package oauth

import (
	"crypto/rand"
	"encoding/base64"
	json2 "encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"lemonaid-backend/db"
	"net/http"
	"os"
	"time"
)

var (
	// Auth Config
	googleOAuthConfig *oauth2.Config
)

func GoogleLogin(c *fiber.Ctx) error {
	googleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URI"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

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

	// Google 로그인 페이지로 리디렉션
	return c.Redirect(googleOAuthConfig.AuthCodeURL(state))
}

func GoogleCallback(c *fiber.Ctx) error {
	sess := store.Get(c)

	state := c.Query("state")
	if sess.Get("state") != state {
		return fiber.ErrUnauthorized
	}

	code := c.Query("code")
	token, err := googleOAuthConfig.Exchange(c.UserContext(), code)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return GoAuthProcessing(c, token)
}

type Error struct{}

func (e *Error) Error() string {
	return "Error occurs while logining oauth.."
}

type GoogleOAuthInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func GoAuthProcessing(c *fiber.Ctx, token *oauth2.Token) error {
	request, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return c.Send([]byte("Maybe error occurs while logining.."))
	}

	request.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := new(http.Client)
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{
				"status": fiber.StatusBadRequest,
			})
	}

	var oauthInfo GoogleOAuthInfo
	err = json2.Unmarshal(body, &oauthInfo)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	var user db.User
	db.DB.
		Select("password").
		Where("email = ?", oauthInfo.Email).
		Find(&user)

	if user.Password == "oauth" {
		_uuid := uuid.New()

		// add session
		sess := new(db.Session)
		db.DB.Where("email = ?", oauthInfo.Email).FirstOrInit(sess)

		sess.Uuid = _uuid.String()
		sess.UserID = user.ID
		sess.Email = oauthInfo.Email
		sess.Expires = time.Now().Add(time.Duration(6) * time.Hour)

		db.DB.Save(sess)

		cookie := new(fiber.Cookie)
		cookie.Name = "lsession"
		cookie.Value = _uuid.String()
		cookie.Expires = time.Now().Add(6 * time.Hour)

		c.Cookie(cookie)

		return c.Redirect("/")
	}

	if user.Password != "oauth" && user.Password != "" {
		return c.Send([]byte("You already have account has same email"))
	}

	sess := new(db.Session)
	db.DB.Where("email = ?", oauthInfo.Email).FirstOrInit(sess)
	_uuid := uuid.New()

	sess.Uuid = _uuid.String()
	sess.OAuthing = 1
	sess.Email = oauthInfo.Email
	sess.Expires = time.Now().Add(time.Duration(6) * time.Hour)

	db.DB.Save(sess)

	cookie := new(fiber.Cookie)
	cookie.Name = "lsession"
	cookie.Value = _uuid.String()
	cookie.Expires = time.Now().Add(6 * time.Hour)

	c.Cookie(cookie)

	return c.Redirect(os.Getenv("OAUTH_GLOBAL_REDIRECT_URI") + "?oauth=true&session=" + _uuid.String())
}
