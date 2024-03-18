package seeder

import "github.com/jmoiron/sqlx"

func (s Seed) PopulateDB() {
	seedFunctions := []func(){
		s.Role,
		s.UserSeeder,
		s.Tag,
		s.PostSeeder,
		s.Permission,
		s.RoleHasPermission,
	}

	for _, seedFunc := range seedFunctions {
		seedFunc()
	}
}

// Seed struct.
type Seed struct {
	db *sqlx.DB
}

// NewSeed return a Seed with a pool of connection to a dabase.
func NewSeed(db *sqlx.DB) Seed {
	return Seed{
		db: db,
	}
}
