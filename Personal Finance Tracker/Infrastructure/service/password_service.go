package services

import (
    "golang.org/x/crypto/bcrypt"
)

type BcryptService struct {
    cost int    
}

func NewBcryptService(cost int) *BcryptService {
    return &BcryptService{cost: cost}
}

func (s *BcryptService) HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

func (s *BcryptService) ComparePassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
