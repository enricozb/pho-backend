package paths

type EXIFMetadata struct {
	Path              string `json:"SourceFile"`
	CreateDate        int64  `json:"CreateDate"`
	MediaGroupUUID    string `json:"MediaGroupUUID"`
	ImageUniqueID     string `json:"ImageUniqueID"`
	ContentIdentifier string `json:"ContentIdentifier"`
}
