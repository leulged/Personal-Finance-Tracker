package usecase

import (
    "errors"
    "strings"
    "time"
    
    "personal-finance-tracker/domain/entities"
    repoInterface "personal-finance-tracker/domain/interface"
    "personal-finance-tracker/Infrastructure/utils"
    "personal-finance-tracker/Infrastructure/service"
)

type UserUsecase struct{
    userRepo repoInterface.UserRepository
    passwordService *services.BcryptService
}

func NewUserUsecase(userRepo repoInterface.UserRepository) *UserUsecase{
	return &UserUsecase{
		userRepo: userRepo,
		passwordService: services.NewBcryptService(12), // Default bcrypt cost
	}
}

type UserInterface interface{
    Register(user *entities.User) (*entities.User, error)
}

func (u *UserUsecase) Register(user *entities.User) (*entities.User, error){
	if err := utils.IsEmailValid(user.Email); err != nil {
		return nil, err
	}
	
	if err := utils.IsValidPassword(user.Password); err != nil {
		return nil, err
	}

	user.Email = strings.ToLower(user.Email)
	existingUser, err := u.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		if err.Error() == "user not found" {
			// User doesn't exist, we can proceed with registration
		} else {
			return nil, errors.New("Database error: " + err.Error())
		}
	} else if existingUser != nil {
		return nil, errors.New("Email already in use")
	}

	// Check if this is the first user (admin)
	userCount, err := u.userRepo.CountUsers()
	if err != nil {
		return nil, errors.New("Database error while counting users: " + err.Error())
	}

	// Set role based on user count
	if userCount == 0 {
		user.Role = "admin" // First user becomes admin
	} else {
		user.Role = "user" // All other users are regular users
	}

	// Set verification status (default to false for new registrations)
	user.IsVerified = false

	// Hash password
	hashedPassword, err := u.passwordService.HashPassword(user.Password)
	if err != nil {
		return nil, errors.New("Failed to hash password")
	}
	user.Password = hashedPassword

	// Set timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Create user
	createUser, err := u.userRepo.CreateUser(user)
	if err != nil {
		return nil, errors.New("Failed to create user: " + err.Error())
	}

	// Don't return password in response
	createUser.Password = ""
	return createUser, nil
}