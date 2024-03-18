package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID            uuid.UUID  `db:"uuid" json:"uuid"`
	Name            string     `db:"name" json:"name"`
	Username        string     `db:"username" json:"username"`
	Email           string     `db:"email" json:"email"`
	Password        string     `db:"password" json:"password"`
	RoleUUID        string     `db:"role_uuid" json:"role_uuid"`
	IsActive        bool       `db:"is_active" json:"is_active"`
	EmailVerifiedAt *time.Time `db:"email_verified_at" json:"email_verified_at"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at"`
}

type Login struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type Register struct {
	Name     string `json:"name" form:"name"`
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	RoleUUID string `json:"role_uuid" form:"role_uuid"`
	Password string `json:"password" form:"password"`
}

type Role struct {
	UUID uuid.UUID `db:"uuid" json:"uuid"`
	Name string    `db:"name" json:"name"`
}

type Permission struct {
	UUID uuid.UUID `db:"uuid" json:"uuid"`
	Name string    `db:"name" json:"name"`
}

type RoleHasPermission struct {
	RoleUUID       uuid.UUID `db:"role_uuid" json:"role_uuid"`
	PermissionUUID uuid.UUID `db:"permission_uuid" json:"permission_uuid"`
}

type ForgotPassword struct {
	Username string `json:"username" form:"username"`
}

type ChangeForgotPassword struct {
	Token    string `json:"token" form:"token"`
	Password string `json:"password" form:"password"`
}
