package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        int    `gorm:"primaryKey"`
	UUID      string `gorm:"not null;type:varchar(50)"`
	Name      string `gorm:"not null;unique;type:varchar(50)"`
	ImageURL  string `gorm:"not null;type:varchar"`
	AdminID   int    `gorm:"not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("Before creating Product")

	if len(p.Name) < 4 {
		err = errors.New("name is too short")
	}

	return
}
