package gin

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/knowlet/router-map/models"
	"github.com/knowlet/router-map/pkg/geoip2"
)

type GeoJSON struct {
	Type       string                 `json:"type"`
	Geometry   Geometry               `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func (s *Service) ListCarsHandler(c *gin.Context) {
	cars, err := s.DAO.Car.List()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	// add random proxy geo
	p := cars[rand.Intn(len(cars))]
	json := []GeoJSON{
		{
			Type: "Feature",
			Geometry: Geometry{
				Type:        "Point",
				Coordinates: []float64{121.531852722168, 25.0477600097656},
			},
			Properties: map[string]interface{}{
				"origin_id":           0,
				"origin_city":         "Taipei",
				"origin_country":      "Taiwan",
				"origin_lon":          121.531852722168,
				"origin_lat":          25.0477600097656,
				"destination_id":      p.ID,
				"destination_city":    p.City,
				"destination_country": p.Country,
				"destination_lon":     p.Longitude,
				"destination_lat":     p.Latitude,
			},
		},
	}
	for _, car := range cars {
		json = append(json, GeoJSON{
			Type: "Feature",
			Geometry: Geometry{
				Type:        "Point",
				Coordinates: []float64{car.Longitude, car.Latitude},
			},
			Properties: map[string]interface{}{
				"name":                fmt.Sprintf("Car #%d", car.ID),
				"car":                 car,
				"color":               "",
				"origin_id":           p.ID,
				"origin_city":         p.City,
				"origin_country":      p.Country,
				"origin_lon":          p.Longitude,
				"origin_lat":          p.Latitude,
				"destination_id":      car.ID,
				"destination_city":    car.City,
				"destination_country": car.Country,
				"destination_lon":     car.Longitude,
				"destination_lat":     car.Latitude,
			},
		})
	}
	c.JSON(http.StatusOK, json)
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

func (s *Service) createCar(json models.Car) (car models.Car, err error) {
	// lookup addr
	addr, err := getUrlIP(json.Url)
	if err != nil {
		return car, err
	}

	// lookup geo location
	region, city, lat, lng, err := geoip2.Getip(addr)
	if err != nil {
		return car, err
	}
	json.Ip = addr
	json.Country = region
	json.City = city
	json.Latitude = float64(lat)
	json.Longitude = float64(lng)
	return s.DAO.Car.Create(json)
}

func (s *Service) NewCarHandler(c *gin.Context) {
	json := models.Car{}
	if err := c.ShouldBindBodyWith(&json, binding.MsgPack); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	car, err := s.createCar(json)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, car)
}

func (s *Service) BatchCarHandler(c *gin.Context) {
	json := []models.Car{}
	if err := c.ShouldBindBodyWith(&json, binding.MsgPack); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cars := []models.Car{}
	for _, car := range json {
		car, err := s.createCar(car)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cars = append(cars, car)
	}

	c.JSON(http.StatusOK, gin.H{
		"input":  len(json),
		"output": len(cars),
	})
}
