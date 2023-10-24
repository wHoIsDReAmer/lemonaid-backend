package v1

import (
	"github.com/gofiber/fiber/v2"
	"lemonaid-backend/db"
	"time"
)

func GetIndex(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":     fiber.StatusOK,
		"serverTime": time.Now(),
	})
}

func authMiddleWare(c *fiber.Ctx) error {
	token := c.Get(fiber.HeaderAuthorization, "")

	var session db.Session
	if rst := db.DB.Select("email").Where("uuid = ?", token).Find(&session); rst.RowsAffected == 0 {
		return c.JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	c.Locals("email", session.Email)
	return c.Next()
}

func Controller(app *fiber.App) {
	app.Get("/api/status", GetIndex)

	app.Post("/api/v1/auth/login", Login)
	app.Post("/api/v1/auth/register", Register)

	app.Use("/api/v1/post/get_job_posts", authMiddleWare)
	app.Post("/api/v1/post/ge_job_posts", GetJobPosts)
	app.Use("/api/v1/post/write_job_post", authMiddleWare)
	app.Put("/api/v1/post/write_job_post", WriteJobPost)
	app.Use("/api/v1/post/remove_job_post", authMiddleWare)
	app.Delete("/api/v1/post/remove_job_post", RemoveJobPost)

	app.Use("/api/v1/post/apply_job_post", authMiddleWare)
	app.Post("/api/v1/post/apply_job_post", ApplyJobPost)

	app.Post("/api/v1/post/get_tours", GetTours)
	app.Use("/api/v1/post/write_tour", authMiddleWare)
	app.Put("/api/v1/post/write_tour", WriteTour)
	app.Use("/api/v1/post/remove_tour", authMiddleWare)
	app.Delete("/api/v1/post/remove_tour", RemoveTour)

	app.Post("/api/v1/post/get_party_and_events", GetPartyAndEvents)
	app.Use("/api/v1/post/write_party_and_events", authMiddleWare)
	app.Put("/api/v1/post/write_party_and_events", WritePartyAndEvents)
	app.Use("/api/v1/post/remove_party_and_events", authMiddleWare)
	app.Delete("/api/v1/post/remove_party_and_events", RemovePartyAndEvents)

	app.Use("/api/v1/auth/logout", authMiddleWare)
	app.Post("/api/v1/auth/logout", Logout)

	app.Use("/api/v1/user/teachers", authMiddleWare)
	app.Post("/api/v1/user/teachers", Teachers)

	app.Use("/api/v1/user/me", authMiddleWare)
	app.Get("/api/v1/user/me", Me)
}
