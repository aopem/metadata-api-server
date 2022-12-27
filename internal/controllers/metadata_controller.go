package controllers

import (
	"io/ioutil"
	"metadata-api-server/internal/services"
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
	bodyData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}

	responseMetadata := mc.MetadataService.CreateMetadata(bodyData)
	c.YAML(http.StatusCreated, responseMetadata)
}

func (mc *MetadataController) GetMetadataByHash(c *gin.Context) {
	responseMetadata := mc.MetadataService.GetMetadataByHash(c.Param("hash"))
	c.YAML(http.StatusOK, responseMetadata)
}

func (mc *MetadataController) GetMetadata(c *gin.Context) {
	responseMetadata := mc.MetadataService.GetMetadata()
	c.YAML(http.StatusOK, responseMetadata)
}
