package jobs

type Status string

const (
	StatusNotStarted Status = "NOT_STARTED"
	StatusScan       Status = "SCAN"
	StatusMetadata   Status = "METADATA"
	StatusDedupe     Status = "DEDUPE"
	StatusConvert    Status = "CONVERT"
	StatusCleanup    Status = "CLEANUP"
	StatusDone       Status = "DONE"
	StatusFailed     Status = "FAILED"
)
