package exif

type EXIFMetadata struct {
	Path      string `json:"SourceFile"`
	Timestamp string `json:"Timestamp"`

	Width  int `json:"ImageWidth"`
	Height int `json:"ImageHeight"`

	LiveID []byte
}
