package processor

import (
	"time"

	"github.com/guzzlerio/corcel/config"
	"github.com/guzzlerio/corcel/core"
)

//CreateJobStream ...
func CreateJobStream(jobs []core.Job, config *config.Configuration) JobStream {

	var jobStream JobStream
	jobStream = CreateJobSequentialStream(jobs)

	if config.Random {
		jobStream = CreateJobRandomStream(jobs)
	}

	if config.Iterations > 0 {
		revolvingStream := CreateJobRevolvingStream(jobStream)
		iterationStream := CreateJobIterationStream(*revolvingStream, len(jobs), config.Iterations)
		jobStream = iterationStream
	}

	if config.Duration > time.Duration(0) {
		jobStream = CreateJobDurationStream(jobStream, config.Duration)
	}

	return jobStream
}
