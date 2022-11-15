package models

type PatchByIdModel struct {
	PatchValue interface{} `json:"patchValue" validate:"required"`
	ColumnName string      `json:"columnName" validate:"required"`
	Id         string      `json:"id" validate:"required"`
}
