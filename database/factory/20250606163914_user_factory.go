package factory
import (
	"gorm.io/gorm"
)

type UserFactory struct {
	db *gorm.DB
}

func NewUserFactory(db *gorm.DB) *UserFactory {
	return &UserFactory{db: db}
}

func (f *UserFactory) Create() error {
	// Implement your factory logic here
	// Example:
	// data := &models.User{
	//     Field1: faker.Word(),
	//     Field2: faker.Email(),
	// }
	// return f.db.Create(data).Error
	return nil
}

func (f *UserFactory) CreateBatch(count int) error {
	for range count {
		if err := f.Create(); err != nil {
			return err
		}
	}
	return nil
}
