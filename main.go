package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Define models

type Province struct {
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	NameEn               string `json:"name_en"`
	FullName             string `json:"full_name"`
	FullNameEn           string `json:"full_name_en"`
	CodeName             string `json:"code_name"`
	AdministrativeUnitID int    `json:"administrative_unit_id"`
}

type District struct {
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	NameEn               string `json:"name_en"`
	FullName             string `json:"full_name"`
	FullNameEn           string `json:"full_name_en"`
	CodeName             string `json:"code_name"`
	ProvinceCode         string `json:"province_code"`
	AdministrativeUnitID int    `json:"administrative_unit_id"`
}

type Ward struct {
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	NameEn               string `json:"name_en"`
	FullName             string `json:"full_name"`
	FullNameEn           string `json:"full_name_en"`
	CodeName             string `json:"code_name"`
	DistrictCode         string `json:"district_code"`
	AdministrativeUnitID int    `json:"administrative_unit_id"`
}

var db *gorm.DB
var err error

func init() {
	// Database connection string for PostgreSQL (update with your credentials)
	connStr := "user=postgres dbname=vietnamese_administrative_units password=1 host=localhost port=5432"
	db, err = gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	db.AutoMigrate(&Province{}, &District{}, &Ward{})
}

func getProvinces(c *gin.Context) {
	var provinces []Province
	if err := db.Find(&provinces).Error; err != nil {
		errorResponse := ClientResponse(http.StatusBadRequest, "Could not fetch provinces", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}
	successResponse := ClientResponse(http.StatusOK, "Provinces fetched successfully", provinces, nil)
	c.JSON(http.StatusOK, successResponse)
}

func getDistricts(c *gin.Context) {

	var districts []District
	province_code := c.Query("province_code")

	if province_code != "" {
		if err := db.Where("province_code = ?", province_code).Find(&districts).Error; err != nil {
			errorResponse := ClientResponse(http.StatusBadRequest, "Could not fetch districts", nil, err.Error())
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
	} else {
		if err := db.Find(&districts).Error; err != nil {
			errorResponse := ClientResponse(http.StatusBadRequest, "Could not fetch districts", nil, err.Error())
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
	}

	successResponse := ClientResponse(http.StatusOK, "Provinces fetched successfully", districts, nil)
	c.JSON(http.StatusOK, successResponse)
}

func getWards(c *gin.Context) {

	var wards []Ward
	district_code := c.Query("district_code")
	if district_code != "" {
		if err := db.Where("district_code = ?", district_code).Find(&wards).Error; err != nil {
			errorResponse := ClientResponse(http.StatusBadRequest, "Could not fetch wards", nil, err.Error())
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
	} else {
		if err := db.Find(&wards).Error; err != nil {
			errorResponse := ClientResponse(http.StatusBadRequest, "Could not fetch wards", nil, err.Error())
			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
	}

	successResponse := ClientResponse(http.StatusOK, "Provinces fetched successfully", wards, nil)
	c.JSON(http.StatusOK, successResponse)
}

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())

	api := r.Group("/api/location")
	{
		api.GET("/province", getProvinces) // List all provinces
		api.GET("/district", getDistricts) // List all districts
		api.GET("/ward", getWards)         // List all wards
	}

	// Start server
	port := ":8090"
	fmt.Println("Server is running on port", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Error      interface{} `json:"error"`
}

func ClientResponse(statusCode int, message string, data interface{}, err interface{}) Response {
	return Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Error:      err,
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		// c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
