package processor

import (
	"net/http"

	"ci.guzzler.io/guzzler/corcel/config"
	"ci.guzzler.io/guzzler/corcel/request"
)

//PlanExecutor ...
type PlanExecutor struct {
	Config *config.Configuration
	Stats  *Statistics
	Bar    ProgressBar
}

func (instance *PlanExecutor) createPlan() Plan {
	plan := Plan{
		Name:     "Plan from urls in file",
		Workers:  instance.Config.Workers,
		WaitTime: instance.Config.WaitTime,
		Jobs:     []Job{},
	}

	job := Job{
		Name:  "Job for the urls in file",
		Steps: []Step{},
	}

	reader := request.NewRequestReader(instance.Config.FilePath)

	stream := request.NewSequentialRequestStream(reader)

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 50,
		},
	}

	for stream.HasNext() {
		request, _ := stream.Next()
		step := Step{}

		action := &HTTPRequestExecutionAction{
			Client:  client,
			URL:     request.URL.String(),
			Method:  request.Method,
			Headers: request.Header,
		}

		step.Action = action
		job.Steps = append(job.Steps, step)
	}

	plan.Jobs = append(plan.Jobs, job)

	return plan
}

// Execute ...
func (instance *PlanExecutor) Execute() error {
	plan := instance.createPlan()
	return plan.Execute()
}
