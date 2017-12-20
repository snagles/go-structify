package database

type Database interface {
	Build(string, string) (*Table, error)
}
