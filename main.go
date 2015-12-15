package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/logger"
)

var (
	//ErrorMappings ...
	ErrorMappings = map[string]ErrorCode{}
)

func check(err error) {
	if err != nil {
		for mapping, errorCode := range ErrorMappings {
			if strings.Contains(fmt.Sprintf("%v", err), mapping) {
				fmt.Println(errorCode.Message)
				os.Exit(errorCode.Code)
			}
		}
		logger.Log.Fatalf("UNKNOWN ERROR: %v", err)
	}
}

type ExecutionId struct {
	value string
}

func (id ExecutionId) String() string {
	return fmt.Sprintf("%s", id.value)
}

func NewExecutionId() ExecutionId {
	id := RandString(16)
	return ExecutionId{id}
}

func RandString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

type Control interface {
	Start(*config.Configuration) (ExecutionId, error)
	Stop() ExecutionOutput
	Status(*ExecutionId) ExecutionOutput
	History() []*ExecutionId
	Events() <-chan string

	//TODO Shouldn't need to expose this out, but required for transition
	Statistics() Statistics
}

type Controller struct {
	stats      *Statistics
	executions map[ExecutionId]*Executor
	bar ProgressBar
}

func (instance *Controller) Start(config *config.Configuration) (ExecutionId, error) {
	instance.stats.Start()
	id := NewExecutionId()
	fmt.Printf("Execution ID: %s\n", id)
	executor := Executor{config, instance.stats, instance.bar}
	instance.executions[id] = &executor
	executor.Execute()
	return id, nil
}
func (instance *Controller) Stop() ExecutionOutput {
	instance.stats.Stop()
	return instance.stats.ExecutionOutput()
}
func (instance *Controller) Status(*ExecutionId) ExecutionOutput {
	return ExecutionOutput{}
}
func (instance *Controller) History() []*ExecutionId {
	return nil
}
func (instance *Controller) Events() <-chan string {
	return nil
}
func (instance *Controller) Statistics() Statistics {
	return *instance.stats
}

func NewControl(bar ProgressBar) Control {
	stats := CreateStatistics()
	control := Controller{stats: stats, executions: make(map[ExecutionId]*Executor), bar: bar}
	return &control
}

type Host interface {
	SetControl(*Control)
}

type ConsoleHost struct {
	Control Control
}

func (host *ConsoleHost) SetControl(control Control) {
	host.Control = control
}

func NewConsoleHost(config *config.Configuration) ConsoleHost {
	host := ConsoleHost{}
	bar := NewProgressBar(100, config)
	control := NewControl(bar)
	host.SetControl(control)
	return host
}

//GenerateExecutionOutput ...
func GenerateExecutionOutput(file string, output ExecutionOutput) {
	outputPath, err := filepath.Abs(file)
	check(err)
	yamlOutput, err := yaml.Marshal(&output)
	check(err)
	err = ioutil.WriteFile(outputPath, yamlOutput, 0644)
	check(err)
}

func main() {
	config, err := config.ParseConfiguration(os.Args[1:])
	if err != nil {
		logger.Log.Fatal(err)
	}

	configureErrorMappings()
	logger.ConfigureLogging(config)

	absolutePath, err := filepath.Abs(config.FilePath)
	check(err)
	file, err := os.Open(absolutePath)
	defer func() {
		err := file.Close()
		if err != nil {
			logger.Log.Printf("Error closing file %v", err)
		}
	}()
	check(err)

	host := NewConsoleHost(config)
	host.Control.Start(config) //will this block?
	output := host.Control.Stop()

	//TODO these should probably be pushed behind the host.Control.Stop afterall the host is a cmd host
	GenerateExecutionOutput("./output.yml", output)

	if config.Summary {
		consoleWriter := ExecutionOutputWriter{output}
		consoleWriter.Write(os.Stdout)
	}
}
