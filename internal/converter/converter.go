package converter

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"net/http"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Convert(r io.Reader, c echo.Context, opts ConvertOptions) error {
	if opts.OutputFormat != FormatJPEG && opts.OutputFormat != FormatPNG && opts.OutputFormat != FormatWebP {
		return ErrUnsupportedFormat
	}
	if (opts.OutputFormat == FormatJPEG || opts.OutputFormat == FormatWebP) && (opts.Quality < 1 || opts.Quality > 100) {
		return ErrInvalidQuality
	}

	img, _, err := image.Decode(r)
	if err != nil {
		slog.Error("Error while decoding Image.")
		return err
	}

	switch opts.OutputFormat {
	case FormatJPEG:
		c.Response().Writer.Header().Set("Content-Type", "image/jpeg")
		return jpeg.Encode(c.Response().Writer, img, &jpeg.Options{Quality: opts.Quality})
	case FormatPNG:
		c.Response().Writer.Header().Set("Content-Type", "image/png")
		return png.Encode(c.Response().Writer, img)
	case FormatWebP:
		c.Response().Writer.Header().Set("Content-Type", "image/webp")
		return webp.Encode(c.Response().Writer, img, &webp.Options{Quality: float32(opts.Quality), Lossless: false})
	}

	return nil
}

func SquareCrop(r io.Reader, c echo.Context) error {
	img, _, err := image.Decode(r)
	if err != nil {
		slog.Error("Error while opening image", "err", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid image file")
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	cropSize := min(width, height)
	cropped := imaging.CropCenter(img, cropSize, cropSize)

	c.Response().Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(c.Response().Writer, cropped, &jpeg.Options{Quality: 90})
	if err != nil {
		slog.Error("Error encoding cropped image", "err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to return image")
	}

	return nil
}

func FitToSquare(r io.Reader, c echo.Context, opts ConvertOptions) error {
	img, format, err := image.Decode(r)
	if err != nil {
		slog.Error("Error decoding image", "err", err, "format", format)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid image file")
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	frameSize := max(width, height)

	dst := imaging.New(frameSize, frameSize, color.NRGBA{255, 255, 255, 255})
	offsetX := (frameSize - width) / 2
	offsetY := (frameSize - height) / 2
	dst = imaging.Paste(dst, img, image.Pt(offsetX, offsetY))

	contentType := "image/" + string(format)
	c.Response().Header().Set("Content-Type", contentType)

	switch opts.OutputFormat {
	case FormatJPEG:
		slog.Info("Returning Square fitted image")
		err = jpeg.Encode(c.Response().Writer, dst, &jpeg.Options{Quality: opts.Quality})
	case FormatPNG:
		slog.Info("Returning Square fitted image")
		err = png.Encode(c.Response().Writer, dst)
	case FormatWebP:
		slog.Info("Returning Square fitted image")
		err = webp.Encode(c.Response().Writer, dst, &webp.Options{Quality: float32(opts.Quality)})
	}
	if err != nil {
		slog.Error("Error encoding image", "err", err, "format", opts.OutputFormat)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode image")
	}
	return nil
}

func ApplyFilter(r io.Reader, c echo.Context, filter FilterSettings) error {
	img, _, err := image.Decode(r)
	if err != nil {
		slog.Error("Error while opening image", "err", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid image file")
	}

	var result image.Image

	switch filter.Name {
	case Blur:
		result = imaging.Blur(img, float64(filter.Intensity))

	case Grayscale:
		res := imaging.Grayscale(img)
		result = imaging.AdjustContrast(res, float64(filter.Intensity))

	default:
		slog.Warn("Unknown filter", "filter", filter.Name)
		return echo.NewHTTPError(http.StatusBadRequest, "Unknown filter name")
	}

	c.Response().Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(c.Response().Writer, result, &jpeg.Options{Quality: 90})
	if err != nil {
		slog.Error("Error encoding image", "err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to return image")
	}

	return nil
}

func Invert(r io.Reader, c echo.Context) error {
	img, _, err := image.Decode(r)
	if err != nil {
		slog.Error("Error while opening image", "err", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid image file")
	}

	result := imaging.Invert(img)

	c.Response().Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(c.Response().Writer, result, &jpeg.Options{Quality: 90})
	if err != nil {
		slog.Error("Error encoding Inverted image", "err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to return image")
	}
	return nil
}

func MakePFP(r io.Reader, c echo.Context) error {
	img, _, err := image.Decode(r)
	if err != nil {
		slog.Error("Error while opening image", "err", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid image file")
	}
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	frameSize := max(width, height)

	dst := imaging.New(frameSize, frameSize, color.NRGBA{255, 255, 255, 255})
	offsetX := (frameSize - width) / 2
	offsetY := (frameSize - height) / 2
	dst = imaging.Paste(dst, img, image.Pt(offsetX, offsetY))

	dst = imaging.Resize(dst, 400, 400, imaging.Lanczos)

	c.Response().Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(c.Response().Writer, dst, &jpeg.Options{Quality: 90})
	if err != nil {
		slog.Error("Error encoding image", "err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to return image")
	}
	return nil
}
