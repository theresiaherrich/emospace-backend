package routes

import (
	"emospaces-backend/config"
	"emospaces-backend/internal/handler"
	"emospaces-backend/internal/repository"
	"emospaces-backend/internal/service"
	"emospaces-backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.SetTrustedProxies(nil)
	r.ForwardedByClientIP = true

	r.Use(func(c *gin.Context) {
		if c.Request.Header.Get("X-Forwarded-Proto") == "https" {
			c.Request.URL.Scheme = "https"
		}
		c.Next()
	})

	db := config.DB

	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(userService)
	userHandler := handler.NewUserHandler(userService)

	moodRepo := repository.NewMoodRepository(db)
	moodService := service.NewMoodService(moodRepo)
	moodHandler := handler.NewMoodHandler(moodService, userService)

	chatRepo := repository.NewChatRepository(db)
	chatService := service.NewChatService(chatRepo)
	aiHandler := handler.NewAIHandler(chatRepo, userRepo, chatService)

	planRepo := repository.NewPlanRepository(db)
	planService := service.NewPlanService(planRepo)
	planHandler := handler.NewPlanHandler(planService)

	txRepo := repository.NewTransactionRepository(db)
	consultanRepo := repository.NewConsultanRepository(db)
	consultanService := service.NewConsultanService(consultanRepo)
	consultanHandler := handler.NewConsultanHandler(consultanService)
	
	paymentService := service.NewPaymentService(userRepo, planRepo, txRepo, consultanRepo)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	journalRepo := repository.NewJournalRepository(config.DB)
	journalService := service.NewJournalService(journalRepo)
	journalHandler := handler.NewJournalHandler(journalService)


	api := r.Group("/api")
	RegisterAuthRoutes(api, authHandler)
	RegisterMoodRoutes(api, moodHandler)
	RegisterAIRoutes(api, aiHandler)
	RegisterPaymentRoutes(api, paymentHandler)
	RegisterPlanRoutes(api, planHandler)
	RegisterUserRoutes(api, userHandler)
	RegisterAdminRoutes(api, userHandler)
	RegisterJournalRoutes(api, journalHandler)
	RegisterConsultanRoutes(api, consultanHandler)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "EmoSpace backend is running 🚀"})
	})

	return r
}
