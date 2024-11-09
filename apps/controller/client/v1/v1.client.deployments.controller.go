package v1

import (
	v1Client "deployment-service/apps/dao/client/v1"
	"deployment-service/apps/repository/adapter"
	"fmt"

	"github.com/gin-gonic/gin"
)

type DeploymentController struct {
	v1EventLoggerDao v1Client.IDeploymentsDao
}

type IDeploymentController interface {
	GetDeploymentsByNamespace(ctx *gin.Context)
	GetTenantKubernetesInfo(ctx *gin.Context)
}

func NewDeploymentController(repository *adapter.Repository) IDeploymentController {
	return &DeploymentController{
		v1EventLoggerDao: v1Client.NewDeploymentsDao(repository),
	}
}

func (ctrl DeploymentController) GetDeploymentsByNamespace(ctx *gin.Context) {
	fmt.Println("getting deployments")
	namespace := ctx.GetString("username")
	ctrl.v1EventLoggerDao.GetDeployments(ctx, namespace)
}

func (ctrl DeploymentController) GetTenantKubernetesInfo(ctx *gin.Context) {
	fmt.Println("getting deployments")
	namespace := ctx.GetString("username")
	ctrl.v1EventLoggerDao.GetTenantKubernetesInfo(ctx, namespace)
}
