package oauth

import (
	"github.com/gofiber/session/v2"
)

var (
	// 세션 스토어
	store = session.New()
)
