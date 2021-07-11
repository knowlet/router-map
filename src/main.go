package main

import (
	"github.com/knowlet/router-map/dao"
	"github.com/knowlet/router-map/pkg/db"
	"github.com/knowlet/router-map/pkg/gin"
)

func main() {
	// Initalise Service with DAO instance (which in turn wraps the connection pool).
	db, err := db.NewGormClient()
	if err != nil {
		panic(err)
	}
	s := &gin.Service{
		DAO: dao.NewDAO(db),
	}
	r := s.SetupRouter()

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()
	// r.Run(":3000") for a hard coded port
}
