package v1

import (
	v1Client "deployment-service/apps/dao/client/v1"
	"deployment-service/apps/repository/adapter"
	"fmt"

	"github.com/gin-gonic/gin"
)

type DeploymentController struct {
	v1DeploymentsDao v1Client.IDeploymentsDao
}

type IDeploymentController interface {
	GetDeploymentsByNamespace(ctx *gin.Context)
	GetTenantKubernetesInfo(ctx *gin.Context)
	GetDeploymentByName(ctx *gin.Context)
	CreateNamespace(ctx *gin.Context)
}

func NewDeploymentController(repository *adapter.Repository) IDeploymentController {
	return &DeploymentController{
		v1DeploymentsDao: v1Client.NewDeploymentsDao(repository),
	}
}

func (ctrl DeploymentController) GetDeploymentsByNamespace(ctx *gin.Context) {
	fmt.Println("getting deployments")
	namespace := ctx.GetString("username")
	ctrl.v1DeploymentsDao.GetDeployments(ctx, namespace)
}

func (ctrl DeploymentController) GetTenantKubernetesInfo(ctx *gin.Context) {
	fmt.Println("getting deployments")
	namespace := ctx.GetString("username")
	ctrl.v1DeploymentsDao.GetTenantKubernetesInfo(ctx, namespace)
}

func (ctrl DeploymentController) GetDeploymentByName(ctx *gin.Context) {
	fmt.Println("getting deployment by name")
	namespace := ctx.GetString("username")
	ctrl.v1DeploymentsDao.GetDeploymentByName(ctx, namespace)
}

func (ctrl DeploymentController) CreateNamespace(ctx *gin.Context) {
	fmt.Println("getting deployment by name")
	namespace := ctx.GetString("username")
	ctrl.v1DeploymentsDao.CreateNamespace(ctx, namespace)
}
