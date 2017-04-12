package config

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

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

func TestSerialisation(t *testing.T) {

	Convey("marshalling", t, func() {

		Convey("Marhsal Duration", func() {
			input := "duration: 5m"
			c := &Configuration{}
			c.parse([]byte(input))

			So(c.Duration, ShouldEqual, time.Duration(5*time.Minute))
		})

		Convey("Marhsal Waittime", func() {
			input := `wait-time: 5m
Workers: 50`
			c := &Configuration{}
			c.parse([]byte(input))

			So(c.WaitTime, ShouldEqual, time.Duration(5*time.Minute))
			So(c.Workers, ShouldEqual, 50)
		})
	})
}

func TestConfiguration(t *testing.T) {
	Convey("Configuration", t, func() {
		var configuration *Configuration
		var args Configuration
		defaultWaitTime := time.Duration(0)
		defaultDuration := time.Duration(0)
		duration3s, _ := time.ParseDuration("3s")
		filename, _ := filepath.Abs(os.Args[0])

		func() {
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
		}()

		Convey("When no config file is found and no command line args are provided", func() {
			Convey("Loading a default configuration", func() {
				Convey("sets duration (--duration)", func() {
					So(configuration.Duration, ShouldEqual, defaultDuration)
				})
				Convey("sets progress (--progress)", func() {
					So(configuration.Progress, ShouldEqual, "logo")
				})
				Convey("sets random (--random)", func() {
					So(configuration.Random, ShouldEqual, false)
				})
				Convey("sets summary (--summary)", func() {
					So(configuration.Summary, ShouldEqual, false)
				})
				Convey("sets workers (--workers)", func() {
					So(configuration.Workers, ShouldEqual, 1)
				})
				Convey("sets wait-time (--wait-time)", func() {
					So(configuration.WaitTime, ShouldEqual, defaultWaitTime)
				})
				Convey("sets summary format (--summary-format)", func() {
					So(configuration.SummaryFormat, ShouldEqual, "console")
				})
				Convey("sets log-level", func() {
					So(configuration.LogLevel, ShouldEqual, logrus.FatalLevel)
				})
			})
		})

		Convey("When config file is not found in pwd", func() {
			Convey("and config file is found in user home", func() {
			})
		})

		Convey("When config file is found in pwd", func() {
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
						/*
							{"set in pwd with invalid value and set in user home config", Configuration{FilePath: filename}, "duration: abc", "duration: 5m", duration5m},
							{"set in pwd and set in user home config with invalid value", Configuration{FilePath: filename}, "duration: 5m", "duration: abc", duration5m},
							{"set in pwd with invalid value and not set in user home config", Configuration{FilePath: filename}, "duration: abc", "", time.Duration(0)},
							{"not set in pwd but set in user home config with invalid value", Configuration{FilePath: filename}, "", "duration: abc", time.Duration(0)},
						*/
					},
				}, {
					context: "random",
					tests: []configurationTest{
						{"passed on cmd but not set in pwd config or user home config", Configuration{Random: true, FilePath: filename}, "", "", true},
						{"passed on cmd and set OFF in pwd config and not set in user home config", Configuration{Random: true, FilePath: filename}, "random: false", "", true},
						{"passed on cmd and set OFF in pwd config and set OFF in user home config", Configuration{Random: true, FilePath: filename}, "random: false", "random: false", true},
						{"set ON in pwd config and set OFF in user home config", Configuration{FilePath: filename}, "random: true", "random: false", true},
						//{"set OFF in pwd config and set ON in user home config", Configuration{FilePath: filename}, "random: false", "random: true", false},
						{"set ON in pwd config and not set in user home config", Configuration{FilePath: filename}, "random: true", "", true},
						{"set OFF in pwd config and not set in user home config", Configuration{FilePath: filename}, "random: false", "", false},
						{"not set in pwd config but set ON in user home config", Configuration{FilePath: filename}, "", "random: true", true},
						{"not set in pwd config but set OFF in user home config", Configuration{FilePath: filename}, "", "random: false", false},
						{"not set in pwd config or user home config", Configuration{FilePath: filename}, "", "", false},
						// unhappy paths
						/*
							{"set in pwd with invalid value and set in user home config", Configuration{FilePath: filename}, "random: abc", "random: true", true},
							{"set in pwd and set in user home config with invalid value", Configuration{FilePath: filename}, "random: true", "random: abc", true},
							{"set in pwd with invalid value and not set in user home config", Configuration{FilePath: filename}, "random: abc", "", false},
							{"not set in pwd but set in user home config with invalid value", Configuration{FilePath: filename}, "", "random: abc", false},
						*/
					},
				}, {
					context: "summary",
					tests: []configurationTest{
						{"passed on cmd but not set in pwd config or user home config", Configuration{Summary: true, FilePath: filename}, "", "", true},
						{"passed on cmd and set OFF in pwd config and not set in user home config", Configuration{Summary: true, FilePath: filename}, "summary: false", "", true},
						{"passed on cmd and set OFF in pwd config and set OFF in user home config", Configuration{Summary: true, FilePath: filename}, "summary: false", "summary: false", true},
						{"set ON in pwd config and set OFF in user home config", Configuration{FilePath: filename}, "summary: true", "summary: false", true},
						{"set OFF in pwd config and not set in user home config", Configuration{FilePath: filename}, "summary: false", "", false},
						// {"set OFF in pwd config and set ON in user home config", Configuration{FilePath: filename}, "summary: false", "summary: true", false},
						{"set ON in pwd config and not set in user home config", Configuration{FilePath: filename}, "summary: true", "", true},
						{"not set in pwd config but set ON in user home config", Configuration{FilePath: filename}, "", "summary: true", true},
						{"not set in pwd config but set OFF in user home config", Configuration{FilePath: filename}, "", "summary: false", false},
						{"not set in pwd config or user home config", Configuration{FilePath: filename}, "", "", false},
						// unhappy paths
						/*
							{"set in pwd with invalid value and set in user home config", Configuration{FilePath: filename}, "summary: abc", "summary: true", true},
							{"set in pwd and set in user home config with invalid value", Configuration{FilePath: filename}, "summary: true", "summary: abc", true},
							{"set in pwd with invalid value and not set in user home config", Configuration{FilePath: filename}, "summary: abc", "", false},
							{"not set in pwd but set in user home config with invalid value", Configuration{FilePath: filename}, "", "summary: abc", false},
						*/
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
						/*
							{"set in pwd with invalid value and set in user home config", Configuration{FilePath: filename}, "wait-time: abc", "wait-time: 5m", duration5m},
							{"set in pwd and set in user home config with invalid value", Configuration{FilePath: filename}, "wait-time: 5m", "wait-time: abc", duration5m},
							{"set in pwd with invalid value and not set in user home config", Configuration{FilePath: filename}, "wait-time: abc", "", time.Duration(0)},
							{"not set in pwd but set in user home config with invalid value", Configuration{FilePath: filename}, "", "wait-time: abc", time.Duration(0)},
						*/
					},
				}, {
					context: "workers",
					tests: []configurationTest{
						{"passed on cmd but not set in pwd config or user home config", Configuration{Workers: 5, FilePath: filename}, "", "", 5},
						{"passed on cmd and set in pwd config and not set in user home config", Configuration{Workers: 5, FilePath: filename}, "workers: 3", "", 5},
						{"passed on cmd and set in pwd config and set in user home config", Configuration{Workers: 5, FilePath: filename}, "workers: 3", "workers: 2", 5},
						{"set in pwd config and not set in user home config", Configuration{FilePath: filename}, "workers: 3", "", 3},
						{"set in pwd config and set in user home config", Configuration{FilePath: filename}, "workers: 3", "workers: 5", 3},
						{"not set in pwd config but set in user home config", Configuration{FilePath: filename}, "", "workers: 3", 3},
						{"not set in pwd config or user home config", Configuration{FilePath: filename}, "", "", 1},
						// unhappy paths
						//{"set in pwd with invalid value and set in user home config", Configuration{FilePath: filename}, "workers: abc", "workers: 5", 5},
						//{"set in pwd and set in user home config with invalid value", Configuration{FilePath: filename}, "workers: 5", "workers: abc", 5},
						//{"set in pwd with invalid value and not set in user home config", Configuration{FilePath: filename}, "workers: abc", "", 1},
						//{"not set in pwd but set in user home config with invalid value", Configuration{FilePath: filename}, "", "workers: abc", 1},
					},
				}, {
					context: "log_level",
					tests: []configurationTest{
						{"set in pwd config and not set in user home config", Configuration{FilePath: filename}, "log-level: 3", "", logrus.WarnLevel},
					},
				},
			}

			Convey("Test Configuration Context Test Cases", func() {
				for _, fixture := range testFixtures {
					// This is that weird thing where if I just used fixture.context in the assertion it had got the one from the next test in the loop!
					context := fixture.context
					for _, test := range fixture.tests {
						//Convey(test.name, func() {
						//func() {
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
						So(err, ShouldBeNil)
						//}()

						//Convey("Parses the yaml and applies the config", func() {
						actual, _ := reflections.GetField(configuration, stringutil.ToUpperCamelCase(context))
						So(actual, ShouldEqual, test.expected)
						//})
						//})
					}
				}
			})
		})

		Convey("When commandline args are provided", func() {
			Convey("overriding the default configuration", func() {
				Convey("for duration (--duration)", func() {
					func() {
						args = Configuration{Duration: duration3s, FilePath: filename}
						configuration, _ = ParseConfiguration(&args)
					}()
					Convey("applies the override", func() {
						duration, _ := time.ParseDuration("3s")
						So(configuration.Duration, ShouldEqual, duration)
					})

					Convey("leaves the default for", func() {
						Convey("random", func() {
							So(configuration.Random, ShouldEqual, false)
						})
						Convey("summary", func() {
							So(configuration.Summary, ShouldEqual, false)
						})
						Convey("workers", func() {
							So(configuration.Workers, ShouldEqual, 1)
						})
						Convey("wait-time", func() {
							So(configuration.WaitTime, ShouldEqual, defaultWaitTime)
						})
					})
				})

				Convey("for file", func() {
					func() {
						args = Configuration{FilePath: filename}
						configuration, _ = ParseConfiguration(&args)
					}()
					Convey("applies the override", func() {
						So(configuration.FilePath, ShouldEqual, filename)
					})

					Convey("leaves the default for", func() {
						Convey("duration", func() {
							So(configuration.Duration, ShouldEqual, defaultDuration)
						})
						Convey("random", func() {
							So(configuration.Random, ShouldEqual, false)
						})
						Convey("summary", func() {
							So(configuration.Summary, ShouldEqual, false)
						})
						Convey("workers", func() {
							So(configuration.Workers, ShouldEqual, 1)
						})
						Convey("wait-time", func() {
							So(configuration.WaitTime, ShouldEqual, defaultWaitTime)
						})
					})
				})

				Convey("for random (--random)", func() {
					func() {
						args = Configuration{Random: true, FilePath: filename}
						configuration, _ = ParseConfiguration(&args)
					}()
					Convey("applies the override", func() {
						So(configuration.Random, ShouldEqual, true)
					})

					Convey("leaves the default for", func() {
						Convey("duration", func() {
							So(configuration.Duration, ShouldEqual, defaultDuration)
						})
						Convey("summary", func() {
							So(configuration.Summary, ShouldEqual, false)
						})
						Convey("workers", func() {
							So(configuration.Workers, ShouldEqual, 1)
						})
						Convey("wait-time", func() {
							So(configuration.WaitTime, ShouldEqual, defaultWaitTime)
						})
					})
				})

				Convey("for summary (--summary)", func() {
					func() {
						args = Configuration{Summary: true, FilePath: filename}
						configuration, _ = ParseConfiguration(&args)
					}()
					Convey("applies the override", func() {
						So(configuration.Summary, ShouldEqual, true)
					})

					Convey("leaves the default for", func() {
						Convey("duration", func() {
							So(configuration.Duration, ShouldEqual, defaultDuration)
						})
						Convey("random", func() {
							So(configuration.Random, ShouldEqual, false)
						})
						Convey("workers", func() {
							So(configuration.Workers, ShouldEqual, 1)
						})
						Convey("wait-time", func() {
							So(configuration.WaitTime, ShouldEqual, defaultWaitTime)
						})
					})
				})

				Convey("for workers (--workers)", func() {
					func() {
						args = Configuration{Workers: 3, FilePath: filename}
						configuration, _ = ParseConfiguration(&args)
					}()
					Convey("applies the override", func() {
						So(configuration.Workers, ShouldEqual, 3)
					})

					Convey("leaves the default for", func() {
						Convey("duration", func() {
							So(configuration.Duration, ShouldEqual, defaultDuration)
						})
						Convey("random", func() {
							So(configuration.Random, ShouldEqual, false)
						})
						Convey("summary", func() {
							So(configuration.Summary, ShouldEqual, false)
						})
						Convey("wait-time", func() {
							So(configuration.WaitTime, ShouldEqual, defaultWaitTime)
						})
					})
				})

				Convey("for wait-time (--wait-time)", func() {
					func() {
						args = Configuration{WaitTime: duration3s, FilePath: filename}
						configuration, _ = ParseConfiguration(&args)
					}()
					Convey("applies the override", func() {
						waitTime, _ := time.ParseDuration("3s")
						So(configuration.WaitTime, ShouldEqual, waitTime)
					})

					Convey("leaves the default for", func() {
						Convey("duration", func() {
							So(configuration.Duration, ShouldEqual, defaultDuration)
						})
						Convey("random", func() {
							So(configuration.Random, ShouldEqual, false)
						})
						Convey("summary", func() {
							So(configuration.Summary, ShouldEqual, false)
						})
						Convey("workers", func() {
							So(configuration.Workers, ShouldEqual, 1)
						})
					})
				})

				//These really need to be moved to test more outside-in
				Convey("for verbosity", func() {
					Convey("not set", func() {
						func() {
							args = Configuration{FilePath: filename}
							configuration, _ = ParseConfiguration(&args)
						}()

						Convey("sets verbosity to Fatal", func() {
							So(configuration.LogLevel, ShouldEqual, logrus.FatalLevel)
						})
					})

					SkipConvey("-v", func() {
						func() {
							//args = []string{"-v", filename}
							configuration, _ = ParseConfiguration(&Configuration{}) //args)
						}()

						Convey("sets verbosity to Fatal", func() {
							So(configuration.LogLevel, ShouldEqual, logrus.WarnLevel)
						})
					})

					SkipConvey("-vv", func() {
						var err error
						func() {
							//args = []string{"-v", filename}
							configuration, err = ParseConfiguration(&Configuration{}) //args)
							So(err, ShouldBeNil)
							So(configuration, ShouldNotBeNil)
						}()

						Convey("sets verbosity to Fatal", func() {
							So(configuration.LogLevel, ShouldEqual, logrus.InfoLevel)
						})
					})
					SkipConvey("-vvv", func() {
						var err error
						func() {
							//args = []string{"-vvv", filename}
							configuration, err = ParseConfiguration(&Configuration{}) //args)
							So(err, ShouldBeNil)
						}()

						Convey("sets verbosity to Fatal", func() {
							So(configuration.LogLevel, ShouldEqual, logrus.DebugLevel)
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
			Convey("setting multiple command line args", func() {
				var err error
				func() {
					args = Configuration{Summary: true, Workers: 50, Duration: duration3s, FilePath: filename}
					configuration, err = ParseConfiguration(&args)
					So(err, ShouldBeNil)
				}()
				Convey("applies the overrides", func() {
					So(configuration.Duration, ShouldEqual, duration3s)
					So(configuration.Summary, ShouldEqual, true)
					So(configuration.Workers, ShouldEqual, 50)
					//Removed the -vv from the setup due to Issue #47
					//So(configuration.LogLevel, ShouldEqual, logrus.InfoLevel)
				})

				Convey("does not override the defaults for other args", func() {
					So(configuration.Random, ShouldEqual, false)
					So(configuration.WaitTime, ShouldEqual, defaultWaitTime)
				})
			})

			Convey("with invalid arg values", func() {
				Convey("missing url file", func() {
					Convey("returns error", func() {
						args = Configuration{}
						_, err := ParseConfiguration(&args)
						So(err, ShouldResemble, errors.New("required argument 'file' not provided"))
					})
				})

				//TODO These are probably crap as they would be testing kingpin now.
				/*
					Convey("for duration", func() {
						Convey("returns error", func() {
							args = string{"--duration", "xs", filename}
							_, err := ParseConfiguration(args)
							So(err).Should(MatchError("time: invalid duration xs"))
						})
					})

					Convey("for workers", func() {
						Convey("returns error", func() {
							args = []string{"--workers", "xs", filename}
							_, err := ParseConfiguration(args)
							So(err).Should(MatchError("strconv.ParseFloat: parsing \"xs\": invalid syntax"))
						})
					})

					Convey("for wait-time", func() {
						Convey("returns error", func() {
							args = []string{"--wait-time", "xs", filename}
							_, err := ParseConfiguration(args)
							So(err).Should(MatchError("time: invalid duration xs"))
						})
					})
				*/
			})

			Convey("providing a HTTP endpoint for the url file", func() {
				var tmpFile, endpoint string
				func() {
					endpoint = "http://some-url/to/download/from"
					createTemporaryFile = func(filePath string) (*os.File, error) {
						hashed := md5.Sum([]byte(filePath))
						file, fileErr := ioutil.TempFile(os.TempDir(), fmt.Sprintf("%x", hashed))
						tmpFile = file.Name()
						return file, fileErr
					}
				}()

				Convey("when the file is downloaded successfully", func() {
					func() {
						downloadURLFileFromEndpoint = func(endpoint string) (io.ReadCloser, error) {
							return ioutil.NopCloser(strings.NewReader("http://something")), nil
						}

						args = Configuration{FilePath: endpoint}
						var err error
						configuration, err = ParseConfiguration(&args)
						So(err, ShouldBeNil)
					}()

					Convey("applies the override", func() {
						So(configuration.FilePath, ShouldEqual, tmpFile)
					})
				})

				Convey("when the download fails", func() {
					func() {
						downloadURLFileFromEndpoint = func(endpoint string) (io.ReadCloser, error) {
							return nil, fmt.Errorf("booom")
						}
					}()

					Convey("returns error", func() {
						args = Configuration{FilePath: endpoint}
						_, err := ParseConfiguration(&args)
						So(err, ShouldResemble, errors.New("unable to download url file from endpoint "+endpoint+" [booom]"))
					})
				})
			})
		})
	})
}

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
