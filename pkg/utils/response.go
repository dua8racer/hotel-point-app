package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIErrorResponse representasi respons error untuk Swagger
type APIErrorResponse struct {
	Status string `json:"status" example:"error"`
	Error  string `json:"error" example:"Invalid input parameter"`
}

// APISuccessResponse representasi respons sukses untuk Swagger
type APISuccessResponse struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message,omitempty" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// Response adalah struktur standar untuk respons API
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SendSuccessResponse mengembalikan respons sukses
func SendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// SendErrorResponse mengembalikan respons error
func SendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Status: "error",
		Error:  message,
	})
}

// SendValidationErrorResponse mengembalikan respons error validasi
func SendValidationErrorResponse(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  "error",
		"message": "Validation failed",
		"errors":  errors,
	})
}

// SendNotFoundResponse mengembalikan respons resource tidak ditemukan
func SendNotFoundResponse(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, Response{
		Status: "error",
		Error:  resource + " not found",
	})
}

// SendUnauthorizedResponse mengembalikan respons unauthorized
func SendUnauthorizedResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Status: "error",
		Error:  "Unauthorized access",
	})
}

// SendServerErrorResponse mengembalikan respons server error
func SendServerErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, Response{
		Status: "error",
		Error:  "Internal server error",
	})
}
