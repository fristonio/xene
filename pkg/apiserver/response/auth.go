package response

import "time"

// OauthLogin is the reponse for a user login to xene.
type OauthLogin struct {
	// LoginURL is the URL to be used for logging in.
	LoginURL string `json:"loginURL" example:"https://xxxx.io/login"`
}

// JWTAuth is the reponse for JWT authentication if succedded.
// It returns the token along with other relevant information in the json
// structure.
type JWTAuth struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`

	ExpiresIn time.Duration `json:"expiresIn" example:"3600"`

	UserEmail string `json:"userEmail" example:"example@example.com"`
	UserName  string `json:"userName" example:"fristonio"`
}
