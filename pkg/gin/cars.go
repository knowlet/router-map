package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) ListCarsHandler(c *gin.Context) {
	cars, err := s.DAO.Car.List()
	if err != nil {
		c.AbortWithStatus(http.StatusNoContent)
	}
	c.JSON(http.StatusOK, cars)
}
