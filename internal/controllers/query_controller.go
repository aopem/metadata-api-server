package controllers

import (
	"log"
	"metadata-api-server/internal/core"
	"metadata-api-server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type QueryController struct {
	QueryService    core.QueryService
	MetadataService core.MetadataService
}

func CreateQueryController(qs core.QueryService, ms core.MetadataService) *QueryController {
	return &QueryController{
		QueryService:    qs,
		MetadataService: ms,
	}
}

func (qc *QueryController) PutMetadataQuery(c *gin.Context) {
	query := &models.Query{}
	response := models.Response{
		StatusCode: http.StatusCreated,
		Data:       nil,
		Errors:     nil,
	}

	// validate query
	if err := c.ShouldBind(query); err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Errors = append(response.Errors, err.Error())
		log.Printf("[ERROR] %s", err.Error())
		c.YAML(response.StatusCode, response)
		return
	}

	// get all matching metadata IDs
	queryResults, err := qc.QueryService.ExecuteQuery(query)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Errors = append(response.Errors, err.Error())
		log.Printf("[ERROR] %s", err.Error())
		c.YAML(response.StatusCode, response)
		return
	}

	for i := range queryResults {
		response.Data = append(response.Data, queryResults[i])
	}
	c.YAML(response.StatusCode, response)
}
