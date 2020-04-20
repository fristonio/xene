package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fristonio/xene/pkg/defaults"
)

var (
	// ErrTokenExpired represents the error if the JWT token is expired.
	ErrTokenExpired = errors.New("Token is either expired or not active yet")
)

// UserClaims are the claims that are verified in the jwt token.
// Email is the email of the user associated with the claim.
// The Claim type also includes default jwt standard claims.
type UserClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`

	// Roles contains the roles assumed by the user.
	Roles []string `json:"roles"`

	jwt.StandardClaims
}

// ServerClaims are the claims that a server uses to establish its identity
// while communication with a peer server.
type ServerClaims struct {
	Address string `json:"address"`
	Name    string `json:"name"`

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
func (j *AuthProvider) NewAuthToken(email, name string) (string, error) {
	expirationTime := time.Now().Add(j.expireInterval)
	c := &UserClaims{
		Email: email,
		Name:  name,

		Roles: []string{"user"},

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(j.Secret)
}

// NewServerAuthToken issues a new JWT authentication token to be used for the provided
// server name and address.
func (j *AuthProvider) NewServerAuthToken(address, name string) (string, error) {
	expirationTime := time.Now().Add(j.expireInterval)
	c := &ServerClaims{
		Address: address,
		Name:    name,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(j.Secret)
}

// GetUserClaimsFromToken parses the token provided in the string, returns an error if there
// is any issue with the token else populate the claim object with the provided
// fields.
func (j *AuthProvider) GetUserClaimsFromToken(token string) (*UserClaims, error) {
	c := &UserClaims{}
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
			return c, ErrTokenExpired
		}
	}

	return c, errors.New("Clould not handle the token")
}

// GetServerClaimsFromToken parses the token against a server claims provided in the string, returns an error if there
// is any issue with the token else populate the claim object with the provided
// fields.
func (j *AuthProvider) GetServerClaimsFromToken(token string) (*ServerClaims, error) {
	c := &ServerClaims{}
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
			return c, ErrTokenExpired
		}
	}

	return c, errors.New("Clould not handle the token")
}
