package authentication

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jna-distribution/service-shipping/internal/entity"
	"github.com/jna-distribution/service-shipping/pkg"
	"time"
)

// CreateToken return new token. Use a little variable for compact.
func CreateToken(userId string, userRole entity.Role, contactId int16) (string, error) {
	env, err := pkg.LoadEnv()
	if err != nil {
		return "", err
	}
	signingKey := []byte(env["JWT_SIGNING_KEY"])
	claims := Claims{
		UserId:   userId,
		UserRole: userRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(2)).Unix(),
			Issuer:    "JNA Distribution",
		},
	}
	if userRole == entity.RoleShop {
		claims.ContactId = contactId
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return ss, nil
}
