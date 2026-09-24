package routes

import (
	"nexa/backend/internal/handlers"
	"nexa/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB, jwtSecret string) {
	r.Use(middleware.CORSMiddleware())

	authHandler := &handlers.AuthHandler{DB: db, JWTSecret: jwtSecret}

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(jwtSecret), authHandler.Me)
		}
	}
}
