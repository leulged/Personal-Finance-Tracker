package router

import (
	"github.com/gin-gonic/gin"
	"personal-finance-tracker/Delivery/handler"
	usecase "personal-finance-tracker/UseCase"
	"personal-finance-tracker/Infrastructure/service"
)

func SetupRouter(
	userUsecase *usecase.UserUsecase,
	jwtService *services.JWTService,
) *gin.Engine {

	router := gin.Default()

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUsecase)

	// Public routes
	router.POST("/register", userHandler.Register)
	// router.POST("/login", userHandler.Login) // Uncomment when login is implemented

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Personal Finance Tracker API is running",
		})
	})

	return router
}