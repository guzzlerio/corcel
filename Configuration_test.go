package main

import (
	"fmt"
    "runtime"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
	//"errors"
)

var _ = Describe("Configuration", func() {

	var configuration *Configuration
	var args []string
	defaultWaitTime := time.Duration(0)
	defaultDuration := time.Duration(0)
	_, filename, _, _ := runtime.Caller(1)

	BeforeEach(func() {
		args = []string{filename}
		configFileReader = func(path string) ([]byte, error) {
			fmt.Println("test filereader")
			return []byte(""), nil
		}
		configuration, _ = ParseConfiguration(args)
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
		It("returns error", func() {
			/*
				var t Configuration
				yaml.Unmarshal([]byte("workers: 3"), &t)
				fmt.Printf("%+v", t)
				Expect(t.Workers).To(Equal(2))
			*/
			yaml := "workers: 3"

			configFileReader = func(path string) ([]byte, error) {
				return []byte(yaml), nil
			}
			configuration, err := ParseConfiguration(args)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(configuration.Workers).To(Equal(3))
		})

		Describe("and config file is found in user home", func() {
		})
	})

	Describe("When commandline args are provided", func() {
		Describe("overriding the default configuration", func() {
			Describe("for duration (--duration)", func() {
				BeforeEach(func() {
					args = []string{"--duration", "3s", filename}
					configuration, _ = ParseConfiguration(args)
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
					configuration, _ = ParseConfiguration(args)
				})
				It("applies the override", func() {
					Expect(configuration.FilePath).To(Equal(filename))
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
					args = []string{"--random", filename}
					configuration, _ = ParseConfiguration(args)
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
					args = []string{"--summary", filename}
					configuration, _ = ParseConfiguration(args)
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
					args = []string{"--workers", "3", filename}
					configuration, _ = ParseConfiguration(args)
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
					args = []string{"--wait-time", "3s", filename}
					configuration, _ = ParseConfiguration(args)
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

		Describe("setting multiple command line args", func() {
			BeforeEach(func() {
				args = []string{"--summary", "--workers", "50", "--duration", "3s", filename}
				configuration, _ = ParseConfiguration(args)
			})
			It("applies the overrides", func() {
				duration, _ := time.ParseDuration("3s")
				Expect(configuration.Duration).To(Equal(duration))
				Expect(configuration.Summary).To(Equal(true))
				Expect(configuration.Workers).To(Equal(50))
			})

			It("does not override the defaults for other args", func() {
				Expect(configuration.Random).To(Equal(false))
				Expect(configuration.WaitTime).To(Equal(defaultWaitTime))
			})
		})

		Describe("with invalid arg values", func() {
			Describe("for duration", func() {
				It("returns error", func() {
					args = []string{"--duration", "xs", filename}
					_, err := ParseConfiguration(args)
					Expect(err).Should(MatchError("Cannot parse the value specified for --duration: 'xs'"))
				})
			})

			Describe("for workers", func() {
				It("returns error", func() {
					args = []string{"--workers", "xs", filename}
					_, err := ParseConfiguration(args)
					Expect(err).Should(MatchError("strconv.ParseFloat: parsing \"xs\": invalid syntax"))
				})
			})

			Describe("for wait-time", func() {
				It("returns error", func() {
					args = []string{"--wait-time", "xs", filename}
					_, err := ParseConfiguration(args)
					Expect(err).Should(MatchError("Cannot parse the value specified for --wait-time: 'xs'"))
				})
			})
		})
	})
})
