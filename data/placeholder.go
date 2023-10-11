package data

type Placeholder struct {
	Id          string `json:"id" gorm:"primary_key"`
	VariableOne string `json:"variable_one"`
	VariableTwo string `json:"variable_two"`
}
