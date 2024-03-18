package dashboard

import (
	"database/sql"

	hash "github.com/arif-x/sqlx-mysql-boilerplate/pkg/hash"

	model "github.com/arif-x/sqlx-mysql-boilerplate/app/model/dashboard"
	repo "github.com/arif-x/sqlx-mysql-boilerplate/app/repository/dashboard"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/database"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/paginate"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// UserIndex func gets all user.
// @Description Get all user.
// @Summary Get all user
// @Tags User
// @Accept json
// @Produce json
// @Param page query integer false "Page"
// @Param limit query integer false "Limit"
// @Param search query string false "Search"
// @Param sort_by query string false "Sort By" Enums(users.uuid, name, email, username)
// @Param sort query string false "Sort" Enums(ASC, DESC)
// @Success 200 {object} response.UsersResponse
// @Failure 400,401,403 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/user [get]
func UserIndex(c *fiber.Ctx) error {
	page, limit, search, sort_by, sort := paginate.Paginate(c)
	repository := repo.NewUserRepo(database.GetDB())

	users, count, err := repository.Index(limit, uint(limit*(page-1)), search, sort_by, sort)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Index(c, page, limit, count, users)
}

// UserShow func gets single user.
// @Description Get single user.
// @Summary Get single user
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID" default(f72cb686-2fc3-4147-8183-f93684780765)
// @Success 200 {object} response.UserResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/user/{id} [get]
func UserShow(c *fiber.Ctx) error {
	ID := c.Params("id")

	repository := repo.NewUserRepo(database.GetDB())
	user, err := repository.Show(ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			return response.InternalServerError(c, err)
		}
	}

	return response.Show(c, user)
}

// UserStore func create user.
// @Description Create user.
// @Summary Create user
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Name" default(Name)
// @Param username formData string true "Username" default(username)
// @Param email formData string true "Email" default(email@gmail.com)
// @Param password formData string true "Password" format(password)
// @Param role_uuid formData string true "Role ID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Success 200 {object} response.UserResponse
// @Failure 400,401,403 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/user [post]
func UserStore(c *fiber.Ctx) error {
	user := &model.StoreUser{}

	if err := c.BodyParser(user); err != nil {
		return response.BadRequest(c, err)
	}

	password, err := hash.Hash([]byte(user.Password))
	if err != nil {
		return response.InternalServerError(c, err)
	}

	user.Password = password

	repository := repo.NewUserRepo(database.GetDB())
	res, err := repository.Store(user)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Store(c, res)
}

// UserUpdate func update user.
// @Description Update user.
// @Summary Update user
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "User ID" default(f72cb686-2fc3-4147-8183-f93684780765)
// @Param name formData string true "Name" default(Name Update)
// @Param username formData string true "Username" default(usernameupdate)
// @Param email formData string true "Email" default(emailupdate@gmail.com)
// @Param password formData string false "Password" format(password)
// @Param role_uuid formData string true "Role ID" default(22863142-1cfe-48cc-9640-ea88926429a4)
// @Success 200 {object} response.UserResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/user/{id} [put]
func UserUpdate(c *fiber.Ctx) error {
	ID := c.Params("id")

	user := &model.UpdateUser{}

	if err := c.BodyParser(user); err != nil {
		return response.BadRequest(c, err)
	}

	if user.Password != "" {
		password, err := hash.Hash([]byte(user.Password))
		if err != nil {
			return response.InternalServerError(c, err)
		}

		user.Password = password
	}

	repository := repo.NewUserRepo(database.GetDB())
	res, err := repository.Update(ID, user)

	if err != nil {
		if err == sql.ErrNoRows {
			return response.NotFound(c, err)
		} else {
			return response.InternalServerError(c, err)
		}
	}

	return response.Update(c, res)
}

// UserDestroy func delete user.
// @Description Delete user.
// @Summary Delete user
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID" default(f72cb686-2fc3-4147-8183-f93684780765)
// @Success 200 {object} response.UserResponse
// @Failure 400,401,403,404 {object} response.ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /api/v1/dashboard/user/{id} [delete]
func UserDestroy(c *fiber.Ctx) error {
	ID := c.Params("id")

	repository := repo.NewUserRepo(database.GetDB())
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
