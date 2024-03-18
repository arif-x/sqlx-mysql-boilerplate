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

// RoleIndex func gets all role.
// @Description Get all role.
// @Summary Get all role
// @Tags Role
// @Accept json
// @Produce json
// @Param page query integer false "Page"
// @Param limit query integer false "Limit"
// @Param search query string false "Search"
// @Param sort_by query string false "Sort By" Enums(id, name)
// @Param sort query string false "Sort" Enums(ASC, DESC)
// @Success 200 {object} response.TagsResponse
// @Failure 400,401,403 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/role [get]
func RoleIndex(c *fiber.Ctx) error {
	page, limit, search, sort_by, sort := paginate.Paginate(c)
	repository := repo.NewRoleRepo(database.GetDB())

	role, count, err := repository.Index(limit, uint(limit*(page-1)), search, sort_by, sort)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Index(c, page, limit, count, role)
}

// RoleShow func gets single role.
// @Description Get single role.
// @Summary Get single role
// @Tags Role
// @Accept json
// @Produce json
// @Param id path string true "Role ID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Success 200 {object} response.RolesResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/role/{id} [get]
func RoleShow(c *fiber.Ctx) error {
	ID := c.Params("id")

	repository := repo.NewRoleRepo(database.GetDB())
	role, err := repository.Show(ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			return response.InternalServerError(c, err)
		}
	}

	return response.Show(c, role)
}

// RoleStore func create role.
// @Description Create role.
// @Summary Create role
// @Tags Role
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Name" default(Role Name)
// @Param is_active formData bool true "Is Active"
// @Success 200 {object} response.RoleResponse
// @Failure 400,401,403 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/role [post]
func RoleStore(c *fiber.Ctx) error {
	role := &model.StoreRole{}

	if err := c.BodyParser(role); err != nil {
		return response.BadRequest(c, err)
	}

	repository := repo.NewRoleRepo(database.GetDB())
	res, err := repository.Store(role)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Store(c, res)
}

// RoleUpdate func update role.
// @Description Update role.
// @Summary Update role
// @Tags Role
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Role ID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Param name formData string true "Name" default(Role Name Update)
// @Param is_active formData bool true "Is Active"
// @Success 200 {object} response.RoleResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/role/{id} [put]
func RoleUpdate(c *fiber.Ctx) error {
	ID := c.Params("id")

	role := &model.UpdateRole{}

	if err := c.BodyParser(role); err != nil {
		return response.BadRequest(c, err)
	}

	repository := repo.NewRoleRepo(database.GetDB())
	res, err := repository.Update(ID, role)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			return response.InternalServerError(c, err)
		}
	}

	return response.Update(c, res)
}

// RoleDestroy func delete role.
// @Description Delete role.
// @Summary Delete role
// @Tags Role
// @Accept json
// @Produce json
// @Param id path string true "Role ID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Success 200 {object} response.RoleResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/role/{id} [delete]
func RoleDestroy(c *fiber.Ctx) error {
	ID := c.Params("id")

	repository := repo.NewRoleRepo(database.GetDB())
	res, err := repository.Destroy(ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			return response.InternalServerError(c, err)
		}
	}

	return response.Destroy(c, res)
}
