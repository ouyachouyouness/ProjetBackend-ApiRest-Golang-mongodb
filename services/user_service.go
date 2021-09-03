package services

import (
	"errors"
	"fmt"
	"goproj/models"
	"goproj/repositories"
	"strings"
	"time"
	// "github.com/google/uuid"
)

type UserService struct {
	userRepository repositories.UserRepository
}

type UserServiceContract interface {
	SignUp(body models.SignUpBody) (string, error)
	IsValidPassword(password string) bool
	GetEncryptedPassword(password string) string
	SignIn(body models.SignInBody) (string, error)
}

///////////////////////////////////NewInstanceOfUserService//////////////////////////////////////////////////

func NewInstanceOfUserService(userRepository repositories.UserRepository) UserService {
	return UserService{userRepository: userRepository}
}

///////////////////////////////////SignUp//////////////////////////////////////////////////

func (u *UserService) SignUp(body models.SignUpBody) (string, error) {
	emailLowerCase := strings.ToLower(body.Email)
	emailTrimmed := strings.Trim(emailLowerCase, " ")

	// Verify password meets sign up requirements
	if !u.IsValidPassword(body.Password) {
		return "", errors.New("error: Your password does not meet requirements.")
	}

	// Check for user
	userExists, err := u.userRepository.DoesUserExist(emailTrimmed)
	if err != nil {
		return "", err
	}
	if userExists {

		return u.signIn(emailTrimmed, body.Password)
	}

	newUser := models.User{
		Email:    emailTrimmed,
		Password: u.GetEncryptedPassword(body.Password),
		Name:     body.Name,
		Created:  time.Now(),
	}
	err = u.userRepository.SaveUser(newUser)
	if err != nil {
		return "", err
	}

	return u.signIn(emailTrimmed, body.Password)
}

///////////////////////////////////IsValidPassword//////////////////////////////////////////////////

func (u *UserService) IsValidPassword(password string) bool {

	fmt.Println("---- IMPORTANT ----")
	fmt.Println("Your sign up validate is failing .")

	return true
}

///////////////////////////////////GetEncryptedPassword//////////////////////////////////////////////////

func (u *UserService) GetEncryptedPassword(password string) string {
	// TODO - See comment in the isValidPassword
	return password
}

///////////////////////////////////SignIn//////////////////////////////////////////////////

func (u *UserService) SignIn(body models.SignInBody) (string, error) {
	emailLowerCase := strings.ToLower(body.Email)
	emailTrimmed := strings.Trim(emailLowerCase, " ")
	return u.signIn(emailTrimmed, body.Password)
}

///////////////////////////////////signIn//////////////////////////////////////////////////

func (u *UserService) signIn(email string, password string) (string, error) {
	// Encrypt password
	encryptedPassword := u.GetEncryptedPassword(password)

	// Grab user
	found, user, err := u.userRepository.GetUserByEmail(email)
	if err != nil {
		return "", err
	}
	if !found {
		return "", errors.New("error: Unauthorized")
	}

	if user.Password != encryptedPassword {
		return "", errors.New("error: Unauthorized")
	}

	// Create session
	now := time.Now()
	expiryDate := now.AddDate(0, 0, 1)
	newSession := models.Session{
		Email:   email,
		Created: now,
		Expiry:  expiryDate,
	}

	token, err := u.userRepository.SaveSession(newSession)
	if err != nil {
		return "", err
	}
	return token, nil
}

///////////////////////////////////Update//////////////////////////////////////////////////

func (c *UserService) Update(session models.Session, carID string, body models.UpdateUser) error {

	err := c.userRepository.Update(session.Email, carID, body)
	if err != nil {
		return err
	}
	return nil
}

///////////////////////////////////Delete//////////////////////////////////////////////////

func (c *UserService) Delete(session models.Session, carID string) error {

	err := c.userRepository.Delete(session.Email, carID)
	if err != nil {
		return err
	}
	return nil
}
