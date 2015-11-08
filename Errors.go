package main

//ErrorCode ...
type ErrorCode struct {
	Code    int
	Message string
}

func configureErrorMappings() {
	ErrorMappings = map[string]ErrorCode{}

	ErrorMappings["socket: too many open files"] = ErrorCode{
		Code:    5001,
		Message: "Your workers value is set to high.  Either increase the system limit for open files or reduce the value of the workers",
	}
	ErrorMappings[`unsupported protocol scheme ""`] = ErrorCode{
		Code:    5002,
		Message: "Your urls in the test specification must be valid urls",
	}
}
