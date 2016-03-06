package httpext

// HttpError is a interface for errors that are transferable over HTTP as JSON.
type HttpError interface {
	StatusCode() int
	Error() string
	InnerError() error
}

// EncapulateIfError encapsulates the error provided as a HttpError if not nil.
// If the error provided is nil, this function returns nil.
func EncapulateIfError(err error, statusCode int) HttpError {
	if err != nil {
		return NewHttpError(err, statusCode)
	}
	return nil
}

// NewHttpError returns an encapsulated error suitable for JSON
// serialization. The error also has a HTTP status code for
// usability.
func NewHttpError(err error, statusCode int) HttpError {
	return &httpErrorImpl{statusCode, err, err.Error()}
}

type httpErrorImpl struct {
	Code         int    `json:"statusCode"`
	Err          error  `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

func (e httpErrorImpl) StatusCode() int {
	return e.Code
}

func (e httpErrorImpl) Error() string {
	return e.Err.Error()
}

func (e httpErrorImpl) InnerError() error {
	return e.Err
}
