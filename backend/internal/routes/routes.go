package routes

import (
	"net/http"

	delivery "nexa/backend/internal/delivery/http"
	"nexa/backend/internal/handlers"
	"nexa/backend/internal/middleware"
	"nexa/backend/internal/repository"
	"nexa/backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(r *gin.Engine, db *gorm.DB, jwtSecret string) {
	r.Use(middleware.CORSMiddleware())

	authHandler := &handlers.AuthHandler{DB: db, JWTSecret: jwtSecret}

	pasienRepo := repository.NewPasienRepository(db)
	pasienUC := usecase.NewPasienUsecase(pasienRepo)
	pasienHandler := delivery.NewPasienHandler(pasienUC)

	rekamRepo := repository.NewRekamMedisRepository(db)
	rekamUC := usecase.NewRekamMedisUsecase(rekamRepo)
	rekamHandler := delivery.NewRekamMedisHandler(rekamUC)

	klaimRepo := repository.NewKlaimRepository(db)
	klaimUC := usecase.NewKlaimUsecase(klaimRepo)
	klaimHandler := delivery.NewKlaimHandler(klaimUC)

	adminHandler := delivery.NewAdminHandler(db)
	adminUserHandler := delivery.NewAdminUserHandler(db)
	masterHandler := delivery.NewMasterHandler(db)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.AuthMiddleware(jwtSecret), authHandler.Me)
		}

		authMw := middleware.AuthMiddleware(jwtSecret)

		pasien := api.Group("/pasien", authMw)
		{
			pasien.GET("", pasienHandler.SearchPasien)
			pasien.GET("/:id", pasienHandler.GetPasien)
			pasien.POST("", middleware.RBACMiddleware("admin_ti", "petugas_rm"), pasienHandler.CreatePasien)
			pasien.PUT("/:id", middleware.RBACMiddleware("admin_ti", "petugas_rm"), pasienHandler.UpdatePasien)
			pasien.DELETE("/:id", middleware.RBACMiddleware("admin_ti", "petugas_rm"), pasienHandler.DeletePasien)
		}

		rekam := api.Group("/rekam-medis", authMw)
		{
			rekam.GET("", rekamHandler.ListRekam)
			rekam.GET("/:id", rekamHandler.GetRekam)
			rekam.POST("", middleware.RBACMiddleware("admin_ti", "petugas_rm", "dokter_dpjp"), rekamHandler.CreateRekam)
			rekam.POST("/:id/diagnosis", middleware.RBACMiddleware("admin_ti", "dokter_dpjp"), rekamHandler.AddDiagnosis)
			rekam.POST("/:id/tindakan", middleware.RBACMiddleware("admin_ti", "perawat", "dokter_dpjp"), rekamHandler.AddTindakan)
		}

		klaim := api.Group("/klaim", authMw)
		{
			klaim.GET("", klaimHandler.SearchKlaim)
			klaim.GET("/:id", klaimHandler.GetKlaim)
			klaim.POST("", middleware.RBACMiddleware("admin_ti", "petugas_casemix"), klaimHandler.CreateKlaim)
			klaim.PATCH("/:id/status", middleware.RBACMiddleware("admin_ti", "petugas_casemix"), klaimHandler.UpdateStatus)
			klaim.DELETE("/:id", middleware.RBACMiddleware("admin_ti"), klaimHandler.DeleteKlaim)
		}

		admin := api.Group("/admin", authMw, middleware.RBACMiddleware("admin_ti"))
		{
			admin.GET("/users", adminHandler.GetUsers)
			admin.POST("/users", adminUserHandler.CreateUser)
			admin.PUT("/users/:id", adminUserHandler.UpdateUser)
			admin.DELETE("/users/:id", adminUserHandler.DeleteUser)
			admin.GET("/audit-trail", adminHandler.GetAuditTrail)
			admin.GET("/submissions", adminHandler.GetSubmissions)
		}

		master := api.Group("/master", authMw)
		{
			// ICD-10
			master.GET("/icd10", masterHandler.ListICD10)
			master.GET("/icd10/search", masterHandler.SearchICD10)
			master.GET("/icd10/:kode", masterHandler.GetICD10)
			master.POST("/icd10", middleware.RBACMiddleware("admin_ti"), masterHandler.CreateICD10)
			master.PUT("/icd10/:kode", middleware.RBACMiddleware("admin_ti"), masterHandler.UpdateICD10)
			master.DELETE("/icd10/:kode", middleware.RBACMiddleware("admin_ti"), masterHandler.DeleteICD10)
			// ICD-9 CM (Tindakan)
			master.GET("/icd9", masterHandler.ListICD9)
			master.GET("/icd9/search", masterHandler.SearchICD9)
			master.GET("/icd9/:kode", masterHandler.GetICD9)
			master.POST("/icd9", middleware.RBACMiddleware("admin_ti"), masterHandler.CreateICD9)
			master.PUT("/icd9/:kode", middleware.RBACMiddleware("admin_ti"), masterHandler.UpdateICD9)
			master.DELETE("/icd9/:kode", middleware.RBACMiddleware("admin_ti"), masterHandler.DeleteICD9)
			// CBGs
			master.GET("/cbgs", masterHandler.ListCBGS)
			master.GET("/cbgs/:kode", masterHandler.GetCBGS)
			master.POST("/cbgs", middleware.RBACMiddleware("admin_ti"), masterHandler.CreateCBGS)
			master.PUT("/cbgs/:kode", middleware.RBACMiddleware("admin_ti"), masterHandler.UpdateCBGS)
			master.DELETE("/cbgs/:kode", middleware.RBACMiddleware("admin_ti"), masterHandler.DeleteCBGS)
			// Existing (keep for backward compatibility)
			master.GET("/roles", masterHandler.GetRoles)
			master.GET("/units", masterHandler.GetUnits)
			master.POST("/units", middleware.RBACMiddleware("admin_ti"), masterHandler.CreateUnit)
			master.PUT("/units/:id", middleware.RBACMiddleware("admin_ti"), masterHandler.UpdateUnit)
			master.DELETE("/units/:id", middleware.RBACMiddleware("admin_ti"), masterHandler.DeleteUnit)
		}

		// Rekam Medis - Update and Delete
		rekam.PUT("/:id", middleware.RBACMiddleware("admin_ti", "petugas_rm", "dokter_dpjp"), rekamHandler.UpdateRekam)
		rekam.DELETE("/:id", middleware.RBACMiddleware("admin_ti", "dokter_dpjp"), rekamHandler.DeleteRekam)

		// Dashboard
		dashboardHandler := delivery.NewDashboardHandler(db)
		api.GET("/dashboard/stats", authMw, dashboardHandler.GetStats)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Nexa API running"})
	})
}
