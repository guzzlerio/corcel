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
	defaultWaitTime := time.Duration(0)
	defaultDuration := time.Duration(0)

	BeforeEach(func() {
		args = []string{"file-path.yml"}
		configuration = ParseConfiguration(args)
	})

	Describe("When no config file is found and no command line args are provided", func() {
		Describe("Loading a default configuration", func() {
			It("sets duration (--duration)", func() {
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
				It("applies the override", func() {
					duration, _ := time.ParseDuration("3s")
					Expect(configuration.Duration).To(Equal(duration))
				})

				Describe("leaves the default for", func() {
					It("random", func() {
						Expect(configuration.Random).To(Equal(false))
					})
					It("summary", func() {
						Expect(configuration.Summary).To(Equal(false))
					})
					It("workers", func() {
						Expect(configuration.Workers).To(Equal(1))
					})
					It("wait-time", func() {
						Expect(configuration.WaitTime).To(Equal(defaultWaitTime))
					})
				})
			})

			Describe("for file", func() {
				BeforeEach(func() {
					args = []string{"./path/to/file"}
					configuration = ParseConfiguration(args)
				})
				It("applies the override", func() {
					Expect(configuration.FilePath).To(Equal("./path/to/file"))
				})

				Describe("leaves the default for", func() {
					It("duration", func() {
						Expect(configuration.Duration).To(Equal(defaultDuration))
					})
					It("random", func() {
						Expect(configuration.Random).To(Equal(false))
					})
					It("summary", func() {
						Expect(configuration.Summary).To(Equal(false))
					})
					It("workers", func() {
						Expect(configuration.Workers).To(Equal(1))
					})
					It("wait-time", func() {
						Expect(configuration.WaitTime).To(Equal(defaultWaitTime))
					})
				})
			})

			Describe("for random (--random)", func() {
				BeforeEach(func() {
					args = []string{"--random", "./path/to/file"}
					configuration = ParseConfiguration(args)
				})
				It("applies the override", func() {
					Expect(configuration.Random).To(Equal(true))
				})

				Describe("leaves the default for", func() {
					It("duration", func() {
						Expect(configuration.Duration).To(Equal(defaultDuration))
					})
					It("summary", func() {
						Expect(configuration.Summary).To(Equal(false))
					})
					It("workers", func() {
						Expect(configuration.Workers).To(Equal(1))
					})
					It("wait-time", func() {
						Expect(configuration.WaitTime).To(Equal(defaultWaitTime))
					})
				})
			})

			Describe("for summary (--summary)", func() {
				BeforeEach(func() {
					args = []string{"--summary", "./path/to/file"}
					configuration = ParseConfiguration(args)
				})
				It("applies the override", func() {
					Expect(configuration.Summary).To(Equal(true))
				})

				Describe("leaves the default for", func() {
					It("duration", func() {
						Expect(configuration.Duration).To(Equal(defaultDuration))
					})
					It("random", func() {
						Expect(configuration.Random).To(Equal(false))
					})
					It("workers", func() {
						Expect(configuration.Workers).To(Equal(1))
					})
					It("wait-time", func() {
						Expect(configuration.WaitTime).To(Equal(defaultWaitTime))
					})
				})
			})

			Describe("for workers (--workers)", func() {
				BeforeEach(func() {
					args = []string{"--workers", "3", "./path/to/file"}
					configuration = ParseConfiguration(args)
				})
				It("applies the override", func() {
					Expect(configuration.Workers).To(Equal(3))
				})

				Describe("leaves the default for", func() {
					It("duration", func() {
						Expect(configuration.Duration).To(Equal(defaultDuration))
					})
					It("random", func() {
						Expect(configuration.Random).To(Equal(false))
					})
					It("summary", func() {
						Expect(configuration.Summary).To(Equal(false))
					})
					It("wait-time", func() {
						Expect(configuration.WaitTime).To(Equal(defaultWaitTime))
					})
				})
			})

			Describe("for wait-time (--wait-time)", func() {
				BeforeEach(func() {
					args = []string{"--wait-time", "3s", "./path/to/file"}
					configuration = ParseConfiguration(args)
				})
				It("applies the override", func() {
					waitTime, _ := time.ParseDuration("3s")
					Expect(configuration.WaitTime).To(Equal(waitTime))
				})

				Describe("leaves the default for", func() {
					It("duration", func() {
						Expect(configuration.Duration).To(Equal(defaultDuration))
					})
					It("random", func() {
						Expect(configuration.Random).To(Equal(false))
					})
					It("summary", func() {
						Expect(configuration.Summary).To(Equal(false))
					})
					It("workers", func() {
						Expect(configuration.Workers).To(Equal(1))
					})
				})
			})
		})
	})
})
