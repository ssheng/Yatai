package controllersv1

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/bentoml/yatai/api-server/services"
	"github.com/bentoml/yatai/api-server/transformers/transformersv1"
	"github.com/bentoml/yatai/common/utils"
	"github.com/bentoml/yatai/schemas/schemasv1"
)

type deploymentRevisionController struct {
	baseController
}

var DeploymentRevisionController = deploymentRevisionController{}

type ListDeploymentRevisionSchema struct {
	schemasv1.ListQuerySchema
	GetDeploymentSchema
}

func (c *deploymentRevisionController) List(ctx *gin.Context, schema *ListDeploymentRevisionSchema) (*schemasv1.DeploymentRevisionListSchema, error) {
	deployment, err := schema.GetDeployment(ctx)
	if err != nil {
		return nil, err
	}

	if err = DeploymentController.canView(ctx, deployment); err != nil {
		return nil, err
	}

	deploymentRevisions, total, err := services.DeploymentRevisionService.List(ctx, services.ListDeploymentRevisionOption{
		BaseListOption: services.BaseListOption{
			Start:  utils.UintPtr(schema.Start),
			Count:  utils.UintPtr(schema.Count),
			Search: schema.Search,
		},
		DeploymentId: deployment.ID,
	})
	if err != nil {
		return nil, errors.Wrap(err, "list deploymentRevisions")
	}

	deploymentRevisionSchemas, err := transformersv1.ToDeploymentRevisionSchemas(ctx, deploymentRevisions)
	return &schemasv1.DeploymentRevisionListSchema{
		BaseListSchema: schemasv1.BaseListSchema{
			Total: total,
			Start: schema.Start,
			Count: schema.Count,
		},
		Items: deploymentRevisionSchemas,
	}, err
}