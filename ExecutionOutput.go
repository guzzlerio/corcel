package main

type ExecutionSummary struct{
	TotalBytesSent uint64 `yaml:"totalBytesSent"`
}

type ExecutionOutput struct{
	Summary ExecutionSummary `yaml:"summary"`
}
