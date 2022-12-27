package models

type MetadataStore struct {
	Id       string    `yaml:"id" binding:"required"`
	Metadata *Metadata `yaml:"metadata" binding:"required,dive"`
}
