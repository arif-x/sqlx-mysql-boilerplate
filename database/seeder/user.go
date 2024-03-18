package seeder

import (
	"fmt"
	"log"
	"strconv"

	"github.com/google/uuid"

	"time"

	hash "github.com/arif-x/sqlx-mysql-boilerplate/pkg/hash"
)

func (s Seed) UserSeeder() {
	password, _ := hash.Hash([]byte("password"))

	superadmin_q := "SELECT uuid FROM roles WHERE name = 'Superadmin' LIMIT 1"
	var superadmin_role_uuid uuid.UUID
	_ = s.db.QueryRow(superadmin_q).Scan(&superadmin_role_uuid)

	inactive_q := "SELECT uuid FROM roles WHERE name = 'Inactive' LIMIT 1"
	var inactive_role_uuid uuid.UUID
	_ = s.db.QueryRow(inactive_q).Scan(&inactive_role_uuid)

	_, err := s.db.Exec(`INSERT INTO users(uuid, name, username, email, role_uuid, password, email_verified_at, is_active, created_at) VALUES (?,?,?,?,?,?,?,?,?)`,
		uuid.New(),
		"superadmin",
		"superadmin",
		"superadmin@gmail.com",
		superadmin_role_uuid,
		password,
		time.Now(),
		true,
		time.Now(),
	)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		_, err := s.db.Exec(`INSERT INTO users(uuid, name, username, email, role_uuid, password, created_at) VALUES (?,?,?,?,?,?,?)`,
			uuid.New(),
			"Name "+strconv.Itoa(i+1),
			"username"+strconv.Itoa(i+1),
			"email"+strconv.Itoa(i+1)+"@gmail.com",
			inactive_role_uuid,
			password,
			time.Now(),
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("User has successfully seeded")
}
