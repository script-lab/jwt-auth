package model

import (
	"gorm.io/gorm"
	"time"
)

type (
	User struct {
		ID        uint   `json:"id" param:"id" gorm:"primaryKey"`
		Name      string `json:"name" gorm:"not null"`
		Email     string `json:"email" gorm:"unique" gorm:"not null"`
		Password  []byte `json:"password" gorm:"not null"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}
)
