package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"time"

	"github.com/rcrowley/go-metrics"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var (
	logEnabled = false
	Log        *log.Logger
)

func check(err error) {
	if err != nil {
		Log.Panic(err)
	}
}

func configureLogging() {
	//TODO: refine this to work with levels or replace
	//with a package which already handles this
	flags := log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	prefix := "cns: "
	if logEnabled {
		Log = log.New(os.Stdout, prefix, flags)
	} else {
		//Send all the output to dev null
		Log = log.New(ioutil.Discard, prefix, flags)
	}
}

func main() {
	filePath := kingpin.Flag("file", "Urls file").Short('f').String()
	kingpin.Parse()

	configureLogging()

	absolutePath, err := filepath.Abs(*filePath)
	check(err)
	file, err := os.Open(absolutePath)
	check(err)

	defer file.Close()
	scanner := bufio.NewScanner(file)

	client := &http.Client{}
	requestAdapter := NewRequestAdapter()

	sampleSize := 1024
	hBytesSent := metrics.NewHistogram(metrics.NewUniformSample(sampleSize))
	hBytesReceived := metrics.NewHistogram(metrics.NewUniformSample(sampleSize))
	hResponseTime :=  metrics.NewHistogram(metrics.NewUniformSample(sampleSize))
	mBytesSent := metrics.NewMeter()
	mBytesReceived := metrics.NewMeter()

	for scanner.Scan() {
		line := scanner.Text()
		request, err := requestAdapter.Create(line)
		check(err)
		start := time.Now()
		response, err := client.Do(request)
		duration := time.Since(start) / time.Millisecond
		check(err)
		requestBytes, _ := httputil.DumpRequest(request, true)
		responseBytes, _ := httputil.DumpResponse(response, true)

		hBytesSent.Update(int64(len(requestBytes)))
		hBytesReceived.Update(int64(len(responseBytes)))

		mBytesSent.Mark(int64(len(requestBytes)))
		mBytesReceived.Mark(int64(len(responseBytes)))

		hResponseTime.Update(int64(duration))
	}

	summaryPath, err := filepath.Abs("./output.yml")
	check(err)

	output := ExecutionOutput{
		Summary: ExecutionSummary{
			ResponseTime: ResponseTimeStats{
				Sum:    hResponseTime.Sum(),
				Max:    hResponseTime.Max(),
				Mean:   hResponseTime.Mean(),
				Min:    hResponseTime.Min(),
				P50:    hResponseTime.Percentile(50),
				P75:    hResponseTime.Percentile(75),
				P95:    hResponseTime.Percentile(95),
				P99:    hResponseTime.Percentile(99),
				StdDev: hResponseTime.StdDev(),
				Var:    hResponseTime.Variance(),
			},
			Bytes: BytesSummary{
				Sent: BytesStats{
					Sum:    hBytesSent.Sum(),
					Max:    hBytesSent.Max(),
					Mean:   hBytesSent.Mean(),
					Min:    hBytesSent.Min(),
					P50:    hBytesSent.Percentile(50),
					P75:    hBytesSent.Percentile(75),
					P95:    hBytesSent.Percentile(95),
					P99:    hBytesSent.Percentile(99),
					StdDev: hBytesSent.StdDev(),
					Var:    hBytesSent.Variance(),
					Rate:   mBytesSent.RateMean(),
				},
				Received: BytesStats{
					Sum:    hBytesReceived.Sum(),
					Max:    hBytesReceived.Max(),
					Mean:   hBytesReceived.Mean(),
					Min:    hBytesReceived.Min(),
					P50:    hBytesReceived.Percentile(50),
					P75:    hBytesReceived.Percentile(75),
					P95:    hBytesReceived.Percentile(95),
					P99:    hBytesReceived.Percentile(99),
					StdDev: hBytesReceived.StdDev(),
					Var:    hBytesReceived.Variance(),
					Rate:   mBytesReceived.RateMean(),
				},
			},
		},
	}

	yamlOutput, err := yaml.Marshal(&output)
	check(err)
	err = ioutil.WriteFile(summaryPath, yamlOutput, 0644)
	check(err)
}
