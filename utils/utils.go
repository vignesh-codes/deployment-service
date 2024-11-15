package utils

import (
	"crypto/md5"
	"deployment-service/logger"
	"deployment-service/utils/response"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func GetRepoNameFromURL(url string) (string, error) {
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid URL format")
	}
	return parts[len(parts)-2] + "/" + parts[len(parts)-1], nil
}

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

// converts a struct to a map[string]interface{}
func StructToMap(mystruct any) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(mystruct)
	json.Unmarshal(inrec, &inInterface)

	// iterate through inrecs
	for field, val := range inInterface {
		fmt.Println("KV Pair: ", field, val)
	}
	return inInterface
}
