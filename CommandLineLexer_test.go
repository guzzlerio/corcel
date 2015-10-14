package main

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Command Line Lexer", func() {

	var lexer *CommandLineLexer

	BeforeEach(func(){
		lexer = NewCommandLineLexer()
	})

	It("Lex a line twice", func(){
			method := "GET"
			applicationJson := "Content-type:application/json"
			applicationSoapXml := "Content-type:application/soap+xml"
			line := fmt.Sprintf("http://127.0.0.1:8000/A -X %s -H \"%s\" -H \"%s\"", method, applicationJson, applicationSoapXml)

		result := lexer.Lex(line)
		Expect(result).To(Equal([]string{"http://127.0.0.1:8000/A", "-X", method, "-H", applicationJson, "-H", applicationSoapXml}))

		result = lexer.Lex(line)
		Expect(result).To(Equal([]string{"http://127.0.0.1:8000/A", "-X", method, "-H", applicationJson, "-H", applicationSoapXml}))

	})

	It("Lex one single dashed flag", func(){
		result := lexer.Lex("-X POST")
		Expect(result).To(Equal([]string{"-X","POST"}))
	})

	It("Lex a flag value which contains a space", func(){
		result := lexer.Lex("-H \"Content-type: application/json\" -H \"Accept: text/xml\"")
		Expect(result).To(Equal([]string{"-H", "Content-type: application/json","-H","Accept: text/xml"}))
	})

	It("Lex a flag value which is inside single quotes", func(){
		result := lexer.Lex("-H 'Content-type: application/json' -H 'Accept: text/xml'")
		Expect(result).To(Equal([]string{"-H", "Content-type: application/json","-H","Accept: text/xml"}))
	})

	It("Lex a flag value which contains a non-flag at beginning", func(){
		result := lexer.Lex(" A -H \"Content-type: application/json\"")
		Expect(result).To(Equal([]string{"A", "-H", "Content-type: application/json"}))
	})

	It("Lex a flag value with double quotes inside single quotes",func(){
		result := lexer.Lex("-d '{\"name\":\"bob\"}'")
		Expect(result).To(Equal([]string{"-d", "{\"name\":\"bob\"}"}))
	})

	It("Lex a flag value with single quotes inside double quotes",func(){
		result := lexer.Lex("-d \"{'name':'bob'}\"")
		Expect(result).To(Equal([]string{"-d", "{'name':'bob'}"}))
	})

	It("Lex a double dashed flag", func(){
		result := lexer.Lex("--cert ca.pem")
		Expect(result).To(Equal([]string{"--cert","ca.pem"}))
	})

	It("Lex a double dashed flag without value", func(){
		result := lexer.Lex("--dummy --cert ca.pem")
		Expect(result).To(Equal([]string{"--dummy","--cert","ca.pem"}))
	})

	It("Lex extra spaces between flags", func(){
		result := lexer.Lex("       --dummy     --cert     ca.pem    -X     POST    ")
		Expect(result).To(Equal([]string{"--dummy","--cert","ca.pem","-X","POST"}))
	})

	It("Lex commands amongst flags", func(){
		result := lexer.Lex("--dummy command --cert ca.pem")
		Expect(result).To(Equal([]string{"--dummy","command","--cert","ca.pem"}))
	})

	It("Lex flags with a colon separator between flag and value", func(){
		result := lexer.Lex("-X:POST --cert:ca.pem")
		Expect(result).To(Equal([]string{"-X","POST","--cert","ca.pem"}))
	})
})
