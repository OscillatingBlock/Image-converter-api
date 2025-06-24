package api

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"imageConverter/internal/converter"
)

func ResponseWithError(c echo.Context, ErrorResponseCode int, message string, err error) error {
	slog.Error(message, "err", err)
	return c.JSON(ErrorResponseCode, map[string]string{
		"error": message,
	})
}

func processImageFile(c echo.Context) (*os.File, func(), error) {
	img, err := c.FormFile("file")
	if err != nil {
		return nil, nil, ResponseWithError(c, http.StatusBadRequest, "Error receiving image file", err)
	}

	src, err := img.Open()
	if err != nil {
		return nil, nil, ResponseWithError(c, http.StatusBadRequest, "Error opening image file", err)
	}

	tmpFile, err := os.CreateTemp("", "image-"+img.Filename)
	if err != nil {
		src.Close()
		return nil, nil, ResponseWithError(c, http.StatusBadRequest, "Error creating temp file", err)
	}

	if _, err = io.Copy(tmpFile, src); err != nil {
		src.Close()
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, nil, ResponseWithError(c, http.StatusBadRequest, "Error copying file", err)
	}
	src.Close()

	tmpReader, err := os.Open(tmpFile.Name())
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, nil, ResponseWithError(c, http.StatusBadRequest, "Error opening temp file", err)
	}

	cleanup := func() {
		tmpReader.Close()
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return tmpReader, cleanup, nil
}

func Convert(c echo.Context) error {
	tmpReader, cleanup, err := processImageFile(c)
	if err != nil {
		return err
	}
	defer cleanup()

	outputFormat := converter.OutputFormat(strings.ToLower(c.FormValue("output_format")))
	if outputFormat == "" {
		return ResponseWithError(c, http.StatusBadRequest, "Output format is required", nil)
	}

	qualityStr := c.FormValue("quality")
	quality := 80
	if qualityStr != "" {
		q, err := strconv.Atoi(qualityStr)
		if err != nil {
			return ResponseWithError(c, http.StatusBadRequest, "Error converting quality string to int", err)
		}
		quality = q
	}

	opts := converter.ConvertOptions{
		OutputFormat: outputFormat,
		Quality:      quality,
	}
	if err := converter.Convert(tmpReader, c, opts); err != nil {
		return ResponseWithError(c, http.StatusInternalServerError, "Error while converting image", err)
	}

	return nil
}

func SquareCropHandler(c echo.Context) error {
	tmpReader, cleanup, err := processImageFile(c)
	if err != nil {
		return err
	}
	defer cleanup()

	if err := converter.SquareCrop(tmpReader, c); err != nil {
		return ResponseWithError(c, http.StatusInternalServerError, "Error while cropping image", err)
	}

	return nil
}

func FitToSquareHandler(c echo.Context) error {
	tmpReader, cleanup, err := processImageFile(c)
	if err != nil {
		return err
	}
	defer cleanup()

	outputFormat := converter.OutputFormat(strings.ToLower(c.FormValue("output_format")))
	if outputFormat == "" {
		return ResponseWithError(c, http.StatusBadRequest, "Output format is required", nil)
	}
	if outputFormat != converter.FormatJPEG && outputFormat != converter.FormatPNG && outputFormat != converter.FormatWebP {
		return ResponseWithError(c, http.StatusBadRequest, "Unsupported output format", nil)
	}

	qualityStr := c.FormValue("quality")
	quality := 80
	if qualityStr != "" {
		q, err := strconv.Atoi(qualityStr)
		if err != nil {
			return ResponseWithError(c, http.StatusBadRequest, "Error converting quality string to int", err)
		}
		quality = q
	}

	opts := converter.ConvertOptions{
		OutputFormat: outputFormat,
		Quality:      quality,
	}

	if err := converter.FitToSquare(tmpReader, c, opts); err != nil {
		return ResponseWithError(c, http.StatusInternalServerError, "Error while fitting to square", err)
	}

	return nil
}

func InvertHandler(c echo.Context) error {
	tmpReader, cleanup, err := processImageFile(c)
	if err != nil {
		return err
	}
	defer cleanup()

	if err := converter.Invert(tmpReader, c); err != nil {
		return ResponseWithError(c, http.StatusInternalServerError, "Error while inverting image", err)
	}

	return nil
}

func ApplyFilterHandler(c echo.Context) error {
	tmpReader, cleanup, err := processImageFile(c)
	if err != nil {
		return err
	}
	defer cleanup()

	filterName := converter.FilterName(strings.ToLower(c.FormValue("filter_name")))
	intensityStr := c.FormValue("intensity")
	intensity := 10

	if intensityStr != "" {
		i, err := strconv.Atoi(intensityStr)
		if err != nil {
			return ResponseWithError(c, http.StatusBadRequest, "Error converting intensity string to int", err)
		}
		intensity = i
	}

	filter := converter.FilterSettings{
		Name:      filterName,
		Intensity: intensity,
	}

	if err := converter.ApplyFilter(tmpReader, c, filter); err != nil {
		return ResponseWithError(c, http.StatusInternalServerError, "Error while applying filter", err)
	}

	return nil
}
