package seeder

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s Seed) Role() {
	var arr = []string{
		"Superadmin",
		"Verified",
		"Inactive",
	}
	for i := 0; i < len(arr); i++ {
		_, err := s.db.Exec(`INSERT INTO roles(uuid, name, created_at) VALUES (?,?,?)`,
			uuid.New(),
			arr[i],
			time.Now(),
		)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Role has successfully seeded")
}
