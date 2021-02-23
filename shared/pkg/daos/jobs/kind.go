package jobs

type JobKind string

const (
	JobScan    JobKind = "SCAN"
	JobConvert JobKind = "CONVERT"
	JobMove    JobKind = "MOVE"
)
