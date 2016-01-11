package errormanager

import (
	"fmt"
	"os"
	"strings"

	"ci.guzzler.io/guzzler/corcel/logger"
)

//Check ...
func Check(err error) {
	if err != nil {
		Log(err)
	}
}

//ErrorCode ...
type ErrorCode struct {
	Code    int
	Message string
}

//ErrorManager ...
var mappings map[string]ErrorCode

//New ...
func configure() {
	mappings = make(map[string]ErrorCode)

	mappings["socket: too many open files"] = ErrorCode{
		Code:    5001,
		Message: "Your workers value is set to high.  Either increase the system limit for open files or reduce the value of the workers",
	}
	mappings[`unsupported protocol scheme ""`] = ErrorCode{
		Code:    5002,
		Message: "Your urls in the test specification must be valid urls",
	}
	mappings["too many open files"] = ErrorCode{
		Code:    5003,
		Message: "Too many workers man!",
	}
}

//Log ...
func Log(err interface{}) {
	if mappings == nil {
		configure()
	}
	for mapping, errorCode := range mappings {
		if strings.Contains(fmt.Sprintf("%v", err), mapping) {
			fmt.Println(errorCode.Message)
			os.Exit(errorCode.Code)
		}
	}
	logger.Log.Fatalf("UNKNOWN ERROR: %v", err)
}
