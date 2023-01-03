package models

type Maintainer struct {
	Name  string `yaml:"name" binding:"required"`
	Email string `yaml:"email" binding:"required,email"`
}

type Metadata struct {
	Title       string       `yaml:"title" binding:"required"`
	Version     string       `yaml:"version" binding:"required"`
	Maintainers []Maintainer `yaml:"maintainers" binding:"required,dive"`
	Company     string       `yaml:"company" binding:"required"`
	Website     string       `yaml:"website" binding:"required,url"`
	Source      string       `yaml:"source" binding:"required"`
	License     string       `yaml:"license" binding:"required"`
	Description string       `yaml:"description" binding:"required"`
}

type MetadataStore struct {
	Id       string    `yaml:"id" binding:"omitempty,required,alphanum,len=8"`
	Metadata *Metadata `yaml:"metadata" binding:"required,dive"`
}
