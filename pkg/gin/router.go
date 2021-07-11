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

	r.GET("/cars", s.ListCarsHandler)
	return r
}
