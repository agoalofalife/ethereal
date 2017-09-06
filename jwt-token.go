package ethereal

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

type EtherealClaims struct {
	jwt.StandardClaims
}

// get key jwt
func JWTKEY() []byte {
	return []byte(config("AUTH.JWT_KEY_HMAC").(string))
}

// handler check error
func handlerErrorToken(err error) (message error) {
	var locale =  config("L18N.LOCALE").(string)

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return errors.New(string(ConstructorI18N().T(locale, "jwtToken.ValidationErrorMalformed")))
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return errors.New(string(ConstructorI18N().T(locale, "jwtToken.ValidationErrorExpired")))
		} else {
			return errors.New(string(ConstructorI18N().T(locale, "jwtToken.ErrorBase")) + err.Error())
		}
	} else {
		return errors.New(string(ConstructorI18N().T(locale, "jwtToken.ErrorBase")) + err.Error())
	}
	return
}

func compareToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JWTKEY(), nil
	})
}
