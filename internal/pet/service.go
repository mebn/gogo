package pet

import (
	"gogo/internal/user"

	"gorm.io/gorm"
)

type Pet struct {
	ID      uint      `json:"id" gorm:"primaryKey"`
	OwnerID uint      `json:"ownerId" gorm:"not null;index"`
	Owner   user.User `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Name    string    `json:"name"`
	Age     uint      `json:"age"`
}

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(ownerID uint, name string, age uint) (*Pet, error) {
	pet := &Pet{
		OwnerID: ownerID,
		Name:    name,
		Age:     age,
	}

	if err := s.db.Create(pet).Error; err != nil {
		return nil, err
	}

	return pet, nil
}

func (s *Service) Get(id uint) (*Pet, error) {
	var pet Pet
	if err := s.db.First(&pet, id).Error; err != nil {
		return nil, err
	}

	return &pet, nil
}

func (s *Service) Update(id uint, ownerID uint, name string, age uint) (*Pet, error) {
	pet, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	pet.OwnerID = ownerID
	pet.Name = name
	pet.Age = age

	if err := s.db.Save(pet).Error; err != nil {
		return nil, err
	}

	return pet, nil
}
