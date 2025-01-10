package database

type Seeder interface {
	Seed() error
}
