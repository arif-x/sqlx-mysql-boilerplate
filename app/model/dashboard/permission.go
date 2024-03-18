package dashboard

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	UUID      uuid.UUID  `db:"uuid" json:"uuid"`
	Name      string     `db:"name" json:"name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type ShowPermission struct {
	UUID      uuid.UUID  `db:"uuid" json:"uuid"`
	Name      string     `db:"name" json:"name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type StorePermission struct {
	Name string `json:"name" form:"name"`
}

type UpdatePermission struct {
	Name string `json:"name" form:"name"`
}
