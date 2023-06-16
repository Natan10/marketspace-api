package services

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/natan10/marketspace-api/models"
)

type TokenService struct {
}

var (
	tokenFamily = "HS256"
	tokenSecret = "marketspace_secret"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	TokenAuth = NewToken()
}

func NewToken() *jwtauth.JWTAuth {
	return jwtauth.New(tokenFamily, []byte(tokenSecret), nil)
}

func (tk TokenService) EncodeToken(user *models.User) (t jwt.Token, tokenString string, err error) {
	claims := map[string]interface{}{
		"sub":      user.Id,
		"username": user.Name,
		"email":    user.Email,
		"exp":      time.Now().Add(2 * time.Hour),
	}

	return TokenAuth.Encode(claims)
}
