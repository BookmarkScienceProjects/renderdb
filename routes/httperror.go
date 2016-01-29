package routes

type HttpError interface {
	StatusCode() int
	Error() error
}

func EncapulateIfError(err error, statusCode int) HttpError {
	if err != nil {
		return NewHttpError(err, statusCode)
	}
	return nil
}

func NewHttpError(err error, statusCode int) HttpError {
	return &MyHtpError{statusCode, err}
}

type MyHtpError struct {
	Code int   `json:"statusCode"`
	Err  error `json:"error"`
}

func (e MyHtpError) StatusCode() int {
	return e.Code
}

func (e MyHtpError) Error() error {
	return e.Err
}
