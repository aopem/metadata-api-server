package controllers

import (
	"metadata-api-server/internal/services"
	"metadata-api-server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueryController struct {
	QueryService    *services.QueryService
	MetadataService *services.MetadataService
}

func CreateQueryController(qs *services.QueryService, ms *services.MetadataService) *QueryController {
	return &QueryController{
		QueryService:    qs,
		MetadataService: ms,
	}
}

func (qc *QueryController) PutMetadataQuery(c *gin.Context) {
	query := &models.Query{}

	// validate query
	if err := c.ShouldBind(query); err != nil {
		c.YAML(http.StatusBadRequest, nil)
		return
	}

	// get all matching metadata IDs
	queryResults := qc.QueryService.ExecuteQuery(query)
	c.YAML(http.StatusCreated, queryResults)
}
