package registry

import "gorm.io/gorm"

type Repositories struct {
	// regist your repositories...
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		// regist your repositories...
	}
}
