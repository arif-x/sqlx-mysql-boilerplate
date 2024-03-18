package main

import (
	"github.com/arif-x/sqlx-mysql-boilerplate/cmd"
	_ "github.com/arif-x/sqlx-mysql-boilerplate/docs"
)

// Swagger Config
// @title SQLX MySQL Boilerplate API
// @version 1.0
// @description SQLX MySQL Boilerplate API Swag.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
