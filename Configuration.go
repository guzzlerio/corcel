package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/imdario/mergo"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

//Configuration ...
type Configuration struct {
	Duration time.Duration `yaml:"duration"`
	FilePath string
	Random   bool          `yaml:"random"`
	Summary  bool          `yaml:"summary"`
	Workers  int           `yaml:"workers"`
	WaitTime time.Duration `yaml:"wait-time"`
}

func parseConfiguration(args []string) (*Configuration, error) {
	config := Configuration{}
	defaults := defaultConfig()
	cmd, err := cmdConfig(args)
	if err != nil {
		return nil, err
	}

	pwd, err := pwdConfig()
	if err != nil {
		return nil, err
	}
	usr, err := userDirConfig()
	if err != nil {
		return nil, err
	}

	/*
		log.Printf(" default: %+v\n", defaults)
		log.Printf("     cmd: %+v\n", cmd)
		log.Printf("     pwd: %+v\n", pwd)
		log.Printf("     usr: %+v\n", usr)
	*/

	if err := mergo.Merge(&config, &cmd); err != nil {
		return nil, err
	}
	if err := mergo.Merge(&config, &pwd); err != nil {
		return nil, err
	}
	if err := mergo.Merge(&config, &usr); err != nil {
		return nil, err
	}
	if err := mergo.Merge(&config, &defaults); err != nil {
		return nil, err
	}
	//log.Printf(" config: %+v\n", config)
	return &config, err
}

func cmdConfig(args []string) (Configuration, error) {
	config := Configuration{}
	CommandLine := kingpin.New("corcel", "")

	CommandLine.Version(applicationVersion)

	CommandLine.Arg("file", "Urls file").Required().StringVar(&config.FilePath)
	summary := CommandLine.Flag("summary", "Output summary to STDOUT").Bool()
	waitTimeArg := CommandLine.Flag("wait-time", "Time to wait between each execution").Default("0s").String()
	workers := CommandLine.Flag("workers", "The number of workers to execute the requests").Int()
	random := CommandLine.Flag("random", "Select the url at random for each execution").Bool()
	durationArg := CommandLine.Flag("duration", "The duration of the run e.g. 10s 10m 10h etc... valid values are  ms, s, m, h").String()

	_, err := CommandLine.Parse(args)

	if err != nil {
		//log.Println("Unable to parse the kingpin args")
		return Configuration{}, err
	}
	waitTime, err := time.ParseDuration(*waitTimeArg)
	if err != nil {
		return Configuration{}, fmt.Errorf("Cannot parse the value specified for --wait-time: '%v'", *waitTimeArg)
	}
	var duration time.Duration
	//remove this if when issue #17 is completed
	if *durationArg != "" {
		duration, err = time.ParseDuration(*durationArg)
		if err != nil {
			return Configuration{}, fmt.Errorf("Cannot parse the value specified for --duration: '%v'", *durationArg)
		}
	}

	if err = config.handleHTTPEndpointForURLFile(); err != nil {
		return Configuration{}, err
	}

	if _, err := os.Stat(config.FilePath); os.IsNotExist(err) {
		return Configuration{}, fmt.Errorf("required argument 'file' not provided")
	}

	absolutePath, err := filepath.Abs(config.FilePath)
	check(err)
	file, err := os.Open(absolutePath)
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing file %v", err)
		}
	}()
	check(err)

	return Configuration{
		Duration: duration,
		FilePath: config.FilePath,
		Random:   *random,
		Summary:  *summary,
		Workers:  *workers,
		WaitTime: waitTime,
	}, err
}

func pwdConfig() (Configuration, error) {
	pwd, _ := os.Getwd()
	// can we get the application name programatically?
	filename := path.Join(pwd, fmt.Sprintf(".%src", "corcel"))

	contents, err := configFileReader(filename)
	if err != nil {
		return Configuration{}, err
	}
	var config Configuration
	if err := config.parse(contents); err != nil {
		return Configuration{}, err
	}
	return config, nil
}

func userDirConfig() (Configuration, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return Configuration{}, err
	}
	home, err := homedir.Expand(dir)
	// can we get the application name programatically?
	filename := path.Join(home, fmt.Sprintf(".%src", "corcel"))

	contents, err := configFileReader(filename)
	if err != nil {
		return Configuration{}, err
	}
	var config Configuration
	if err := config.parse(contents); err != nil {
		return Configuration{}, err
	}
	return config, nil
}

func defaultConfig() Configuration {
	waitTime := time.Duration(0)
	duration := time.Duration(0)
	return Configuration{
		Duration: duration,
		Random:   false,
		Summary:  false,
		Workers:  1,
		WaitTime: waitTime,
	}
}

func (c *Configuration) handleHTTPEndpointForURLFile() error {
	u, e := url.ParseRequestURI(c.FilePath)
	if e == nil && u.Scheme != "" {
		log.Printf("Dowloading URL file from %s ...\n", c.FilePath)
		file, _ := createTemporaryFile(c.FilePath)
		out, _ := os.Create(file.Name())
		defer func() {
			out.Close()
		}()

		body, e := downloadURLFileFromEndpoint(c.FilePath)
		if e != nil {
			return fmt.Errorf("unable to download url file from endpoint %s [%s]", c.FilePath, e)
		}
		defer body.Close()
		_, _ = io.Copy(out, body)
		//log.Printf("url file downloaded: %s\n", file.Name())
		c.FilePath = file.Name()
	}
	return nil
}

func (c *Configuration) parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		//log.Println("Unable to parse config file")
		return nil
	}
	/*
	   if c.Hostname == "" {
	       return errors.New("Kitchen config: invalid `hostname`")
	   }
	*/
	return nil
}

var configFileReader = func(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		//log.Println("Config file not found")
		return nil, nil
	}
	//log.Println("file exists; processing...")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil
	}
	return data, nil
}

var downloadURLFileFromEndpoint = func(endpoint string) (io.ReadCloser, error) {
	resp, e := http.Get(endpoint)
	if e != nil {
		return nil, e
	}
	return resp.Body, nil
}

var createTemporaryFile = func(filePath string) (*os.File, error) {
	hashed := md5.Sum([]byte(filePath))
	return ioutil.TempFile(os.TempDir(), fmt.Sprintf("%x", hashed))
}
