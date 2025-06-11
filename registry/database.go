package registry

import (
	"webservices/packages/structers"

	"gorm.io/gorm"
)

var Database = &structers.DatabaseRegistry{
	Enums: []structers.Enum{
		// example:
		// {Name: "gendre", Values: []string{"male", "female"}},
	},
	Tables: []string{
		// tables for backup
	},
	Models: []any{
		// models for auto-migrate, example:
		// model.User{},
	},
	Extensions: []string{
		// example:
		// "uuid-ossp",
	},
	Factories: []func(*gorm.DB) error{
		// register factory for seeding, example:
		// func(db *gorm.DB) error {
		// 	return factory.NewUserFactory(db).CreateBatch(1)
		// },
	},
}
