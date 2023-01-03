package models

type Response struct {
	StatusCode int           `yaml:"statusCode"`
	Data       []interface{} `yaml:"data"`
	Errors     []interface{} `yaml:"errors"`
}
