package user

import "gorm.io/gorm"

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(name string, age uint) (*User, error) {
	user := &User{
		Name: name,
		Age:  age,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Get(id uint) (*User, error) {
	var user User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) Update(id uint, name string, age uint) (*User, error) {
	user, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	user.Name = name
	user.Age = age

	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
