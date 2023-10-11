package data

import "gorm.io/gorm"

type Placeholder struct {
	gorm.Model
	VariableOne string `json:"variable_one"`
	VariableTwo string `json:"variable_two"`
}
