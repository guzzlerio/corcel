package request

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCommand_Line_Lexer(t *testing.T) {
	BeforeTest()

	defer AfterTest()
	Convey("Command Line Lexer", t, func() {

		var lexer *CommandLineLexer

		func() {
			lexer = NewCommandLineLexer()
		}()

		Convey("Lex a line twice", func() {
			method := "GET"
			applicationJSON := "Content-type:application/json"
			applicationSoapXML := "Content-type:application/soap+xml"
			line := fmt.Sprintf("http://127.0.0.1:8000/A -X %s -H \"%s\" -H \"%s\"", method, applicationJSON, applicationSoapXML)

			result := lexer.Lex(line)
			So(result, ShouldResemble, []string{"http://127.0.0.1:8000/A", "-X", method, "-H", applicationJSON, "-H", applicationSoapXML})

			result = lexer.Lex(line)
			So(result, ShouldResemble, []string{"http://127.0.0.1:8000/A", "-X", method, "-H", applicationJSON, "-H", applicationSoapXML})

		})

		Convey("Lex one single dashed flag", func() {
			result := lexer.Lex("-X POST")
			So(result, ShouldResemble, []string{"-X", "POST"})
		})

		Convey("Lex a flag value which contains a space", func() {
			result := lexer.Lex("-H \"Content-type: application/json\" -H \"Accept: text/xml\"")
			So(result, ShouldResemble, []string{"-H", "Content-type: application/json", "-H", "Accept: text/xml"})
		})

		Convey("Lex a flag value which is inside single quotes", func() {
			result := lexer.Lex("-H 'Content-type: application/json' -H 'Accept: text/xml'")
			So(result, ShouldResemble, []string{"-H", "Content-type: application/json", "-H", "Accept: text/xml"})
		})

		Convey("Lex a flag value which contains a non-flag at beginning", func() {
			result := lexer.Lex(" A -H \"Content-type: application/json\"")
			So(result, ShouldResemble, []string{"A", "-H", "Content-type: application/json"})
		})

		Convey("Lex a flag value with double quotes inside single quotes", func() {
			result := lexer.Lex("-d '{\"name\":\"bob\"}'")
			So(result, ShouldResemble, []string{"-d", "{\"name\":\"bob\"}"})
		})

		Convey("Lex a flag value with single quotes inside double quotes", func() {
			result := lexer.Lex("-d \"{'name':'bob'}\"")
			So(result, ShouldResemble, []string{"-d", "{'name':'bob'}"})
		})

		Convey("Lex a double dashed flag", func() {
			result := lexer.Lex("--cert ca.pem")
			So(result, ShouldResemble, []string{"--cert", "ca.pem"})
		})

		Convey("Lex a double dashed flag without value", func() {
			result := lexer.Lex("--dummy --cert ca.pem")
			So(result, ShouldResemble, []string{"--dummy", "--cert", "ca.pem"})
		})

		Convey("Lex extra spaces between flags", func() {
			result := lexer.Lex("       --dummy     --cert     ca.pem    -X     POST    ")
			So(result, ShouldResemble, []string{"--dummy", "--cert", "ca.pem", "-X", "POST"})
		})

		Convey("Lex commands amongst flags", func() {
			result := lexer.Lex("--dummy command --cert ca.pem")
			So(result, ShouldResemble, []string{"--dummy", "command", "--cert", "ca.pem"})
		})

		Convey("Lex flags with a colon separator between flag and value", func() {
			result := lexer.Lex("-X:POST --cert:ca.pem")
			So(result, ShouldResemble, []string{"-X", "POST", "--cert", "ca.pem"})
		})
	})
}
