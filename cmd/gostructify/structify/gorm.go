package structify

import (
	"fmt"
	"strings"

	"github.com/fatih/structtag"
	"github.com/snagles/gostructify/cmd/gostructify/database"
)

// TODO: add more advanced gorm tags for migrations such as `gorm:"type:varchar(100);unique"`

// GormTags adds gorm specific annotations
func GormTags(c database.Column) structtag.Tag {
	return structtag.Tag{Key: "gorm", Name: "column:" + c.Name, Options: []string{}}
}

func GormTableName(structname string, tablename string) string {
	return fmt.Sprintf(`
  // TableName manually overrides gorms defaults of taking the struct name and pluralizing it 
  // http://jinzhu.me/gorm/models.html#conventions
	func (%s *%s) TableName() string {
			return "%s"
	}`+"\n", strings.ToLower(string(structname[0])), structname, tablename)
}
