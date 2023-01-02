package controllers

import (
	"log"
	"metadata-api-server/internal/core"
	"metadata-api-server/internal/utils"
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
	metadataStore := &models.MetadataStore{}
	response := models.Response{
		StatusCode: http.StatusCreated,
		Data:       nil,
		Errors:     nil,
	}

	// validate payload
	if err := c.BindYAML(metadataStore); err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Errors = append(response.Errors, err.Error())
		log.Printf("[ERROR] %s", err.Error())
		c.YAML(response.StatusCode, response)
		return
	}

	// attempt to get metadata so that any changes can be detected later
	// if this is a request to update already existing data
	expected, _ := mc.MetadataService.GetMetadataById(metadataStore.Id)

	createdMetadata, err := mc.MetadataService.CreateMetadata(metadataStore)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Errors = append(response.Errors, err.Error())
		log.Printf("[ERROR] %s", err.Error())
		c.YAML(response.StatusCode, response)
		return
	}

	// determine correct status code to send back,
	// based on whether metadata record already existed
	// and whether or not it has been modified
	if expected != nil {
		if utils.MetadataEqual(createdMetadata, expected) {
			response.StatusCode = http.StatusNotModified
		} else {
			response.StatusCode = http.StatusOK
		}
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
