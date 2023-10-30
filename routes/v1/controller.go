package v1

import (
	"github.com/gofiber/fiber/v2"
	"time"
)

func GetIndex(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":     fiber.StatusOK,
		"serverTime": time.Now(),
	})
}

func Controller(app *fiber.App) {
	app.Get("/api/status", GetIndex)

	app.Post("/api/v1/auth/login", Login)
	app.Post("/api/v1/auth/register", Register)

	app.Use("/api/v1/post/get_job_posts", authMiddleWare)
	app.Post("/api/v1/post/get_job_posts", GetJobPosts)
	app.Use("/api/v1/post/write_job_post", adminMiddleWare)
	app.Put("/api/v1/post/write_job_post", WriteJobPost)
	app.Use("/api/v1/post/update_job_post", adminMiddleWare)
	app.Post("/api/v1/post/update_job_post", UpdateJobPost)
	app.Use("/api/v1/post/remove_job_post", adminMiddleWare)
	app.Delete("/api/v1/post/remove_job_post", RemoveJobPost)

	app.Use("/api/v1/post/apply_job_post", authMiddleWare)
	app.Post("/api/v1/post/apply_job_post", ApplyJobPost)

	app.Post("/api/v1/post/get_tours", GetTours)
	app.Use("/api/v1/post/write_tour", adminMiddleWare)
	app.Put("/api/v1/post/write_tour", WriteTour)
	app.Use("/api/v1/post/update_tour", adminMiddleWare)
	app.Post("/api/v1/post/update_tour", UpdateTour)
	app.Use("/api/v1/post/remove_tour", adminMiddleWare)
	app.Delete("/api/v1/post/remove_tour", RemoveTour)

	app.Post("/api/v1/post/get_party_and_events", GetPartyAndEvents)
	app.Use("/api/v1/post/write_party_and_events", adminMiddleWare)
	app.Put("/api/v1/post/write_party_and_events", WritePartyAndEvents)
	app.Use("/api/v1/post/update_party_and_events", adminMiddleWare)
	app.Post("/api/v1/post/update_party_and_events", UpdatePartyAndEvents)
	app.Use("/api/v1/post/remove_party_and_events", adminMiddleWare)
	app.Delete("/api/v1/post/remove_party_and_events", RemovePartyAndEvents)

	app.Use("/api/v1/auth/logout", authMiddleWare)
	app.Post("/api/v1/auth/logout", Logout)

	app.Use("/api/v1/user/teachers", authMiddleWare)
	app.Post("/api/v1/user/teachers", Teachers)

	app.Use("/api/v1/user/me", authMiddleWare)
	app.Get("/api/v1/user/me", Me)

	app.Use("/api/v1/auth/get_approval_queue", adminMiddleWare)
	app.Get("/api/v1/auth/get_approval_queue", GetApprovalQueue)

	app.Use("/api/v1/auth/accept_user", adminMiddleWare)
	app.Put("/api/v1/auth_accept_user", AcceptUser)

	app.Use("/api/v1/search/search_posts_and_teachers", authMiddleWare)
	app.Get("/api/v1/search/search_posts_and_teachers", SearchPostAndTeachers)
}
