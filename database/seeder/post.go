package seeder

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	model "github.com/arif-x/sqlx-mysql-boilerplate/app/model/dashboard"
	"github.com/google/uuid"

	"time"
)

func (s Seed) PostSeeder() {
	user, err := s.db.QueryContext(context.Background(), `SELECT uuid, name, email, username, created_at, updated_at FROM users`)
	if err != nil {
		log.Fatal(err)
	}

	defer user.Close()
	var users []model.User
	for user.Next() {
		var i model.User
		if err := user.Scan(
			&i.UUID,
			&i.Name,
			&i.Email,
			&i.Username,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			log.Fatal(err)
		}
		users = append(users, i)
	}
	if err := user.Close(); err != nil {
		log.Fatal(err)
	}
	if err := user.Err(); err != nil {
		log.Fatal(err)
	}

	tag, err := s.db.QueryContext(context.Background(), `SELECT uuid, name, created_at, updated_at FROM tags`)
	if err != nil {
		log.Fatal(err)
	}

	defer tag.Close()
	var categories []model.Tag
	for tag.Next() {
		var i model.Tag
		if err := tag.Scan(
			&i.UUID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			log.Fatal(err)
		}
		categories = append(categories, i)
	}
	if err := tag.Close(); err != nil {
		log.Fatal(err)
	}
	if err := tag.Err(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(users); i++ {
		for j := 0; j < 3; j++ {
			ShuffleTag(categories)
			_, err := s.db.Exec(`INSERT INTO posts(uuid, tag_uuid, user_uuid, title, thumbnail, content, slug, keyword, created_at) VALUES (?,?,?,?,?,?,?,?,?)`,
				uuid.New(),
				categories[0].UUID,
				users[i].UUID,
				"Title "+strconv.Itoa(i+1)+"-"+strconv.Itoa(j+1),
				"https://4.bp.blogspot.com/-JU8lLIDYcq4/UkWR38K8pAI/AAAAAAAAQxw/Z-3UaPjKgBw/s1600/images.jpg",
				"Post Title "+strconv.Itoa(i+1)+"-"+strconv.Itoa(j+1)+" By "+users[i].Username,
				"title-"+strconv.Itoa(i+1)+"-"+strconv.Itoa(j+1),
				"Title 1, Title",
				time.Now(),
			)
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("Post has successfully seeded")
}

func ShuffleTag(r []model.Tag) {
	for i := len(r) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		r[i], r[j] = r[j], r[i]
	}
}
