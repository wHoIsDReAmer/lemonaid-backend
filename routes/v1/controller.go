package v1

import (
	"github.com/gofiber/fiber/v2"
	"lemonaid-backend/routes/v1/oauth"
	"lemonaid-backend/routes/v1/pay"
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

	app.Get("/api/v1/post/job_post", GetJobPosts)
	app.Use("/api/v1/post/job_post", authMiddleWare)
	app.Post("/api/v1/post/job_post", WriteJobPost)
	app.Put("/api/v1/post/job_post", UpdateJobPost)
	app.Delete("/api/v1/post/job_post", RemoveJobPost)

	app.Use("/api/v1/post/pending_job_post", adminMiddleWare)
	app.Get("/api/v1/post/pending_job_post", GetPendingJobPosts)
	app.Put("/api/v1/post/pending_job_post", AcceptPendingJobPost)
	app.Delete("/api/v1/post/pending_job_post", DenyPendingJobPost)

	app.Use("/api/v1/post/apply_job_post", authMiddleWare)
	app.Post("/api/v1/post/apply_job_post", ApplyJobPost)

	app.Get("/api/v1/post/popular_job_post", GetPopularJobPosts)

	// idk pending
	//app.Get("/api/v1/post/apply_job_post", ApplyJobPostApprovalQueue)
	//app.Put("/api/v1/post/apply_job_post", AcceptApplyJobPost)
	//app.Delete("/api/v1/post/apply_job_post", DenyApplyJobPost)

	app.Get("/api/v1/post/tour", GetTours)
	app.Use("/api/v1/post/tour", adminMiddleWare)
	app.Post("/api/v1/post/tour", WriteTour)
	app.Put("/api/v1/post/tour", UpdateTour)
	app.Delete("/api/v1/post/tour", RemoveTour)

	app.Get("/api/v1/post/party_and_events", GetPartyAndEvents)

	app.Use("/api/v1/post/party_and_events", adminMiddleWare)
	app.Post("/api/v1/post/party_and_events", WritePartyAndEvents)
	app.Put("/api/v1/post/party_and_events", UpdatePartyAndEvents)
	app.Delete("/api/v1/post/party_and_events", RemovePartyAndEvents)

	app.Use("/api/v1/post/job_images_upload", authMiddleWare)
	app.Post("/api/v1/post/job_images_upload", UploadImageToJobPost)

	app.Use("/api/v1/post/pending_job_images_upload", authMiddleWare)
	app.Post("/api/v1/post/pending_job_images_upload", UploadImageToPendingJobPost)

	app.Use("/api/v1/post/images_upload", adminMiddleWare)
	app.Post("/api/v1/post/images_upload", UploadImageToPost)

	app.Use("/api/v1/auth/logout", authMiddleWare)
	app.Post("/api/v1/auth/logout", Logout)

	app.Use("/api/v1/user/users", adminMiddleWare)
	app.Get("/api/v1/user/users", Users)

	//app.Use("/api/v1/user/teachers", authMiddleWare)
	app.Get("/api/v1/user/teachers", Teachers)

	app.Use("/api/v1/user/resume", authMiddleWare)
	app.Get("/api/v1/user/resume", ResumeDownload)

	app.Use("/api/v1/user/me", authMiddleWare)
	app.Get("/api/v1/user/me", Me)

	app.Use("/api/v1/user", adminMiddleWare)
	app.Get("/api/v1/user", User)
	app.Put("/api/v1/user", UserEdit)

	app.Use("/api/v1/auth/approval_user", adminMiddleWare)
	app.Get("/api/v1/auth/approval_user", UserApprovalQueue)
	app.Put("/api/v1/auth/approval_user", AcceptUser)
	app.Delete("/api/v1/auth/approval_user", DenyUser)

	//app.Use("/api/v1/search/search_posts_and_teachers", authMiddleWare)
	app.Get("/api/v1/search/search_posts_and_teachers", SearchPostAndTeachers)

	app.Get("/api/v1/oauth/google/login", oauth.GoogleLogin)
	app.Get("/api/v1/oauth/google", oauth.GoogleCallback)

	app.Get("/api/v1/oauth/naver/login", oauth.NaverLogin)
	app.Get("/api/v1/oauth/naver", oauth.NaverCallback)

	app.Get("/api/v1/oauth/kakao/login", oauth.KakaoLogin)
	app.Get("/api/v1/oauth/kakao", oauth.KakaoCallback)

	app.Get("/api/v1/oauth/facebook/login", oauth.FacebookLogin)
	app.Get("/api/v1/oauth/facebook", oauth.FacebookCallback)

	app.Post("/api/v1/pay/payapp_feedback", pay.PayAppFeedback)
}
