package dashboard

import (
	"context"
	"fmt"

	model "github.com/arif-x/sqlx-mysql-boilerplate/app/model/dashboard"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/database"
)

type SyncPermissionRepository interface {
	Show(uuid string) (model.ShowSyncPermission, error)
	Update(uuid string, request *model.UpdateSyncPermission) (model.ShowSyncPermission, error)
}

type SyncPermissionRepo struct {
	db *database.DB
}

func (repo *SyncPermissionRepo) Show(uuid string) (model.ShowSyncPermission, error) {
	_select := `
	roles.uuid,
    roles.name,
    IFNULL(
        (
            SELECT JSON_ARRAYAGG(
                JSON_OBJECT(
                    'uuid', permission_uuid_p,
                    'name', permission_name
                )
            ) 
            FROM (
                SELECT 
                    role_has_permissions.*, 
                    permissions.uuid AS permission_uuid_p, 
                    permissions.name AS permission_name 
                FROM 
                    role_has_permissions
                    LEFT JOIN permissions ON permissions.uuid = role_has_permissions.permission_uuid
                WHERE 
                    role_has_permissions.role_uuid = ?
            ) AS permissions
        ), '[]'
    ) AS permission
	`

	query := fmt.Sprintf(`SELECT %s FROM roles WHERE roles.uuid = ?`, _select)
	var items model.ShowSyncPermission
	err := repo.db.QueryRowContext(context.Background(), query, uuid, uuid).Scan(
		&items.UUID,
		&items.Name,
		&items.Permission,
	)
	if err != nil {
		return model.ShowSyncPermission{}, err
	}

	return items, nil
}

func (repo *SyncPermissionRepo) Update(uuid string, request *model.UpdateSyncPermission) (model.ShowSyncPermission, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return model.ShowSyncPermission{}, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	defer tx.Commit()

	rhp := `DELETE FROM role_has_permissions WHERE role_uuid = ?`
	_, rerr := tx.ExecContext(context.Background(), rhp, uuid)

	if rerr != nil {
		tx.Rollback()
		return model.ShowSyncPermission{}, rerr
	}

	for i := 0; i < len(request.PermissionUUID); i++ {
		check_query := fmt.Sprintf(`SELECT uuid, name FROM permissions WHERE uuid = ?`)
		var items model.Permission
		cerr := tx.QueryRowContext(context.Background(), check_query, request.PermissionUUID[i]).Scan(
			&items.UUID,
			&items.Name,
		)
		if cerr != nil {
			tx.Rollback()
			return model.ShowSyncPermission{}, cerr
		}

		query := `INSERT INTO role_has_permissions (role_uuid, permission_uuid) VALUES (?, ?)`
		_, ierr := tx.ExecContext(context.Background(), query, uuid, request.PermissionUUID[i])
		if ierr != nil {
			tx.Rollback()
			return model.ShowSyncPermission{}, ierr
		}
	}

	_select := `
	roles.uuid,
    roles.name,
    IFNULL(
        (
            SELECT JSON_ARRAYAGG(
                JSON_OBJECT(
                    'uuid', permission_uuid_p,
                    'name', permission_name
                )
            ) 
            FROM (
                SELECT 
                    role_has_permissions.*, 
                    permissions.uuid AS permission_uuid_p, 
                    permissions.name AS permission_name 
                FROM 
                    role_has_permissions
                    LEFT JOIN permissions ON permissions.uuid = role_has_permissions.permission_uuid
                WHERE 
                    role_has_permissions.role_uuid = ?
            ) AS permissions
        ), '[]'
    ) AS permission
	`

	query := fmt.Sprintf(`SELECT %s FROM roles WHERE roles.uuid = ?`, _select)
	var items model.ShowSyncPermission
	serr := tx.QueryRowContext(context.Background(), query, uuid, uuid).Scan(
		&items.UUID,
		&items.Name,
		&items.Permission,
	)
	if serr != nil {
		return model.ShowSyncPermission{}, serr
	}

	return items, nil
}

func NewSyncPermissionRepo(db *database.DB) SyncPermissionRepository {
	return &SyncPermissionRepo{db}
}
