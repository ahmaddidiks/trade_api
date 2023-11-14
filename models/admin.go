package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Admin struct {
	ID        uint       `gorm:"primaryKey"`
	UUID      string     `gorm:"not null"`
	Name      string     `gorm:"not null;unique" json:"name" form:"name" valid:"required~Your name is required"`
	Email     string     `gorm:"not null;unique" json:"email" form:"email" valid:"required~Your email is required, email~Invalid email format"`
	Password  string     `gorm:"not null;" json:"password" form:"password" valid:"required~Your password is required, minstringlength(6)~Password has to have a minimum length of 6 characters"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (a *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(a)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}

func (a *Admin) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(a)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}
