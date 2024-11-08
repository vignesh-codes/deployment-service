package utils

import (
	"crypto/md5"
	"deployment-service/logger"
	"deployment-service/utils/response"
	"encoding/binary"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func GenerateUUID(str string) int64 {
	// Generate MD5 hash of the input string
	hash := md5.Sum([]byte(str))

	// Extract the first 8 bytes of the MD5 hash
	b := hash[:8]

	// Convert the byte slice to an int64 value
	uuid := int64(binary.BigEndian.Uint64(b))

	// Get the current timestamp in nanoseconds
	timestamp := time.Now().UnixNano()

	// Combine the timestamp and UUID to create a unique value
	uniqueID := (timestamp << 16) | (uuid & 0x0000FFFF)

	// Take the absolute value of the unique ID to ensure it's positive
	uniqueID = int64(math.Abs(float64(uniqueID)))

	return uniqueID
}

func BindJSON(c *gin.Context, req interface{}) bool {
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.FullPath())
		reply := response.UnsupportedMediaType(msg)
		c.JSON(reply.Status(), reply)
		return false
	}

	if err := c.ShouldBind(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {

			var errList []string
			for _, err := range errs {
				e := fmt.Sprintf("field=%s tag=%s required=%s kind=%s value=%v", err.Field(), err.Tag(), err.Param(), err.Kind(), err.Value())
				errList = append(errList, e)
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"status_code": http.StatusBadRequest,
				"message":     "Validation errors",
				"errors":      errList,
			})
			return false
		}

		if err.Error() == "EOF" {
			status := response.BadRequest("Empty Body")
			c.JSON(status.Status(), status)
			return false
		}

		fallback := response.InternalServerError(logger.BindJSONtoStruct, "BindJSON::fallback", err)
		c.JSON(fallback.Status(), fallback)
		return false
	}

	return true
}