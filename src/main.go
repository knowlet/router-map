package main

import (
	"github.com/knowlet/router-map/pkg/gin"
)

func main() {
	r := gin.SetupRouter()

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()
	// r.Run(":3000") for a hard coded port
}
