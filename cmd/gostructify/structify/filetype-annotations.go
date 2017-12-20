package structify

import (
	"strings"

	"github.com/fatih/structtag"
	"github.com/snagles/gostructify/cmd/gostructify/database"
)

// JSONTags adds json specific annotations
func JSONTags(c database.Column) structtag.Tag {
	return structtag.Tag{Key: "json", Name: strings.ToLower(c.Name), Options: []string{"omitempty"}}
}

// CSVTags adds csv specific annotations
func CSVTags(c database.Column) structtag.Tag {
	return structtag.Tag{Key: "csv", Name: strings.ToLower(c.Name), Options: []string{"omitempty"}}
}

// XMLTags adds xml specific annotations
func XMLTags(c database.Column) structtag.Tag {
	return structtag.Tag{Key: "xml", Name: strings.ToLower(c.Name), Options: []string{"omitempty"}}
}
