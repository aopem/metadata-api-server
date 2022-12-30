package models

type Query struct {
	Title           string `form:"title" binding:"omitempty"`
	Version         string `form:"version" binding:"omitempty"`
	MaintainerName  string `form:"maintainerName" binding:"omitempty"`
	MaintainerEmail string `form:"maintainerEmail" binding:"omitempty,email"`
	Company         string `form:"company" binding:"omitempty"`
	Website         string `form:"website" binding:"omitempty,url"`
	Source          string `form:"source" binding:"omitempty"`
	License         string `form:"license" binding:"omitempty"`
	Description     string `form:"description" binding:"omitempty"`
}
