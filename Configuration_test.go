package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

    "github.com/oleiade/reflections"
    "github.com/naoina/go-stringutil"
)

var _ = Describe("Configuration", func() {

	var configuration *Configuration
	var args []string
	defaultWaitTime := time.Duration(0)
	defaultDuration := time.Duration(0)
	filename, _ := filepath.Abs(os.Args[0])

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
		var (
			yaml          string
			configuration *Configuration
			err           error
		)
        duration5m, _ := time.ParseDuration("5m")

        //TODO Should the url file also be supported in the config files, thus removing it from being a required cmd arg, but required and valid once Configuration is built?
		testFixtures := []configurationTestFixture{
			{
				context: "duration",
				tests: []configurationTest{
                    {"passed on cmd but not set in pwd config or user home config", []string{"--duration", "5m", filename}, "", "", duration5m},
                    {"passed on cmd and set in pwd config and not set in user home config", []string{"--duration", "5m", filename}, "duration: 30s", "", duration5m},
                    {"passed on cmd and set in pwd config and set in user home config", []string{"--duration", "5m", filename}, "duration: 30s", "duration: 1m", duration5m},
                    {"set in pwd config and set in user home config", []string{filename}, "duration: 5m", "duration: 1m", duration5m},
                    {"set in pwd config and not set in user home config", []string{filename}, "duration: 5m", "", duration5m},
                    {"not set in pwd config but set in user home config", []string{filename}, "", "duration: 5m", duration5m},
                    {"not set in pwd config or user home config", []string{filename}, "", "", time.Duration(0)},
                },
            }, {
				context: "random",
				tests: []configurationTest{
                    {"passed on cmd but not set in pwd config or user home config", []string{"--random", filename}, "", "", true},
                    {"passed on cmd and set OFF in pwd config and not set in user home config", []string{"--random", filename}, "random: false", "", true},
                    {"passed on cmd and set OFF in pwd config and set OFF in user home config", []string{"--random", filename}, "random: false", "random: false", true},
                    {"set ON in pwd config and set OFF in user home config", []string{filename}, "random: true", "random: false", true},
                    {"set OFF in pwd config and set ON in user home config", []string{filename}, "random: false", "random: true", false},
                    {"set ON in pwd config and not set in user home config", []string{filename}, "random: true", "", true},
                    {"set OFF in pwd config and not set in user home config", []string{filename}, "random: false", "", false},
                    {"not set in pwd config but set ON in user home config", []string{filename}, "", "random: true", true},
                    {"not set in pwd config but set OFF in user home config", []string{filename}, "", "random: false", false},
                    {"not set in pwd config or user home config", []string{filename}, "", "", false},
                },
            }, {
				context: "summary",
				tests: []configurationTest{
                    {"passed on cmd but not set in pwd config or user home config", []string{"--summary", filename}, "", "", true},
                    {"passed on cmd and set OFF in pwd config and not set in user home config", []string{"--summary", filename}, "summary: false", "", true},
                    {"passed on cmd and set OFF in pwd config and set OFF in user home config", []string{"--summary", filename}, "summary: false", "summary: false", true},
                    {"set ON in pwd config and set OFF in user home config", []string{filename}, "summary: true", "summary: false", true},
                    {"set OFF in pwd config and set ON in user home config", []string{filename}, "summary: false", "summary: true", false},
                    {"set ON in pwd config and not set in user home config", []string{filename}, "summary: true", "", true},
                    {"set OFF in pwd config and not set in user home config", []string{filename}, "summary: false", "", false},
                    {"not set in pwd config but set ON in user home config", []string{filename}, "", "summary: true", true},
                    {"not set in pwd config but set OFF in user home config", []string{filename}, "", "summary: false", false},
                    {"not set in pwd config or user home config", []string{filename}, "", "", false},
                },
            }, {
                //TODO change this to "wait-time" when fixed in the stringutil library
				context: "wait_time",
				tests: []configurationTest{
                    {"passed on cmd but not set in pwd config or user home config", []string{"--wait-time", "5m", filename}, "", "", duration5m},
                    {"passed on cmd and set in pwd config and not set in user home config", []string{"--wait-time", "5m", filename}, "wait-time: 30s", "", duration5m},
                    {"passed on cmd and set in pwd config and set in user home config", []string{"--wait-time", "5m", filename}, "wait-time: 30s", "wait-time: 1m", duration5m},
                    {"set in pwd config and set in user home config", []string{filename}, "wait-time: 5m", "wait-time: 1m", duration5m},
                    {"set in pwd config and not set in user home config", []string{filename}, "wait-time: 5m", "", duration5m},
                    {"not set in pwd config but set in user home config", []string{filename}, "", "wait-time: 5m", duration5m},
                    {"not set in pwd config or user home config", []string{filename}, "", "", time.Duration(0)},
                },
            }, {
				context: "workers",
				tests: []configurationTest{
                    {"passed on cmd but not set in pwd config or user home config", []string{"--workers", "5", filename}, "", "", 5},
                    {"passed on cmd and set in pwd config and not set in user home config", []string{"--workers", "5", filename}, "workers: 3", "", 5},
                    {"passed on cmd and set in pwd config and set in user home config", []string{"--workers", "3", filename}, "workers: 3", "workers: 2", 5},
					{"set in pwd config and not set in user home config", []string{filename}, "workers: 3", "", 3},
                    {"set in pwd config and set in user home config", []string{filename}, "workers: 3", "workers: 5", 3},
					{"not set in pwd config but set in user home config", []string{filename}, "", "workers: 3", 3},
					{"not set in pwd config or user home config", []string{filename}, "", "", 1},
				},
			},
		}

		for _, fixture := range testFixtures {
            // This is that weird thing where if I just used fixture.context in the assertion it had got the one from the next test in the loop!
            context := fixture.context
			Context("for " + context, func() {
				for _, test := range fixture.tests {
					Context(test.name, func() {
						BeforeEach(func() {
							configFileReader = func(path string) ([]byte, error) {
								pwd, _ := os.Getwd()
								if strings.Contains(path, pwd) {
									yaml = test.pwdYaml
								} else {
									yaml = test.usrYaml
								}
								return []byte(yaml), nil
							}
							configuration, err = ParseConfiguration(test.cmdArgs)
							Expect(err).ShouldNot(HaveOccurred())
						})

						It("Parses the yaml and applies the config", func() {
                            fmt.Printf("Looking for %s in %+v", context, configuration)
                            actual, _ := reflections.GetField(configuration, stringutil.ToUpperCamelCase(context))
							Expect(actual).To(Equal(test.expected))
						})
					})
				}
			})
		}
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

type configurationTestFixture struct {
	context string
	tests   []configurationTest
}

type configurationTest struct {
	name     string
	cmdArgs  []string
	pwdYaml  string
	usrYaml  string
	expected interface{}
}
