package domain

type APIErrorResponse struct{
	Error APIError `json:"error"`
	
}
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	
}
