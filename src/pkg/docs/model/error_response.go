package model

type ErrorResponse struct {
	Message string  `json:"message"`
	Errors  []error `json:"validationErrors"`
}

func (e *ErrorResponse) SetMessage(m string) *ErrorResponse {
	e.Message = m
	return e
}

func (e *ErrorResponse) AddError(err error) *ErrorResponse {
	e.Errors = append(e.Errors, err)
	return e
}
