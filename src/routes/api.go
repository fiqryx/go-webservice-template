package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// This route group starts with `/api`. You can adjust it in `cmd/serve.go`
// in the `NewServer()` function. See: cmd/serve.go:NewServer
func Api(db *gorm.DB, r *gin.RouterGroup) {
	// your route here...

	r.GET("/example", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
