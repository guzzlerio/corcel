package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/guzzlerio/rizo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"ci.guzzler.io/guzzler/corcel/global"
	"ci.guzzler.io/guzzler/corcel/processor"
	"ci.guzzler.io/guzzler/corcel/statistics"
	. "ci.guzzler.io/guzzler/corcel/utils"
)

func GetPeopleTestRequest() map[string]interface{} {
	return map[string]interface{}{
		"type":          "HttpRequest",
		"requesTimeout": 150,
		"method":        "GET",
		"url":           TestServer.CreateURL("/people"),
		"httpHeaders": map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func HTTPStatusExactAssertion(code int) map[string]interface{} {
	return map[string]interface{}{
		"type":     "ExactAssertion",
		"key":      "http:response:status",
		"expected": code,
	}
}

func ExecutePlanBuilder(planBuilder *processor.YamlPlanBuilder) error {
	file, err := planBuilder.Build()
	if err != nil {
		return err
	}
	exePath, err := filepath.Abs("./corcel")
	if err != nil {
		return err
	}
	defer func() {
		err := os.Remove(file.Name())
		if err != nil {
			panic(err)
		}
	}()
	args := []string{"--plan"}
	cmd := exec.Command(exePath, append(args, file.Name())...)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	return err
}

var _ = Describe("ExecutionPlan", func() {

	BeforeEach(func() {
		TestServer.Clear()
		factory := rizo.HTTPResponseFactory(func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
		})

		//Add more methods to the TestServer so that there is not a need
		//to explicitly create a factory.  It would be good to be able to
		//do TestServer.ReturnHttpStatus(200).For(rizo.RequestWithPath("/people"))
		TestServer.Use(factory).For(rizo.RequestWithPath("/people"))
	})

	AfterEach(func() {
		TestServer.Clear()
	})

	/*/\/\/\/\/\/\/\//\/\/\/\/\/\/\/\/\/\/\/\\/\/
	 *
	 * REFER TO ISSUE #50
	 *
	 *\/\/\/\/\/\/\//\/\/\/\/\/\/\/\/\/\/\/\\/\/
	 */

	for _, numberOfWorkers := range global.NumberOfWorkersToTest {
		name := fmt.Sprintf("SetWorkers for %v workers", numberOfWorkers)
		func(workers int) {
			It(name, func() {
				planBuilder := processor.NewYamlPlanBuilder()

				planBuilder.
					SetWorkers(workers).
					CreateJob().
					CreateStep().
					ToExecuteAction(GetPeopleTestRequest())

				err := ExecutePlanBuilder(planBuilder)
				Expect(err).To(BeNil())

				var executionOutput statistics.AggregatorSnapShot
				UnmarshalYamlFromFile("./output.yml", &executionOutput)
				var summary = statistics.CreateSummary(executionOutput)

				Expect(summary.TotalErrors).To(Equal(float64(0)))
				Expect(summary.TotalRequests).To(Equal(float64(workers)))
				Expect(len(TestServer.Requests)).To(Equal(workers))
			})
		}(numberOfWorkers)
	}

	It("SetWaitTime", func() {
		numberOfSteps := 6
		waitTime := "500ms"

		planBuilder := processor.NewYamlPlanBuilder()
		planBuilder.SetWaitTime(waitTime)
		jobBuilder := planBuilder.CreateJob()

		for i := 0; i < numberOfSteps; i++ {
			jobBuilder.CreateStep().ToExecuteAction(GetPeopleTestRequest())
		}

		err := ExecutePlanBuilder(planBuilder)
		Expect(err).To(BeNil())

		var executionOutput statistics.AggregatorSnapShot
		UnmarshalYamlFromFile("./output.yml", &executionOutput)
		var summary = statistics.CreateSummary(executionOutput)

		actual, _ := time.ParseDuration(summary.RunningTime)
		seconds := actual.Seconds()
		seconds = math.Floor(seconds)
		Expect(seconds).To(Equal(float64(3)))
	})

	PIt("Wait Time")
	PIt("Duration")
	PIt("Random")

	PDescribe("HttpRequest", func() {})

	PDescribe("Assertions", func() {

		//ASSERTION FAILURES ARE NOT CURRENTLY COUNTING AS ERRORS IN THE SUMMARY OUTPUT
		PIt("ExactAssertion")
	})

	It("Name", func() {

		planBuilder := processor.NewYamlPlanBuilder()
		planBuilder.CreateJob().
			CreateStep().
			ToExecuteAction(GetPeopleTestRequest()).
			WithAssertion(HTTPStatusExactAssertion(201))

		err := ExecutePlanBuilder(planBuilder)
		CheckErr(err)

		var executionOutput statistics.AggregatorSnapShot
		UnmarshalYamlFromFile("./output.yml", &executionOutput)
		var summary = statistics.CreateSummary(executionOutput)

		Expect(summary.TotalRequests).To(BeNumerically(">", 0))
		Expect(summary.TotalErrors).To(Equal(float64(0)))
	})
})
