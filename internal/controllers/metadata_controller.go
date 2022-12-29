package controllers

import (
	"metadata-api-server/internal/services"
	"metadata-api-server/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MetadataController struct {
	MetadataService *services.MetadataService
}

func CreateMetadataController(ms *services.MetadataService) *MetadataController {
	return &MetadataController{
		MetadataService: ms,
	}
}

func (mc *MetadataController) PutMetadata(c *gin.Context) {
	metadata := &models.Metadata{}

	// validate metadata payload
	if err := c.BindYAML(metadata); err != nil {
		c.YAML(http.StatusBadRequest, nil)
		return
	}

	responseMetadata := mc.MetadataService.CreateMetadata(metadata)
	c.YAML(http.StatusCreated, responseMetadata)
}

func (mc *MetadataController) DeleteMetadataById(c *gin.Context) {
	responseMetadata := mc.MetadataService.DeleteMetadataById(c.Param("id"))
	c.YAML(http.StatusGone, responseMetadata)
}

func (mc *MetadataController) GetMetadataById(c *gin.Context) {
	responseMetadata := mc.MetadataService.GetMetadataById(c.Param("id"))
	c.YAML(http.StatusOK, responseMetadata)
}

func (mc *MetadataController) GetMetadata(c *gin.Context) {
	responseMetadata := mc.MetadataService.GetMetadata()
	c.YAML(http.StatusOK, responseMetadata)
}
