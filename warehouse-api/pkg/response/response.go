package response

import "github.com/gin-gonic/gin"

// APIResponse adalah format standar semua response di project ini
type APIResponse struct{
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`    // omitempty → tidak muncul jika nil
	Error   interface{} `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`    // untuk pagination
}

type Meta struct {
    Page      int   `json:"page"`
    Limit     int   `json:"limit"`
    Total     int64 `json:"total"`
    TotalPage int   `json:"total_page"`
}

// Success response
func OK(c *gin.Context, message string, data interface{}) {
    c.JSON(200, APIResponse{
        Success: true,
        Message: message,
        Data:    data,
    })
}

// Created response (POST berhasil)
func Created(c *gin.Context, message string, data interface{}) {
    c.JSON(201, APIResponse{
        Success: true,
        Message: message,
        Data:    data,
    })
}

// Error responses
func BadRequest(c *gin.Context, message string, err interface{}) {
    c.JSON(400, APIResponse{
        Success: false,
        Message: message,
        Error:   err,
    })
}

func Unauthorized(c *gin.Context, message string) {
    c.JSON(401, APIResponse{
        Success: false,
        Message: message,
    })
}

func Forbidden(c *gin.Context, message string) {
    c.JSON(403, APIResponse{
        Success: false,
        Message: message,
    })
}

func NotFound(c *gin.Context, message string) {
    c.JSON(404, APIResponse{
        Success: false,
        Message: message,
    })
}

func InternalError(c *gin.Context, message string) {
    c.JSON(500, APIResponse{
        Success: false,
        Message: message,
    })
}

// WithPagination untuk response list + meta pagination
func WithPagination(c *gin.Context, message string, data interface{}, meta Meta) {
    c.JSON(200, APIResponse{
        Success: true,
        Message: message,
        Data:    data,
        Meta:    &meta,
    })
}