package gin

import (
	"errors"
	"fmt"
	"io"
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
	json := []GeoJSON{}
	for _, car := range cars {
		json = append(json, GeoJSON{
			Type: "Feature",
			Geometry: Geometry{
				Type:        "Point",
				Coordinates: []float64{car.Longitude, car.Latitude},
			},
			Properties: map[string]interface{}{
				"name": fmt.Sprintf("Car #%d", car.ID),
				"car":  car,
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
	if json.Country == "" {
		json.Country = region
	}
	if json.City == "" {
		json.City = city
	}
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

type deleteRequest struct {
	ID uint `json:"id" binding:"required"`
}

func (s *Service) DeleteCarHandler(c *gin.Context) {
	var json deleteRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.DAO.Car.Delete(json.ID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, json)
}

// GetProvincesHandler returns a list of provinces
func (s *Service) GetProvincesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.DAO.Car.GetProvinces())
}

// ExportCarsHandler return a csv file of cars with filter of province
func (s *Service) ExportCarsHandler(c *gin.Context) {
	province := c.Param("province")
	cars, err := s.DAO.Car.GetCars(province)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	csv := "id,url,ip,user,pass,country,province,city,unit,longitude,latitude,vendor,protocol\n"
	for _, car := range cars {
		csv += fmt.Sprintf("%d,%s,%s,%s,%s,%s,%s,%s,%s,%f,%f,%s,%s\n",
			car.ID,
			car.Url,
			car.Ip,
			car.User,
			car.Pass,
			car.Country,
			car.Province,
			car.City,
			car.Unit,
			car.Longitude,
			car.Latitude,
			car.Vendor,
			car.Protocol,
		)
	}
	c.Header("Content-Disposition", "attachment; filename=cars_"+province+".csv")
	c.Data(http.StatusOK, "text/csv", []byte(csv))
}

// UpdateCarHandler updates a car
func (s *Service) UpdateCarHandler(c *gin.Context) {
	json := models.Car{}
	if err := c.ShouldBindBodyWith(&json, binding.MsgPack); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	car, err := s.DAO.Car.Update(&json)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, car)
}
