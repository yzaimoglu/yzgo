package data

import "errors"

var (
	ErrDatabaseSelect    = errors.New("could not get data from database")
	ErrDatabaseCreate    = errors.New("could not create data in database")
	ErrInterfaceToObject = errors.New("could not convert interface to object")
	ErrValidateBody      = errors.New("could not validate body")
	ErrNoResultsFound    = errors.New("no results found")
)
