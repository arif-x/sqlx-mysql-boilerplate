package dashboard

import (
	"database/sql"

	model "github.com/arif-x/sqlx-mysql-boilerplate/app/model/dashboard"
	repo "github.com/arif-x/sqlx-mysql-boilerplate/app/repository/dashboard"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/database"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/paginate"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// PermissionIndex func gets all permission.
// @Description Get all permission.
// @Summary Get all permission
// @Tags Permission
// @Accept json
// @Produce json
// @Param page query integer false "Page"
// @Param limit query integer false "Limit"
// @Param search query string false "Search"
// @Param sort_by query string false "Sort By" Enums(id, name)
// @Param sort query string false "Sort" Enums(ASC, DESC)
// @Success 200 {object} response.PermissionsResponse
// @Failure 400,401,403 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/permission [get]
func PermissionIndex(c *fiber.Ctx) error {
	page, limit, search, sort_by, sort := paginate.Paginate(c)
	repository := repo.NewPermissionRepo(database.GetDB())

	permission, count, err := repository.Index(limit, uint(limit*(page-1)), search, sort_by, sort)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Index(c, page, limit, count, permission)
}

// PermissionShow func gets single permission.
// @Description Get single permission.
// @Summary Get single permission
// @Tags Permission
// @Accept json
// @Produce json
// @Param id path string true "Permission ID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Success 200 {object} response.PermissionResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/permission/{id} [get]
func PermissionShow(c *fiber.Ctx) error {
	ID := c.Params("id")

	repository := repo.NewPermissionRepo(database.GetDB())
	permission, err := repository.Show(ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			return response.InternalServerError(c, err)
		}
	}

	return response.Show(c, permission)
}

// PermissionStore func create permission.
// @Description Create permission.
// @Summary Create permission
// @Tags Permission
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Name" default(Permission Name)
// @Success 200 {object} response.PermissionResponse
// @Failure 400,401,403 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/permission [post]
func PermissionStore(c *fiber.Ctx) error {
	permission := &model.StorePermission{}

	if err := c.BodyParser(permission); err != nil {
		return response.BadRequest(c, err)
	}

	repository := repo.NewPermissionRepo(database.GetDB())
	res, err := repository.Store(permission)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Store(c, res)
}

// PermissionUpdate func update permission.
// @Description Update permission.
// @Summary Update permission
// @Tags Permission
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Permission ID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Param name formData string true "Name" default(Permission Name Update)
// @Success 200 {object} response.PermissionResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/permission/{id} [put]
func PermissionUpdate(c *fiber.Ctx) error {
	ID := c.Params("id")

	permission := &model.UpdatePermission{}

	if err := c.BodyParser(permission); err != nil {
		return response.BadRequest(c, err)
	}

	repository := repo.NewPermissionRepo(database.GetDB())
	res, err := repository.Update(ID, permission)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Update(c, res)
}

// PermissionDestroy func delete permission.
// @Description Delete permission.
// @Summary Delete permission
// @Tags Permission
// @Accept json
// @Produce json
// @Param id path string true "Permission ID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Success 200 {object} response.PermissionResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/permission/{id} [delete]
func PermissionDestroy(c *fiber.Ctx) error {
	ID := c.Params("id")

	repository := repo.NewPermissionRepo(database.GetDB())
	res, err := repository.Destroy(ID)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Destroy(c, res)
}
