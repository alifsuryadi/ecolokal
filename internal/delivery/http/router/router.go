package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"database/sql"

	"github.com/alifsuryadi/ecolokal/config"
	"github.com/alifsuryadi/ecolokal/internal/delivery/http/handler"
	"github.com/alifsuryadi/ecolokal/internal/delivery/http/middleware"
	"github.com/alifsuryadi/ecolokal/internal/repository"
	"github.com/alifsuryadi/ecolokal/internal/usecase"
)

func SetupRouter(db *sql.DB, cfg *config.Config) *gin.Engine {
    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    wasteTypeRepo := repository.NewWasteTypeRepository(db)
    pickupRepo := repository.NewPickupRepository(db)
    transactionRepo := repository.NewTransactionRepository(db)
    
    // Initialize usecases
    authUsecase := usecase.NewAuthUsecase(userRepo, cfg)
    userUsecase := usecase.NewUserUsecase(userRepo)
    wasteTypeUsecase := usecase.NewWasteTypeUsecase(wasteTypeRepo)
    pickupUsecase := usecase.NewPickupUsecase(pickupRepo, userRepo, wasteTypeRepo, transactionRepo)
    transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, userRepo)
    
    // Initialize handlers
    authHandler := handler.NewAuthHandler(authUsecase)
    userHandler := handler.NewUserHandler(userUsecase)
    wasteTypeHandler := handler.NewWasteTypeHandler(wasteTypeUsecase)
    pickupHandler := handler.NewPickupHandler(pickupUsecase)
    transactionHandler := handler.NewTransactionHandler(transactionUsecase)
    
    // Setup router
    r := gin.Default()
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // Swagger documentation
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    
    // API routes
    api := r.Group("/api")
    {
        // Auth routes (public)
        auth := api.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
        }
        
        // Login route (as per requirement)
        api.POST("/users/login", authHandler.Login)
        
        // Protected routes
        protected := api.Group("")
        protected.Use(middleware.AuthMiddleware(cfg))
        {
            // User routes
            users := protected.Group("/users")
            {
                users.GET("/profile", userHandler.GetProfile)
                users.GET("/points", userHandler.GetUserPoints)
                users.GET("/role/:role", middleware.AdminOnly(), userHandler.GetUsersByRole)
            }
            
            // Waste type routes
            // Waste type routes
            wasteTypes := protected.Group("/waste-types")
            {
                wasteTypes.GET("", wasteTypeHandler.GetAll)
                wasteTypes.GET("/active", wasteTypeHandler.GetActiveTypes)
                wasteTypes.GET("/:id", wasteTypeHandler.GetByID)
                wasteTypes.POST("", middleware.AdminOnly(), wasteTypeHandler.Create)
                wasteTypes.PUT("/:id", middleware.AdminOnly(), wasteTypeHandler.Update)
                wasteTypes.DELETE("/:id", middleware.AdminOnly(), wasteTypeHandler.Delete)
            }
            
            // Pickup routes
            pickups := protected.Group("/pickups")
            {
                pickups.POST("", pickupHandler.CreatePickupRequest)
                pickups.GET("/my", pickupHandler.GetMyPickups)
                pickups.GET("/pending", middleware.AdminOnly(), pickupHandler.GetPendingPickups)
                pickups.GET("/petugas", middleware.PetugasOrAdmin(), pickupHandler.GetPetugasPickups)
                pickups.GET("/:id", pickupHandler.GetPickupByID)
                pickups.PUT("/:id/status", middleware.PetugasOrAdmin(), pickupHandler.UpdatePickupStatus)
                pickups.PUT("/:id/items", middleware.PetugasOrAdmin(), pickupHandler.UpdatePickupItems)
                pickups.PUT("/:id/assign", middleware.AdminOnly(), pickupHandler.AssignPickupToPetugas)
            }
            
            // Transaction routes
            transactions := protected.Group("/transactions")
            {
                transactions.POST("", middleware.AdminOnly(), transactionHandler.CreateTransaction)
                transactions.GET("/my", transactionHandler.GetMyTransactions)
            }
        }
    }
    
    return r
}