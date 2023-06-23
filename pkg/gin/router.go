package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/knowlet/router-map/dao"
)

type Service struct {
	DAO *dao.DAO
}

func (s Service) SetupRouter() *gin.Engine {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	r := gin.Default()

	r.GET("/", GetIndexHandler)

	api := r.Group("/api")
	api.GET("/cars", s.ListCarsHandler)
	api.POST("/new", s.NewCarHandler)
	api.POST("/delete", s.DeleteCarHandler)
	api.POST("/batch", s.BatchCarHandler)
	api.POST("/check", s.CheckCarHandler)
	return r
}
