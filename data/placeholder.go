package data

type Placeholder struct {
	Id          string `json:"id" gorm:"primary_key"`
	VariableOne string `json:"var_one"`
	VariableTwo string `json:"var_two"`
}
