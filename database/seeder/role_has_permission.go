package seeder

import (
	"context"
	"fmt"
	"log"

	model "github.com/arif-x/sqlx-mysql-boilerplate/app/model/dashboard"
	"github.com/google/uuid"
)

func (s Seed) RoleHasPermission() {
	superadmin_q := "SELECT uuid FROM roles WHERE name = 'Superadmin' LIMIT 1"
	var superadmin_role_uuid uuid.UUID
	_ = s.db.QueryRow(superadmin_q).Scan(&superadmin_role_uuid)

	permission, err := s.db.QueryContext(context.Background(), `SELECT uuid, name, created_at, updated_at FROM permissions`)
	if err != nil {
		log.Fatal(err)
	}

	defer permission.Close()
	var permissions []model.Permission
	for permission.Next() {
		var i model.Permission
		if err := permission.Scan(
			&i.UUID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			log.Fatal(err)
		}
		permissions = append(permissions, i)
	}
	if err := permission.Close(); err != nil {
		log.Fatal(err)
	}
	if err := permission.Err(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(permissions); i++ {
		_, err := s.db.Exec(`INSERT INTO role_has_permissions(role_uuid, permission_uuid) VALUES (?,?)`,
			superadmin_role_uuid,
			permissions[i].UUID,
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Role has permission has successfully seeded")
}
