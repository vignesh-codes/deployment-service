package middlewares

import (
	"deployment-service/apps/repository/adapter"

	"github.com/gin-gonic/gin"
)

type Application struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Service      string `json:"service"`
}

// for internal apis
func ValidateHeaderSecrets(repository *adapter.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// id := c.Request.Header.Get("x-client-id")
		// secret := c.Request.Header.Get("x-client-secret")
		// if id == "" {
		// 	status := response.UnAuthorized(string(response.ErrUnauthorized))
		// 	c.JSON(status.Status(), status)
		// 	c.Abort()
		// 	return
		// }
		// if secret == "" {
		// 	status := response.UnAuthorized(string(response.ErrUnauthorized))
		// 	c.JSON(status.Status(), status)
		// 	c.Abort()
		// 	return
		// }
		username := c.Request.Header.Get("username")
		c.Set("username", username)
		c.Next()
	}
}
