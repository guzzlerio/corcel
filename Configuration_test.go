package main

import (
	//"fmt"
    "time"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration", func() {

	var configuration *Configuration

	BeforeEach(func() {
		configuration = ParseConfiguration()
	})

	Describe("When no config file is found and no command line args are provided", func() {
		Describe("Loading a default configuration", func() {
            It("sets random (--random)", func() {
				Expect(configuration.Random).To(Equal(false))
            })
            It("sets summary (--summary)", func() {
				Expect(configuration.Summary).To(Equal(false))
            })
			It("sets workers (--workers)", func() {
				Expect(configuration.Workers).To(Equal(int64(1)))
			})
            It("sets wait-time (--wait-time)", func() {
                defaultWaitTime := time.Duration(0)
                Expect(configuration.WaitTime).To(Equal(defaultWaitTime))
            })
            It("sets duration (--duration)", func() {
                defaultDuration := time.Duration(0)
                Expect(configuration.WaitTime).To(Equal(defaultDuration))
            })
		})
	})

	Describe("When config file is not found in pwd", func() {
		Describe("and config file is found in user home", func() {
		})
	})

	Describe("When config file is found in pwd", func() {
		Describe("and config file is found in user home", func() {
		})
	})

	Describe("When commandline args are provided", func() {
		Describe("overriding the default configuration", func() {
			Describe("for workers (-w)", func() {
			})
		})
	})

	/*
		It("Lex a line twice", func() {
			result := lexer.Lex(line)
			Expect(result).To(Equal([]string{"http://127.0.0.1:8000/A", "-X", method, "-H", applicationJson, "-H", applicationSoapXml}))

			result = lexer.Lex(line)
			Expect(result).To(Equal([]string{"http://127.0.0.1:8000/A", "-X", method, "-H", applicationJson, "-H", applicationSoapXml}))
		})
	*/
})
