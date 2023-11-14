package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Variant struct {
	ID          uint   `gorm:"primaryKey"`
	UUID        string `gorm:"not null" json:"uuid"`
	VariantName string `gorm:"not null" json:"variant_name" form:"title" valid:"required~Name of varient is required"`
	Quantity    int    `gorm:"not null" json:"quantity" form:"quantity" valid:"required~Quantity of variant is required, numeric~Invalid quantity format"`
	ProductID   uint
	Product     *Product
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func (v *Variant) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(v)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}

func (v *Variant) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(v)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}
