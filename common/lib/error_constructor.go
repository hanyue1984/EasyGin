package lib

// CustomError 自定义类型
type HTTPError struct {
	Code    int
	Message string
}

func CustomError(code int, message string) {
	err := HTTPError{
		Code:    code,
		Message: message,
	}
	panic(err)
}
