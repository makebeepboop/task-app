package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/makebeepboop/task-app/internal/domain/models"
	"time"
)

func NewToken(app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	//claims["user_id"] = User.ID
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	// TODO: replace secret from structure App for more secure
	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
