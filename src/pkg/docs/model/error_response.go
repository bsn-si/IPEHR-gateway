package model

type ErrorResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"validationErrors"`
}

func (e *ErrorResponse) Add(err error) {
	e.Errors = append(e.Errors, err.Error())

	//for err != nil {
	//	if x, ok := err.(interface{ As(any) bool }); ok && x.As(target) {
	//		return true
	//	}
	//	err = Unwrap(err)
	//}
}
