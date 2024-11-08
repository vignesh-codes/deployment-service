package client

import (
	"deployment-service/apps/repository/adapter"
	"deployment-service/logger"

	"github.com/gin-gonic/gin"
)

func V1(group *gin.RouterGroup, repository *adapter.Repository) {
	logger.ConsoleLogger.Debug("Initialising frontend v1 group routes.")

}
