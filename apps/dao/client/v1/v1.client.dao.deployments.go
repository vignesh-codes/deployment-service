package v1

import (
	"deployment-service/apps/repository/adapter"
	"deployment-service/apps/svc"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeploymentDao struct {
	ServiceRepo *svc.ServiceRepository
}

type IDeploymentsDao interface {
	GetDeployments(ctx *gin.Context, namespace string)
	GetTenantKubernetesInfo(ctx *gin.Context, namespace string)
	GetDeploymentByName(ctx *gin.Context, namespace string)
	CreateNamespace(ctx *gin.Context, namespace string)
}

func NewDeploymentsDao(repository *adapter.Repository) IDeploymentsDao {
	return &DeploymentDao{
		ServiceRepo: svc.NewServiceRepo(repository),
	}
}

func (dao DeploymentDao) GetDeployments(ctx *gin.Context, namespace string) {
	var response []map[string]interface{}
	response, _ = dao.ServiceRepo.DeploymentService.GetDeploymentsByNamespace(namespace)

	ctx.JSON(http.StatusOK, map[string]interface{}{"message": response})
	ctx.Abort()
}

func (dao DeploymentDao) GetTenantKubernetesInfo(ctx *gin.Context, namespace string) {
	response, _ := dao.ServiceRepo.DeploymentService.GetTenantKubernetesInfo(namespace)

	ctx.JSON(http.StatusOK, map[string]interface{}{"message": response})
	ctx.Abort()
}

func (dao DeploymentDao) GetDeploymentByName(ctx *gin.Context, namespace string) {
	response, _ := dao.ServiceRepo.DeploymentService.GetDeploymentByName(namespace, ctx.Param("deployment_name"))

	ctx.JSON(http.StatusOK, map[string]interface{}{"message": response})
	ctx.Abort()
}

func (dao DeploymentDao) CreateNamespace(ctx *gin.Context, namespace string) {
	response := dao.ServiceRepo.DeploymentService.CreateNamespace(namespace)

	ctx.JSON(http.StatusOK, map[string]interface{}{"message": "successfully created namespace", "err": response})
	ctx.Abort()
}
