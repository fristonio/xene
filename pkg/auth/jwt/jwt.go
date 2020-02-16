package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fristonio/xene/pkg/defaults"
)

// Claims are the claims that are verified in the jwt token.
// Email is the email of the user associated with the claim.
// The Claim type also includes default jwt standard claims.
type Claims struct {
	Email string `json:"email"`

	jwt.StandardClaims
}

// AuthProvider is the authentication provider for JWT, it contains
// functions to issue, validate and refresh tokens using the JWT
// secret embedded.
type AuthProvider struct {
	// Secret is the JWT secret used to issue the tokens.
	Secret         []byte
	expireInterval time.Duration
}

// NewJWTAuthProvider returns a new instance of Authentication provider
// with the specified secret used for all relevant purposes.
func NewJWTAuthProvider(secret string) *AuthProvider {
	return &AuthProvider{
		Secret:         []byte(secret),
		expireInterval: defaults.JWTExpireInterval,
	}
}

// NewAuthToken issues a new JWT authentication token to be used for the provided
// user email.
func (j *AuthProvider) NewAuthToken(email string) (string, error) {
	expirationTime := time.Now().Add(j.expireInterval)
	c := &Claims{
		Email: email,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(j.Secret)
}

// GetClaimsFromToken parses the token provided in the string, returns an error if there
// is any issue with the token else populate the claim object with the provided
// fields.
func (j *AuthProvider) GetClaimsFromToken(token string) (*Claims, error) {
	c := &Claims{}
	t, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		return j.Secret, nil
	})

	if t.Valid {
		return c, nil
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return c, errors.New("Malformed token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return c, errors.New("Token is either expired or not active yet")
		}
	}

	return c, errors.New("Clould not handle the token")
}