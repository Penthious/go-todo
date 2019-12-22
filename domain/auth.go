package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
)

type RegisterPayload struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Username        string `json:"username"`
}

func (r *RegisterPayload) IsValid() (bool, map[string]string) {
	v := NewValidator()

	v.ValueIsRequired("email", r.Email)
	v.MustBeValidEmail("email", r.Email)

	v.ValueIsRequired("password", r.Password)
	v.MustBeLongerThan("password", r.Password, 6)

	v.ValueIsRequired("confirmPassword", r.ConfirmPassword)
	v.MustMatch(
		ElementMatcher{field: "password", value: r.Password},
		ElementMatcher{field: "confirmPassword", value: r.ConfirmPassword},
	)

	v.ValueIsRequired("username", r.Username)
	v.MustBeLongerThan("username", r.Username, 3)

	return v.IsValid(), v.errors
}

func (d *Domain) Register(payload RegisterPayload) (*User, error) {
	userExist, _ := d.DB.UserRepo.GetByEmail(payload.Email)

	if userExist != nil {
		return nil, ErrUserWithEmailAlreadyExist
	}

	userExist, _ = d.DB.UserRepo.GetByUsername(payload.Username)

	if userExist != nil {
		return nil, ErrUserWithUsernameAlreadyExist
	}

	password, err := d.setPassword(payload.Password)

	if err != nil {
		return nil, err
	}

	data := &User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: *password,
	}

	user, err := d.DB.UserRepo.Create(data)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (d *Domain) setPassword(password string) (*string, error) {
	passwordByte := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}
	password = string(passwordHash)

	return &password, nil
}

func stripBearerPrefixFromToken(token string) (string, error) {
	bearer := "BEARER"

	if len(token) > len(bearer) && strings.ToUpper(token[0:len(bearer)]) == bearer {

		return token[len(bearer)+1:], nil
	}

	// @todo: pass an error instead
	return token, nil

}

var authHeaderExtractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{"Authorization"},
	Filter:    stripBearerPrefixFromToken,
}

var authExtractor = &request.MultiExtractor{
	authHeaderExtractor,
}

func ParseToken(r *http.Request) (*jwt.Token, error) {

	token, error := request.ParseFromRequest(r, authExtractor, func(t *jwt.Token) (interface{}, error) {
		// @todo: what happens if no jwt_secret is set?
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	return token, error
}
