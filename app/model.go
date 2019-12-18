package app

import "context"

// Repository interface describe functions available on top of storage layer.
type Repository interface {

	// OK - returns bool flag whether connection to databse is up and running.
	OK(context.Context) (bool, error)
}
