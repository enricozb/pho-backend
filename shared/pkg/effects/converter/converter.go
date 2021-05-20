package converter

import "fmt"

type Converter interface {
	// Convert copies src to dst, converting if necessary.
	// This call alone might not do any copying, however after Finish is called this conversion will have completed.
	Convert(src, dst string) error

	// Complete any remaining conversion tasks, blocking until all are done.
	Finish() error
}

var SupportedMimeTypes []string

var converters = map[string]func() Converter{}

func registerConverter(mimetype string, c func() Converter) {
	if _, alreadyRegistered := converters[mimetype]; alreadyRegistered {
		panic(fmt.Sprintf("converter already exists for mimetype %s", mimetype))
	}
	converters[mimetype] = c

	SupportedMimeTypes = append(SupportedMimeTypes, mimetype)
}
