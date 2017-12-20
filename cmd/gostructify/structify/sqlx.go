package structify

import (
	"github.com/fatih/structtag"
	"github.com/snagles/gostructify/cmd/gostructify/database"
)

// SQLXTags adds sqlx specific annotations
func SQLXTags(c database.Column) structtag.Tag {
	return structtag.Tag{Key: "db", Name: c.Name, Options: []string{}}
}
