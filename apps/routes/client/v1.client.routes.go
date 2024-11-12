package client

import (
	v1 "deployment-service/apps/controller/client/v1"
	"deployment-service/apps/repository/adapter"
	"deployment-service/logger"
	"deployment-service/middlewares"
	"fmt"

	"github.com/gin-gonic/gin"
)

func V1(group *gin.RouterGroup, repository *adapter.Repository) {
	logger.ConsoleLogger.Debug("Initialising frontend v1 group routes.")
	fmt.Println("Initialising frontend v1 group routes.")
	v1ClientDeploymentsCtrl := v1.NewDeploymentController(repository)
	group.Use(middlewares.ValidateHeaderSecrets(repository))
	{
		group.GET("/deployments/", v1ClientDeploymentsCtrl.GetDeploymentsByNamespace)
		group.GET("/tenant/", v1ClientDeploymentsCtrl.GetTenantKubernetesInfo)
		group.GET("/deployments/:deployment_name", v1ClientDeploymentsCtrl.GetDeploymentByName)
		group.POST("/createns/", v1ClientDeploymentsCtrl.CreateNamespace)
	}
}
