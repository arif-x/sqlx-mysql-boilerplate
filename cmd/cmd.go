package cmd

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/arif-x/sqlx-mysql-boilerplate/config"
	seeds "github.com/arif-x/sqlx-mysql-boilerplate/database/seeder"
	"github.com/arif-x/sqlx-mysql-boilerplate/pkg/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func PrintTextFunc() {
	fmt.Println("SQLX GoFiber Boilerplate. \nVisit https://github.com/arif-x/sqlx-mysql-boilerplate")
}

func ServeFunc() {
	config.LoadAllConfigs(".env")
	server.Serve()
}

func SeedFunc() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("can't load .env file. error: %v", err)
	}
	config.LoadDBCfg()
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.DBCfg().User, config.DBCfg().Password, config.DBCfg().Host, config.DBCfg().Port, config.DBCfg().Name)
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		log.Fatalf("error opening a connection with the database %s\n", err)
	}

	seed := seeds.NewSeed(db)
	seed.PopulateDB()
	fmt.Println("Database seeder has successfully executed")
}

func MigrateMake(fileName string) {
	ext := "sql"
	dir := "database/migration"
	seq := fileName
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting the current working directory:", err)
		return
	}
	run := exec.Command("migrate", "create", "-ext", ext, "-dir", dir, "-seq", seq)
	run.Dir = workDir

	output, err := run.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(string(output))
}

func MigrateUpFunc() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("can't load .env file. error: %v", err)
	}
	config.LoadDBCfg()
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting the current working directory:", err)
		return
	}
	command := fmt.Sprintf(`migrate -path ./database/migration/ -database "mysql://%s:%s@tcp(%s:%d)/%s?parseTime=true" -verbose up`,
		config.DBCfg().User, config.DBCfg().Password, config.DBCfg().Host, config.DBCfg().Port, config.DBCfg().Name)
	run := exec.Command("sh", "-c", command)
	run.Dir = workDir

	output, err := run.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Output:", string(output))
		return
	}

	fmt.Println(string(output))
}

func MigrateDownFunc(step string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("can't load .env file. error: %v", err)
	}
	config.LoadDBCfg()
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting the current working directory:", err)
		return
	}
	number, err := strconv.Atoi(step)
	if err != nil {
		log.Fatal("Invalid number")
	}
	command := fmt.Sprintf(`migrate -path ./database/migration/ -database "mysql://%s:%s@tcp(%s:%d)/%s?parseTime=true" -verbose down %d`,
		config.DBCfg().User, config.DBCfg().Password, config.DBCfg().Host, config.DBCfg().Port, config.DBCfg().Name, number)
	run := exec.Command("sh", "-c", command)
	run.Dir = workDir
	output, err := run.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Output:", string(output))
		return
	}
	fmt.Println(string(output))
}

func MigrateFreshFunc() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("can't load .env file. error: %v", err)
	}
	config.LoadDBCfg()
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting the current working directory:", err)
		return
	}
	command := fmt.Sprintf(`migrate -path ./database/migration/ -database "mysql://%s:%s@tcp(%s:%d)/%s?parseTime=true" -verbose down`,
		config.DBCfg().User, config.DBCfg().Password, config.DBCfg().Host, config.DBCfg().Port, config.DBCfg().Name)
	run := exec.Command("sh", "-c", command)
	run.Dir = workDir
	run.Stdin = strings.NewReader("y\n")

	output, err := run.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Output:", string(output))
		return
	}

	fmt.Println(string(output))
}

func MakeController(fileName string) {
	dir := filepath.Join("app", "http", "controller", filepath.Dir(fileName))

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	filePath := filepath.Join(dir, filepath.Base(fileName)+".go")

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		fmt.Printf("File %s already exists.\n", filePath)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	contentTemplate := `package {{.PackageName}}

	`

	dirName := filepath.Base(filepath.Dir(filepath.Clean(fileName)))

	if dirName == "." {
		dirName = "controller"
	}

	data := struct {
		PackageName    string
		FileName       string
		ControllerName string
	}{
		PackageName:    dirName,
		FileName:       fileName,
		ControllerName: strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName)),
	}

	tmpl, err := template.New("controller").Parse(contentTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	if err := tmpl.Execute(file, data); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("File created successfully: %s\n", filePath)
}

func MakeMiddleware(fileName string) {
	dir := filepath.Join("app", "http", "middleware", filepath.Dir(fileName))

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	filePath := filepath.Join(dir, filepath.Base(fileName)+".go")

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		fmt.Printf("File %s already exists.\n", filePath)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	contentTemplate := `package {{.PackageName}}

	`

	dirName := filepath.Base(filepath.Dir(filepath.Clean(fileName)))

	if dirName == "." {
		dirName = "middleware"
	}

	data := struct {
		PackageName    string
		FileName       string
		ControllerName string
	}{
		PackageName:    dirName,
		FileName:       fileName,
		ControllerName: strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName)),
	}

	tmpl, err := template.New("middleware").Parse(contentTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	if err := tmpl.Execute(file, data); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("File created successfully: %s\n", filePath)
}

func MakeModel(fileName string) {
	dir := filepath.Join("app", "model", filepath.Dir(fileName))

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	filePath := filepath.Join(dir, filepath.Base(fileName)+".go")

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		fmt.Printf("File %s already exists.\n", filePath)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	contentTemplate := `package {{.PackageName}}

	`

	dirName := filepath.Base(filepath.Dir(filepath.Clean(fileName)))

	if dirName == "." {
		dirName = "model"
	}

	data := struct {
		PackageName    string
		FileName       string
		ControllerName string
	}{
		PackageName:    dirName,
		FileName:       fileName,
		ControllerName: strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName)),
	}

	tmpl, err := template.New("model").Parse(contentTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	if err := tmpl.Execute(file, data); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("File created successfully: %s\n", filePath)
}

func MakeRepository(fileName string) {
	dir := filepath.Join("app", "repository", filepath.Dir(fileName))

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	filePath := filepath.Join(dir, filepath.Base(fileName)+".go")

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		fmt.Printf("File %s already exists.\n", filePath)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	contentTemplate := `package {{.PackageName}}

	`

	dirName := filepath.Base(filepath.Dir(filepath.Clean(fileName)))

	if dirName == "." {
		dirName = "repository"
	}

	data := struct {
		PackageName    string
		FileName       string
		ControllerName string
	}{
		PackageName:    dirName,
		FileName:       fileName,
		ControllerName: strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName)),
	}

	tmpl, err := template.New("repository").Parse(contentTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	if err := tmpl.Execute(file, data); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("File created successfully: %s\n", filePath)
}

func GenerateSwag() {
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting the current working directory:", err)
		return
	}
	run := exec.Command("swag", "init", "--parseDependency", "--parseInternal")
	run.Dir = workDir

	run.Stdout = os.Stdout
	run.Stderr = os.Stderr

	err = run.Run()
	if err != nil {
		log.Fatalf("Error running swag: %v", err)
	}
}
