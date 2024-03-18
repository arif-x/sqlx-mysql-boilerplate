package dashboard

import (
	jsonutil "github.com/arif-x/sqlx-mysql-boilerplate/pkg/json"
	"github.com/google/uuid"
)

type ShowSyncPermission struct {
	UUID       uuid.UUID         `db:"uuid" json:"uuid"`
	Name       string            `db:"name" json:"name"`
	Permission *jsonutil.JSONRaw `db:"permission" json:"permission"`
}

type UpdateSyncPermission struct {
	PermissionUUID []string `json:"permission_uuid" form:"permission_uuid"`
}
