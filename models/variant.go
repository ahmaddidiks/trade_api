package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Variant struct {
	ID          int    `gorm:"primaryKey"`
	UUID        string `gorm:"not null;type:varchar(50)"`
	VariantName string `gorm:"not null;type:varchar(50)"`
	Quantity    int    `gorm:"not null"`
	ProductID   int    `gorm:"not null"`
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

func (v *Variant) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("Before creating Variant")

	if len(v.VariantName) < 4 {
		err = errors.New("name is too short")
	}

	return
}
