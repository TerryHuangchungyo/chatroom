package model

import (
	"database/sql"
)

type HubModel struct {
	tableName string
	db        *sql.DB
}
