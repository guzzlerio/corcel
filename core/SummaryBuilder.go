package core

//SummaryBuilder ...
type SummaryBuilder interface {
	Write(summary ExecutionSummary)
}
