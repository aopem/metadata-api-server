package controllers

import (
	"log"
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
	}

	// get all matching metadata IDs
	queryResults, err := qc.QueryService.ExecuteQuery(query)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Errors = append(response.Errors, err.Error())
		log.Printf("[ERROR] %s", err.Error())
	}

	response.Data = queryResults
	c.YAML(response.StatusCode, response)
}
