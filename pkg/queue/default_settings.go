package queue

var defaultSettings = DefaultSettings{
	Parallelism:                  1,
	TtlAfterSuccessfulCompletion: "",
	TtlAfterFailedCompletion:     "",
	ExecutionTimeout:             "",
	MaxTimeInQueue:               "",
}

type DefaultSettings struct {
	Parallelism                  int64
	TtlAfterSuccessfulCompletion string
	TtlAfterFailedCompletion     string
	ExecutionTimeout             string
	MaxTimeInQueue               string
}

func (d DefaultSettings) GetParallelism() int64 {
	return d.Parallelism
}

func (d DefaultSettings) GetTtlAfterSuccessfulCompletion() string {
	return d.TtlAfterSuccessfulCompletion
}

func (d DefaultSettings) GetTtlAfterFailedCompletion() string {
	return d.TtlAfterFailedCompletion
}

func (d DefaultSettings) GetExecutionTimeout() string {
	return d.ExecutionTimeout
}

func (d DefaultSettings) GetMaxTimeInQueue() string {
	return d.MaxTimeInQueue
}
