package data

type NoSQLDatabase interface {
	GetSingle(query string, vars interface{}) (interface{}, error)
	GetMultiple(query string, vars interface{}) ([]interface{}, error)
	Query(query string, vars interface{}) (interface{}, error)
	QueryWithoutResponse(query string, vars interface{}) (interface{}, error)
	Create(table string, data interface{}) (interface{}, error)
	Patch(id string, data interface{}) (interface{}, error)
	Delete(id string) error
	Init() error
	Close() error
}

type SQLDatabase interface {
	Init() error
	Close() error
}
