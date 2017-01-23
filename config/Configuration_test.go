package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/naoina/go-stringutil"
	"github.com/oleiade/reflections"
)

var _ = Describe("Configuration", func() {
	var configuration *Configuration
	var args Configuration
	defaultWaitTime := time.Duration(0)
	defaultDuration := time.Duration(0)
	duration3s, _ := time.ParseDuration("3s")
	filename, _ := filepath.Abs(os.Args[0])

	BeforeEach(func() {
		//args = Configuration{FilePath: filename}
		cfg := &Configuration{
			FilePath: filename,
		}
		logrus.SetOutput(ioutil.Discard)
		//Log.Out = ioutil.Discard
		configFileReader = func(path string) ([]byte, error) {
			return []byte(""), nil
		}
		configuration, _ = ParseConfiguration(cfg)
	})

	Describe("When no config file is found and no command line args are provided", func() {
		Describe("Loading a default configuration", func() {
			It("sets duration (--duration)", func() {
				Expect(configuration.Duration).To(Equal(defaultDuration))
			})
			It("sets progress (--progress)", func() {
				Expect(configuration.Progress).To(Equal("logo"))
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
			It("sets summary format (--summary-format)", func() {
				Expect(configuration.SummaryFormat).To(Equal("console"))
			})
			It("sets log-level", func() {
				Expect(configuration.LogLevel).To(Equal(logrus.FatalLevel))
			})
		})
	})

	Describe("When config file is not found in pwd", func() {
		Describe("and config file is found in user home", func() {
		})
	})

	Describe("When config file is found in pwd", func() {
		var (
			yaml string
		)
		duration5m, _ := time.ParseDuration("5m")

		testFixtures := []configurationTestFixture{
			{
				context: "duration",
				tests: []configurationTest{
					{"passed on cmd but not set in pwd config or user home config", Configuration{Duration: duration5m, FilePath: filename}, "", "", duration5m},
					{"passed on cmd and set in pwd config and not set in user home config", Configuration{Duration: duration5m, FilePath: filename}, "duration: 30s", "", duration5m},
					{"passed on cmd and set in pwd config and set in user home config", Configuration{Duration: duration5m, FilePath: filename}, "duration: 30s", "duration: 1m", duration5m},
					{"set in pwd config and set in user home config", Configuration{FilePath: filename}, "duration: 5m", "duration: 1m", duration5m},
					{"set in pwd config and not set in user home config", Configuration{FilePath: filename}, "duration: 5m", "", duration5m},
					{"not set in pwd config but set in user home config", Configuration{FilePath: filename}, "", "duration: 5m", duration5m},
					{"not set in pwd config or user home config", Configuration{FilePath: filename}, "", "", time.Duration(0)},
					// unhappy paths
					{"set in pwd with invalid value and set in user home config", Configuration{FilePath: filename}, "duration: abc", "duration: 5m", duration5m},
					{"set in pwd and set in user home config with invalid value", Configuration{FilePath: filename}, "duration: 5m", "duration: abc", duration5m},
					{"set in pwd with invalid value and not set in user home config", Configuration{FilePath: filename}, "duration: abc", "", time.Duration(0)},
					{"not set in pwd but set in user home config with invalid value", Configuration{FilePath: filename}, "", "duration: abc", time.Duration(0)},
				},
			}, {
				context: "random",
				tests: []configurationTest{
					{"passed on cmd but not set in pwd config or user home config", Configuration{Random: true, FilePath: filename}, "", "", true},
					{"passed on cmd and set OFF in pwd config and not set in user home config", Configuration{Random: true, FilePath: filename}, "random: false", "", true},
					{"passed on cmd and set OFF in pwd config and set OFF in user home config", Configuration{Random: true, FilePath: filename}, "random: false", "random: false", true},
					{"set ON in pwd config and set OFF in user home config", Configuration{FilePath: filename}, "random: true", "random: false", true},
					{"set OFF in pwd config and set ON in user home config", Configuration{FilePath: filename}, "random: false", "random: true", false},
					{"set ON in pwd config and not set in user home config", Configuration{FilePath: filename}, "random: true", "", true},
					{"set OFF in pwd config and not set in user home config", Configuration{FilePath: filename}, "random: false", "", false},
					{"not set in pwd config but set ON in user home config", Configuration{FilePath: filename}, "", "random: true", true},
					{"not set in pwd config but set OFF in user home config", Configuration{FilePath: filename}, "", "random: false", false},
					{"not set in pwd config or user home config", Configuration{FilePath: filename}, "", "", false},
					// unhappy paths
					{"set in pwd with invalid value and set in user home config", Configuration{FilePath: filename}, "random: abc", "random: true", true},
					{"set in pwd and set in user home config with invalid value", Configuration{FilePath: filename}, "random: true", "random: abc", true},
					{"set in pwd with invalid value and not set in user home config", Configuration{FilePath: filename}, "random: abc", "", false},
					{"not set in pwd but set in user home config with invalid value", Configuration{FilePath: filename}, "", "random: abc", false},
				},
			}, {
				context: "summary",
				tests: []configurationTest{
					{"passed on cmd but not set in pwd config or user home config", Configuration{Summary: true, FilePath: filename}, "", "", true},
					{"passed on cmd and set OFF in pwd config and not set in user home config", Configuration{Summary: true, FilePath: filename}, "summary: false", "", true},
					{"passed on cmd and set OFF in pwd config and set OFF in user home config", Configuration{Summary: true, FilePath: filename}, "summary: false", "summary: false", true},
					{"set ON in pwd config and set OFF in user home config", Configuration{FilePath: filename}, "summary: true", "summary: false", true},
					{"set OFF in pwd config and set ON in user home config", Configuration{FilePath: filename}, "summary: false", "summary: true", false},
					{"set ON in pwd config and not set in user home config", Configuration{FilePath: filename}, "summary: true", "", true},
					{"set OFF in pwd config and not set in user home config", Configuration{FilePath: filename}, "summary: false", "", false},
					{"not set in pwd config but set ON in user home config", Configuration{FilePath: filename}, "", "summary: true", true},
					{"not set in pwd config but set OFF in user home config", Configuration{FilePath: filename}, "", "summary: false", false},
					{"not set in pwd config or user home config", Configuration{FilePath: filename}, "", "", false},
					// unhappy paths
					{"set in pwd with invalid value and set in user home config", Configuration{FilePath: filename}, "summary: abc", "summary: true", true},
					{"set in pwd and set in user home config with invalid value", Configuration{FilePath: filename}, "summary: true", "summary: abc", true},
					{"set in pwd with invalid value and not set in user home config", Configuration{FilePath: filename}, "summary: abc", "", false},
					{"not set in pwd but set in user home config with invalid value", Configuration{FilePath: filename}, "", "summary: abc", false},
				},
			}, {
				//TODO change this to "wait-time" when fixed in the stringutil library
				context: "wait_time",
				tests: []configurationTest{
					{"passed on cmd but not set in pwd config or user home config", Configuration{WaitTime: duration5m, FilePath: filename}, "", "", duration5m},
					{"passed on cmd and set in pwd config and not set in user home config", Configuration{WaitTime: duration5m, FilePath: filename}, "wait-time: 30s", "", duration5m},
					{"passed on cmd and set in pwd config and set in user home config", Configuration{WaitTime: duration5m, FilePath: filename}, "wait-time: 30s", "wait-time: 1m", duration5m},
					{"set in pwd config and set in user home config", Configuration{FilePath: filename}, "wait-time: 5m", "wait-time: 1m", duration5m},
					{"set in pwd config and not set in user home config", Configuration{FilePath: filename}, "wait-time: 5m", "", duration5m},
					{"not set in pwd config but set in user home config", Configuration{FilePath: filename}, "", "wait-time: 5m", duration5m},
					{"not set in pwd config or user home config", Configuration{FilePath: filename}, "", "", time.Duration(0)},
					// unhappy paths
					{"set in pwd with invalid value and set in user home config", Configuration{FilePath: filename}, "wait-time: abc", "wait-time: 5m", duration5m},
					{"set in pwd and set in user home config with invalid value", Configuration{FilePath: filename}, "wait-time: 5m", "wait-time: abc", duration5m},
					{"set in pwd with invalid value and not set in user home config", Configuration{FilePath: filename}, "wait-time: abc", "", time.Duration(0)},
					{"not set in pwd but set in user home config with invalid value", Configuration{FilePath: filename}, "", "wait-time: abc", time.Duration(0)},
				},
			}, {
				context: "workers",
				tests: []configurationTest{
					{"passed on cmd but not set in pwd config or user home config", Configuration{Workers: 5, FilePath: filename}, "", "", 5},
					{"passed on cmd and set in pwd config and not set in user home config", Configuration{Workers: 5, FilePath: filename}, "workers: 3", "", 5},
					{"passed on cmd and set in pwd config and set in user home config", Configuration{Workers: 3, FilePath: filename}, "workers: 3", "workers: 2", 5},
					{"set in pwd config and not set in user home config", Configuration{FilePath: filename}, "workers: 3", "", 3},
					{"set in pwd config and set in user home config", Configuration{FilePath: filename}, "workers: 3", "workers: 5", 3},
					{"not set in pwd config but set in user home config", Configuration{FilePath: filename}, "", "workers: 3", 3},
					{"not set in pwd config or user home config", Configuration{FilePath: filename}, "", "", 1},
					// unhappy paths
					{"set in pwd with invalid value and set in user home config", Configuration{FilePath: filename}, "workers: abc", "workers: 5", 5},
					{"set in pwd and set in user home config with invalid value", Configuration{FilePath: filename}, "workers: 5", "workers: abc", 5},
					{"set in pwd with invalid value and not set in user home config", Configuration{FilePath: filename}, "workers: abc", "", 1},
					{"not set in pwd but set in user home config with invalid value", Configuration{FilePath: filename}, "", "workers: abc", 1},
				},
			}, {
				context: "log_level",
				tests: []configurationTest{
					{"set in pwd config and not set in user home config", Configuration{FilePath: filename}, "log-level: 3", "", logrus.WarnLevel},
				},
			},
		}

		for _, fixture := range testFixtures {
			// This is that weird thing where if I just used fixture.context in the assertion it had got the one from the next test in the loop!
			context := fixture.context
			Context("for "+context, func() {
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
							var err error
							configuration, err = ParseConfiguration(&test.cmdArgs)
							Expect(err).ShouldNot(HaveOccurred())
						})

						It("Parses the yaml and applies the config", func() {
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
					args = Configuration{Duration: duration3s, FilePath: filename}
					configuration, _ = ParseConfiguration(&args)
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
					args = Configuration{FilePath: filename}
					configuration, _ = ParseConfiguration(&args)
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
					args = Configuration{Random: true, FilePath: filename}
					configuration, _ = ParseConfiguration(&args)
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
					args = Configuration{Summary: true, FilePath: filename}
					configuration, _ = ParseConfiguration(&args)
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
					args = Configuration{Workers: 3, FilePath: filename}
					configuration, _ = ParseConfiguration(&args)
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
					args = Configuration{WaitTime: duration3s, FilePath: filename}
					configuration, _ = ParseConfiguration(&args)
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

			//These really need to be moved to test more outside-in
			Describe("for verbosity", func() {
				Context("not set", func() {
					BeforeEach(func() {
						args = Configuration{FilePath: filename}
						configuration, _ = ParseConfiguration(&args)
					})

					It("sets verbosity to Fatal", func() {
						Expect(configuration.LogLevel).To(Equal(logrus.FatalLevel))
					})
				})

				PContext("-v", func() {
					BeforeEach(func() {
						//args = []string{"-v", filename}
						configuration, _ = ParseConfiguration(&Configuration{}) //args)
					})

					It("sets verbosity to Fatal", func() {
						Expect(configuration.LogLevel).To(Equal(logrus.WarnLevel))
					})
				})

				PContext("-vv", func() {
					var err error
					BeforeEach(func() {
						//args = []string{"-v", filename}
						configuration, err = ParseConfiguration(&Configuration{}) //args)
						Expect(err).To(BeNil())
						Expect(configuration).ToNot(BeNil())
					})

					It("sets verbosity to Fatal", func() {
						Expect(configuration.LogLevel).To(Equal(logrus.InfoLevel))
					})
				})
				PContext("-vvv", func() {
					var err error
					BeforeEach(func() {
						//args = []string{"-vvv", filename}
						configuration, err = ParseConfiguration(&Configuration{}) //args)
						Expect(err).To(BeNil())
					})

					It("sets verbosity to Fatal", func() {
						Expect(configuration.LogLevel).To(Equal(logrus.DebugLevel))
					})
				})
			})
		})
		/*
					   TODO:
					   There does not seem to be support for the current implementation
					   of multiple flags i.e. -vvvv when verbose is set to Bool() throws
					   an error as it tries to parse 'v' as a bool.  Pending them for now
			       until we can find a way to achieve this.

		*/
		Describe("setting multiple command line args", func() {
			var err error
			BeforeEach(func() {
				args = Configuration{Summary: true, Workers: 50, Duration: duration3s, FilePath: filename}
				configuration, err = ParseConfiguration(&args)
				Expect(err).To(BeNil())
			})
			It("applies the overrides", func() {
				Expect(configuration.Duration).To(Equal(duration3s))
				Expect(configuration.Summary).To(Equal(true))
				Expect(configuration.Workers).To(Equal(50))
				//Removed the -vv from the setup due to Issue #47
				//Expect(configuration.LogLevel).To(Equal(logrus.InfoLevel))
			})

			It("does not override the defaults for other args", func() {
				Expect(configuration.Random).To(Equal(false))
				Expect(configuration.WaitTime).To(Equal(defaultWaitTime))
			})
		})

		Describe("with invalid arg values", func() {
			Describe("missing url file", func() {
				It("returns error", func() {
					args = Configuration{}
					_, err := ParseConfiguration(&args)
					Expect(err).Should(MatchError("required argument 'file' not provided"))
				})
			})

			//TODO These are probably crap as they would be testing kingpin now.
			/*
				Describe("for duration", func() {
					It("returns error", func() {
						args = string{"--duration", "xs", filename}
						_, err := ParseConfiguration(args)
						Expect(err).Should(MatchError("time: invalid duration xs"))
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
						Expect(err).Should(MatchError("time: invalid duration xs"))
					})
				})
			*/
		})

		Describe("providing a HTTP endpoint for the url file", func() {
			var tmpFile, endpoint string
			BeforeEach(func() {
				endpoint = "http://some-url/to/download/from"
				createTemporaryFile = func(filePath string) (*os.File, error) {
					hashed := md5.Sum([]byte(filePath))
					file, fileErr := ioutil.TempFile(os.TempDir(), fmt.Sprintf("%x", hashed))
					tmpFile = file.Name()
					return file, fileErr
				}
			})

			Context("when the file is downloaded successfully", func() {
				BeforeEach(func() {
					downloadURLFileFromEndpoint = func(endpoint string) (io.ReadCloser, error) {
						return ioutil.NopCloser(strings.NewReader("http://something")), nil
					}

					args = Configuration{FilePath: endpoint}
					var err error
					configuration, err = ParseConfiguration(&args)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("applies the override", func() {
					Expect(configuration.FilePath).To(Equal(tmpFile))
				})
			})

			Context("when the download fails", func() {
				BeforeEach(func() {
					downloadURLFileFromEndpoint = func(endpoint string) (io.ReadCloser, error) {
						return nil, fmt.Errorf("booom")
					}
				})

				It("returns error", func() {
					args = Configuration{FilePath: endpoint}
					_, err := ParseConfiguration(&args)
					Expect(err).Should(MatchError("unable to download url file from endpoint " + endpoint + " [booom]"))
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
	cmdArgs  Configuration
	pwdYaml  string
	usrYaml  string
	expected interface{}
}
