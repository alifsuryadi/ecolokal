package middleware

import (
	"strings"

	"github.com/alifsuryadi/ecolokal/config"
	"github.com/alifsuryadi/ecolokal/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            utils.ErrorResponse(c, 401, "Authorization header required")
            c.Abort()
            return
        }
        
        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        
        claims, err := utils.ValidateJWT(tokenString, cfg)
        if err != nil {
            utils.ErrorResponse(c, 401, "Invalid token")
            c.Abort()
            return
        }
        
        c.Set("userID", claims.UserID)
        c.Set("userEmail", claims.Email)
        c.Set("userRole", claims.Role)
        
        c.Next()
    }
}

func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("userRole")
        if !exists || role != "admin" {
            utils.ErrorResponse(c, 403, "Admin access required")
            c.Abort()
            return
        }
        c.Next()
    }
}

func PetugasOrAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("userRole")
        if !exists || (role != "petugas" && role != "admin") {
            utils.ErrorResponse(c, 403, "Petugas or Admin access required")
            c.Abort()
            return
        }
        c.Next()
    }
}