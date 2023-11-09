package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        int    `gorm:"primaryKey"`
	UUID      string `gorm:"not null;type:varchar(50)"`
	Name      string `gorm:"not null;unique;type:varchar(50)"`
	Email     string `gorm:"not null;unique;type:varchar(50)"`
	Password  string `gorm:"not null;type:varchar(50)"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (a *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("Before creating Admin")

	if len(a.Name) < 4 {
		err = errors.New("name is too short")
	}

	return
}
