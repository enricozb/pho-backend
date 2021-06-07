package jobs

type JobKind string

const (
	JobScan            JobKind = "SCAN"
	JobMetadata        JobKind = "METADATA"
	JobMetadataHash    JobKind = "METADATA_HASH"
	JobMetadataEXIF    JobKind = "METADATA_EXIF"
	JobMetadataMonitor JobKind = "METADATA_MONITOR"
	JobDedupe          JobKind = "DEDUPE"
	JobConvert         JobKind = "CONVERT"
	JobThumbnail       JobKind = "THUMBNAIL"
	JobCleanup         JobKind = "CLEANUP"
)

var MetadataJobKinds = []JobKind{
	JobMetadataHash,
	JobMetadataEXIF,
}
