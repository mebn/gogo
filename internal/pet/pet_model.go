package pet

import "gogo/internal/user"

type Pet struct {
	ID      uint      `json:"id" gorm:"primaryKey"`
	OwnerID uint      `json:"ownerId" gorm:"not null;index"`
	Owner   user.User `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Name    string    `json:"name"`
	Age     uint      `json:"age"`
}
