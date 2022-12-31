package controllers

import (
	"log"
	"metadata-api-server/internal/core"
	"metadata-api-server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MetadataController struct {
	MetadataService core.MetadataService
}

func CreateMetadataController(ms core.MetadataService) *MetadataController {
	return &MetadataController{
		MetadataService: ms,
	}
}

func (mc *MetadataController) PutMetadata(c *gin.Context) {
	metadata := &models.Metadata{}
	response := models.Response{
		StatusCode: http.StatusCreated,
		Data:       nil,
		Errors:     nil,
	}

	// validate metadata payload
	if err := c.BindYAML(metadata); err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Errors = append(response.Errors, err.Error())
		log.Printf("[ERROR] %s", err.Error())
		c.YAML(response.StatusCode, response)
		return
	}

	createdMetadata, err := mc.MetadataService.CreateMetadata(metadata)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Errors = append(response.Errors, err.Error())
		log.Printf("[ERROR] %s", err.Error())
		c.YAML(response.StatusCode, response)
		return
	}

	response.Data = createdMetadata
	c.YAML(response.StatusCode, response)
}

func (mc *MetadataController) DeleteMetadataById(c *gin.Context) {
	response := models.Response{
		StatusCode: http.StatusGone,
		Data:       nil,
		Errors:     nil,
	}

	deletedMetadata, err := mc.MetadataService.DeleteMetadataById(c.Param("id"))
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Errors = append(response.Errors, err.Error())
		log.Printf("[ERROR] %s", err.Error())
		c.YAML(response.StatusCode, response)
		return
	}

	response.Data = deletedMetadata
	c.YAML(response.StatusCode, response)
}

func (mc *MetadataController) GetMetadataById(c *gin.Context) {
	response := models.Response{
		StatusCode: http.StatusOK,
		Data:       nil,
		Errors:     nil,
	}

	metadata, err := mc.MetadataService.GetMetadataById(c.Param("id"))
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Errors = append(response.Errors, err.Error())
		log.Printf("[ERROR] %s", err.Error())
		c.YAML(response.StatusCode, response)
		return
	}

	response.Data = metadata
	c.YAML(response.StatusCode, response)
}

func (mc *MetadataController) GetMetadata(c *gin.Context) {
	response := models.Response{
		StatusCode: http.StatusOK,
		Data:       nil,
		Errors:     nil,
	}

	metadataList, err := mc.MetadataService.GetMetadata()
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Errors = append(response.Errors, err.Error())
		log.Printf("[ERROR] %s", err.Error())
		c.YAML(response.StatusCode, response)
		return
	}

	response.Data = metadataList
	c.YAML(response.StatusCode, response)
}
