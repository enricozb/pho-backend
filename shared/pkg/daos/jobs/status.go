package jobs

type Status string

const (
	StatusScan    Status = "SCANNING"
	StatusConvert Status = "CONVERTING"
	StatusMove    Status = "MOVING"
)
