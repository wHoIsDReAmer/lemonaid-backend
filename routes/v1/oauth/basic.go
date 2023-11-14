package oauth

import (
	"github.com/gofiber/session/v2"
	"lemonaid-backend/db"
	"time"
)

var (
	// 세션 스토어
	store = session.New()
)

func CreateOAuthSession(uuid string, email string, oauthing int8, userId uint) {
	sess := new(db.Session)
	db.DB.Where("email = ?", email).FirstOrInit(sess)

	sess.Uuid = uuid
	sess.OAuthing = oauthing
	sess.UserID = userId
	sess.Email = email
	sess.Expires = time.Now().Add(time.Duration(6) * time.Hour)
	db.DB.Save(sess)
}
