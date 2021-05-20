package jobs

type JobKind string

const (
	JobScan              JobKind = "SCAN"
	JobMetadata          JobKind = "METADATA"
	JobMetadataHash      JobKind = "METADATA_HASH"
	JobMetadataTimestamp JobKind = "METADATA_TIMESTAMP"
	JobMetadataLive      JobKind = "METADATA_LIVE"
	JobMetadataEXIF      JobKind = "METADATA_EXIF"
	JobMetadataMonitor   JobKind = "METADATA_MONITOR"
	JobDedupe            JobKind = "DEDUPE"
	JobConvert           JobKind = "CONVERT"
	JobCleanup           JobKind = "CLEANUP"
)

var MetadataJobKinds = []JobKind{
	JobMetadataHash,
	JobMetadataEXIF,
}
