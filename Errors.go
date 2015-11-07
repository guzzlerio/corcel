package main

type ErrorCode struct{
    Code int
    Message string
}

func ConfigureErrorMappings(){
    ErrorMappings = map[string]ErrorCode{}

    ErrorMappings["socket: too many open files"] = ErrorCode{
        Code : 5001,
        Message : "Your workers value is set to high.  Either increase the system limit for open files or reduce the value of the workers",
    }
}
