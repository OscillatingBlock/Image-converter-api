package converter

import "errors"

type OutputFormat string

const (
	FormatJPEG OutputFormat = "jpeg"
	FormatPNG  OutputFormat = "png"
	FormatWebP OutputFormat = "webp"
)

type ConvertOptions struct {
	OutputFormat OutputFormat
	Quality      int
}

var (
	ErrUnsupportedFormat = errors.New("unsupported output format")
	ErrInvalidQuality    = errors.New("quality must be between 1 and 100")
)

type FilterName string

const (
	Blur      FilterName = "blur"
	Grayscale FilterName = "grayscale"
)

type FilterSettings struct {
	Name      FilterName
	Intensity int
}
