package main

import (
	"github.com/jmoiron/sqlx"
)

type Seeder struct {
	db *sqlx.DB
}

func (s *Seeder) seed() {

}
