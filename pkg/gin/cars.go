package gin

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/knowlet/router-map/models"
	"github.com/knowlet/router-map/pkg/geoip2"
)

func (s *Service) ListCarsHandler(c *gin.Context) {
	cars, err := s.DAO.Car.List()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, cars)
}

func getUrlIP(rawurl string) (string, error) {
	// parse a hostname and path without a scheme is invalid
	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return "", err
	}
	// check if hostname is empty
	if u.Hostname() == "" {
		return "", errors.New("empty hostname")
	}
	// lookup hostname
	addrs, err := net.LookupIP(u.Hostname())
	if err != nil {
		return "", err
	}
	// get first ip
	for _, addr := range addrs {
		ipv4 := addr.To4()
		if ipv4 == nil {
			continue
		}
		return ipv4.String(), nil
	}
	// no ip found
	return "", errors.New("could not infer host IP")
}

func connectionCheck(url string) (code string, body []byte, err error) {
	client := http.Client{
		Timeout: 120 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return code, nil, err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	return resp.Status, body, err
}

func (s *Service) CheckCarHandler(c *gin.Context) {
	json := struct {
		Url string `json:"url" binding:"required"`
	}{}
	// read json
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, body, err := connectionCheck(json.Url)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": code, "html": body})
}

func (s *Service) NewCarHandler(c *gin.Context) {
	json := models.Car{}
	// read json
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// lookup addr
	addr, err := getUrlIP(json.Url)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	// lookup geo location
	region, city, lat, lng, err := geoip2.Getip(addr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	json.Ip = addr
	json.Country = region
	json.City = city
	json.Latitude = float64(lat)
	json.Longitude = float64(lng)
	car, err := s.DAO.Car.Create(json)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, car)
}
