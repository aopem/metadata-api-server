package models

type MetadataStore struct {
	Id       string    `yaml:"id"`
	Metadata *Metadata `yaml:"metadata"`
}
