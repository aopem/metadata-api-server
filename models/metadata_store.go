package models

type MetadataStore struct {
	Hash     string    `yaml:"hash"`
	Metadata *Metadata `yaml:"metadata"`
}
