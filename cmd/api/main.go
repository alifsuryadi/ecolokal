package main

import (
	"log"

	"github.com/alifsuryadi/ecolokal/config"
	"github.com/alifsuryadi/ecolokal/internal/delivery/http/router"
	"github.com/alifsuryadi/ecolokal/pkg/database"

	_ "github.com/alifsuryadi/ecolokal/docs" // swagger docs
)

// @title EcoLokal API
// @version 1.0
// @description API untuk Bank Sampah dan Sistem Jemput Sampah Terjadwal
// @termsOfService http://swagger.io/terms/

// @contact.name EcoLokal Support
// @contact.email support@ecolokal.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
    // Load configuration
    cfg := config.LoadConfig()
    
    // Initialize database
    db, err := database.NewPostgresDB(cfg)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Setup router
    r := router.SetupRouter(db, cfg)
    
    // Start server
    log.Printf("Server starting on port %s", cfg.AppPort)
    if err := r.Run(":" + cfg.AppPort); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}