package auth

import (
	"time"

	"github.com/JesseNicholas00/HaloSuster/repos/auth"
	"github.com/golang-jwt/jwt/v4"
)

func (svc *authServiceImpl) generateToken(
	user auth.User,
) (res string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
		Data: jwtSubClaims{
			UserId: user.Id,
			Nip:    user.Nip,
		},
	})
	res, err = token.SignedString(svc.jwtSecret)
	return
}
