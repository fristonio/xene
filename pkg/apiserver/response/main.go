package response

// HTTPMessage is the general HTTP response for a request.
type HTTPMessage struct {
	Message string `json:"message" example:"Messsage in response to your request"`
}

// HTTPError is the response type for any, general simple error that occur for a HTTP request.
type HTTPError struct {
	Error string `json:"error" example:"Invalid authentication type provided."`
}
