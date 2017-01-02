package errormanager

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/satori/go.uuid"
)

var (
	panicNotRecover = false
)

//PanicNotRecover ...
func PanicNotRecover() {
	panicNotRecover = true
}

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

const (
	//LogMessageVaidURLs ...
	LogMessageVaidURLs = "Your urls in the test specification must be absolute and valid urls"
)

//New ...
func configure() {
	mappings = make(map[string]ErrorCode)

	mappings["socket: too many open files"] = ErrorCode{
		Code:    1,
		Message: "Your workers value is set to high.  Either increase the system limit for open files or reduce the value of the workers",
	}
	mappings[`unsupported protocol scheme ""`] = ErrorCode{
		Code:    2,
		Message: LogMessageVaidURLs,
	}
	mappings["too many open files"] = ErrorCode{
		Code:    3,
		Message: "Too many workers man!",
	}
	mappings["invalid URI for request"] = ErrorCode{
		Code:    4,
		Message: LogMessageVaidURLs,
	}
}

//HandlePanic ...
func HandlePanic() {
	if panicNotRecover {
		return
	}

	if err := recover(); err != nil { //catch

		for mapping, errorCode := range mappings {
			if strings.Contains(fmt.Sprintf("%v", err), mapping) {
				fmt.Println(errorCode.Message)
				os.Exit(errorCode.Code)
			}
		}

		var id = uuid.NewV4().String()
		ioutil.WriteFile(fmt.Sprintf("/tmp/%v", id), []byte(fmt.Sprintf("%v \n\n %s", err, string(debug.Stack()))), 0644)
		fmt.Println(fmt.Sprintf("An unexpected error has occurred.  The error has been logged to /tmp/%v", id))
		os.Exit(255)
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
			//		os.Exit(errorCode.Code)
		}
	}
	//panic(err)
	panic("BOOM")
	//logger.Log.Fatalf("UNKNOWN ERROR: %v", err)
}
