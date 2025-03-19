package dtos

type BaseResponse struct {
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Errors  []Error     `json:"errors,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

type Error struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func WithSuccess(message string, code int, result interface{}) *BaseResponse {
	return &BaseResponse{
		Message: message,
		Success: true,
		Code:    code,
		Result:  result,
	}
}

func WithError(message string, code int, errors ...Error) *BaseResponse {
	return &BaseResponse{
		Message: message,
		Success: false,
		Code:    code,
		Errors:  errors,
	}
}
