package converter

type Converter interface {
	// Convert copies src to dst, converting if necessary.
	// This call alone might not do any copying, however after Finish is called this conversion will have completed.
	Convert(src, dst string) error

	// Complete any remaining conversion tasks, blocking until all are done.
	Finish() error
}

var converters = map[string]Converter{
	"image/png":  NewIdentityConverter(),
	"image/jpeg": NewIdentityConverter(),
	"image/heic": nil,

	"video/quicktime": nil,
}
