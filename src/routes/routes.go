package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// This is route group. You can adjust it in `cmd/serve.go:NewServer`
func App(db *gorm.DB, r *gin.RouterGroup) {
	// example:
	// repos := registry.NewRepositories(db)
	// services := registry.NewServices(repos)
	// controllers := registry.NewControllers(services)

	// other route here...

	r.GET("/example", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
