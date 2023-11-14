package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Product struct {
	ID        uint   `gorm:"primaryKey"`
	UUID      string `gorm:"not null;type:varchar(50)"`
	Name      string `gorm:"not null;unique" json:"name" form:"name" valid:"required~Your product name is required"`
	ImageURL  string `gorm:"not null;type:varchar"`
	AdminID   uint
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (p Product) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}

func (p Product) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}
