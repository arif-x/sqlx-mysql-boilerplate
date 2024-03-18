package dashboard

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	UUID      uuid.UUID  `db:"uuid" json:"uuid"`
	Name      string     `db:"name" json:"name"`
	Slug      string     `db:"slug" json:"slug"`
	IsActive  bool       `db:"is_active" json:"is_active"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type TagShow struct {
	UUID      uuid.UUID  `db:"uuid" json:"uuid"`
	Name      string     `db:"name" json:"name"`
	Slug      string     `db:"slug" json:"slug"`
	IsActive  bool       `db:"is_active" json:"is_active"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type StoreTag struct {
	Name     string `json:"name" form:"name"`
	Slug     string `json:"slug" form:"slug"`
	IsActive bool   `json:"is_active" form:"is_active"`
}

type UpdateTag struct {
	Name     string `json:"name" form:"name"`
	Slug     string `json:"slug" form:"slug"`
	IsActive bool   `json:"is_active" form:"is_active"`
}
