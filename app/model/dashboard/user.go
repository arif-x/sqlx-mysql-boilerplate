package dashboard

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID      uuid.UUID  `db:"uuid" json:"uuid"`
	Name      string     `db:"name" json:"name"`
	Username  string     `db:"username" json:"username"`
	Email     string     `db:"email" json:"email"`
	RoleUUID  *uuid.UUID `db:"role_uuid" json:"role_uuid"`
	RoleName  *string    `db:"role_name" json:"role_name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type UserShow struct {
	UUID      uuid.UUID  `db:"uuid" json:"uuid"`
	Name      string     `db:"name" json:"name"`
	Username  string     `db:"username" json:"username"`
	Email     string     `db:"email" json:"email"`
	RoleUUID  *uuid.UUID `db:"role_uuid" json:"role_uuid"`
	RoleName  *string    `db:"role_name" json:"role_name"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type StoreUser struct {
	Name     string `json:"name" form:"name"`
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	RoleUUID string `json:"role_uuid" form:"role_uuid"`
}

type UpdateUser struct {
	Name     string `json:"name" form:"name"`
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	RoleUUID string `json:"role_uuid" form:"role_uuid"`
}
