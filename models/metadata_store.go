package models

type MetadataStore struct {
	Id       string    `yaml:"id" binding:"required,alphanum,len=8"`
	Metadata *Metadata `yaml:"metadata" binding:"required,dive"`
}
