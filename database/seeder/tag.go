package seeder

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (s Seed) Tag() {
	for i := 0; i < 3; i++ {
		_, err := s.db.Exec(`INSERT INTO tags(uuid, name, slug, created_at) VALUES (?,?,?,?)`,
			uuid.New(),
			"Tag "+strconv.Itoa(i+1),
			"tag-"+strconv.Itoa(i+1),
			time.Now(),
		)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Tag has successfully seeded")
}
