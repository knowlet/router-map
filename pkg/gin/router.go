package gin

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	r := gin.Default()

	r.GET("/", GetIndexHandler)
	return r
}
