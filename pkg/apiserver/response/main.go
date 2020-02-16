package response

// HTTPError is the response type for any, general simple error that occur for a HTTP request.
type HTTPError struct {
	Error string `json:"error" example:"Invalid authentication type provided."`
}
