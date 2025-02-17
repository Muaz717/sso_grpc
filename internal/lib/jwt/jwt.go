package jwt

import (
	"fmt"
	"time"

	"github.com/Muaz717/sso/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	const op = "jwt.NewToken"

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.Id

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, err
}
