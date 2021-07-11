package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) ListCarsHandler(c *gin.Context) {
	cars, err := s.DAO.Car.List()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
	}
	c.JSON(http.StatusOK, cars)
}
