package dashboard

import (
	"database/sql"
	"log"

	model "github.com/arif-x/sqlx-mysql-boilerplate/app/model/dashboard"
	repo "github.com/arif-x/sqlx-mysql-boilerplate/app/repository/dashboard"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/database"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// SyncPermissionShow func gets all permissions that role has.
// @Description Get all permissions that role has.
// @Summary Get all permissions that role has
// @Tags Sync Permission
// @Accept json
// @Produce json
// @Param id path string true "Role ID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Success 200 {object} response.SyncPermissionResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/sync-permission/{id} [get]
func SyncPermissionShow(c *fiber.Ctx) error {
	ID := c.Params("id")
	repository := repo.NewSyncPermissionRepo(database.GetDB())

	data, err := repository.Show(ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			return response.InternalServerError(c, err)
		}
	}

	return response.Show(c, data)
}

// SyncPermissionUpdate func update permissions per role.
// @Description update permissions per role.
// @Summary update permissions per role
// @Tags Sync Permission
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Role ID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Param permission_uuid formData []string true "permission_uuid" collectionFormat(multi)
// @Success 200 {object} response.SyncPermissionResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/sync-permission/{id} [put]
func SyncPermissionUpdate(c *fiber.Ctx) error {
	ID := c.Params("id")
	req := &model.UpdateSyncPermission{}

	if err := c.BodyParser(req); err != nil {
		return response.BadRequest(c, err)
	}

	repository := repo.NewSyncPermissionRepo(database.GetDB())
	res, err := repository.Update(ID, req)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			log.Println(err)
			return response.InternalServerError(c, err)
		}
	}

	return response.Update(c, res)
}
