package main

import (
	//"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Configuration", func() {

	var configuration *Configuration
	var args []string

	BeforeEach(func() {
        args = []string{"file-path.yml"}
		configuration = ParseConfiguration(args)
	})

	Describe("When no config file is found and no command line args are provided", func() {
		Describe("Loading a default configuration", func() {
			It("sets duration (--duration)", func() {
				defaultDuration := time.Duration(0)
				Expect(configuration.WaitTime).To(Equal(defaultDuration))
			})
			It("sets random (--random)", func() {
				Expect(configuration.Random).To(Equal(false))
			})
			It("sets summary (--summary)", func() {
				Expect(configuration.Summary).To(Equal(false))
			})
			It("sets workers (--workers)", func() {
				Expect(configuration.Workers).To(Equal(1))
			})
			It("sets wait-time (--wait-time)", func() {
				defaultWaitTime := time.Duration(0)
				Expect(configuration.WaitTime).To(Equal(defaultWaitTime))
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
			Describe("for duration (--duration)", func() {
				BeforeEach(func() {
					args = []string{"--duration", "3s", "./path/to/file"}
					configuration = ParseConfiguration(args)
				})
				It("sets the value", func() {
					duration, _ := time.ParseDuration("3s")
					Expect(configuration.Duration).To(Equal(duration))
				})
			})

			Describe("for file", func() {
				BeforeEach(func() {
					args = []string{"./path/to/file"}
					configuration = ParseConfiguration(args)
				})
				It("sets the value", func() {
					Expect(configuration.FilePath).To(Equal("./path/to/file"))
				})
			})

			Describe("for random (--random)", func() {
				BeforeEach(func() {
					args = []string{"--random", "./path/to/file"}
					configuration = ParseConfiguration(args)
				})
				It("sets the value", func() {
					Expect(configuration.Random).To(Equal(true))
				})
			})

			Describe("for summary (--summary)", func() {
				BeforeEach(func() {
					args = []string{"--summary", "./path/to/file"}
					configuration = ParseConfiguration(args)
				})
				It("sets the value", func() {
					Expect(configuration.Summary).To(Equal(true))
				})
			})

			Describe("for workers (--workers)", func() {
				BeforeEach(func() {
					args = []string{"--workers", "3", "./path/to/file"}
					configuration = ParseConfiguration(args)
				})
				It("sets the value", func() {
					Expect(configuration.Workers).To(Equal(3))
				})
			})

			Describe("for wait-time (--wait-time)", func() {
				BeforeEach(func() {
					args = []string{"--wait-time", "3s", "./path/to/file"}
					configuration = ParseConfiguration(args)
				})
				It("sets the value", func() {
					waitTime, _ := time.ParseDuration("3s")
					Expect(configuration.WaitTime).To(Equal(waitTime))
				})
			})
		})
	})
})
