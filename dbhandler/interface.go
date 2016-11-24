package dbhandler

import "github.com/jmoiron/sqlx"

// ConfType contains the creation of the handler
type ConfType interface {
	NewDBHandler() (*sqlx.DB, error)
}
