package dataset

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	ThumbnailSmall  = 150
	ThumbnailMedium = 400
	ThumbnailBig    = 800
)

type ThumbnailSize int

const (
	SizeSmall ThumbnailSize = iota
	SizeMedium
	SizeBig
)

func (s ThumbnailSize) String() string {
	switch s {
	case SizeSmall:
		return "small"
	case SizeMedium:
		return "medium"
	case SizeBig:
		return "big"
	default:
		return "small"
	}
}

func (s ThumbnailSize) MaxDimension() int {
	switch s {
	case SizeSmall:
		return ThumbnailSmall
	case SizeMedium:
		return ThumbnailMedium
	case SizeBig:
		return ThumbnailBig
	default:
		return ThumbnailSmall
	}
}

func ParseThumbnailSize(s string) ThumbnailSize {
	switch strings.ToLower(s) {
	case "small":
		return SizeSmall
	case "medium":
		return SizeMedium
	case "big":
		return SizeBig
	default:
		return SizeSmall
	}
}

func resizeImage(src image.Image, maxDim int) image.Image {
	bounds := src.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	var newWidth, newHeight int
	if width > height {
		newWidth = maxDim
		newHeight = (height * maxDim) / width
	} else {
		newHeight = maxDim
		newWidth = (width * maxDim) / height
	}

	if newWidth >= width && newHeight >= height {
		return src
	}

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := (x * width) / newWidth
			srcY := (y * height) / newHeight
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}

	return dst
}

func GenerateThumbnail(srcPath, dstPath string, size ThumbnailSize) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("could not open source image: %w", err)
	}
	defer src.Close()

	img, format, err := image.Decode(src)
	if err != nil {
		return fmt.Errorf("could not decode image: %w", err)
	}

	thumbnail := resizeImage(img, size.MaxDimension())

	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("could not create destination directory: %w", err)
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}
	defer dst.Close()

	switch format {
	case "jpeg", "jpg":
		if err := jpeg.Encode(dst, thumbnail, &jpeg.Options{Quality: 85}); err != nil {
			return fmt.Errorf("could not encode JPEG: %w", err)
		}
	case "png":
		if err := png.Encode(dst, thumbnail); err != nil {
			return fmt.Errorf("could not encode PNG: %w", err)
		}
	default:
		if err := jpeg.Encode(dst, thumbnail, &jpeg.Options{Quality: 85}); err != nil {
			return fmt.Errorf("could not encode JPEG: %w", err)
		}
	}

	return nil
}

func CopyOrGenerateThumbnail(srcPath, dstPath string, size ThumbnailSize) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("could not open source image: %w", err)
	}
	defer src.Close()

	cfg, _, err := image.DecodeConfig(src)
	if err != nil {
		return fmt.Errorf("could not decode image config: %w", err)
	}

	maxDim := size.MaxDimension()
	needsResize := cfg.Width > maxDim || cfg.Height > maxDim

	if !needsResize {
		if _, err := src.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("could not seek to start: %w", err)
		}

		dstDir := filepath.Dir(dstPath)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return fmt.Errorf("could not create destination directory: %w", err)
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			return fmt.Errorf("could not create destination file: %w", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return fmt.Errorf("could not copy file: %w", err)
		}
		return nil
	}

	return GenerateThumbnail(srcPath, dstPath, size)
}

func GetThumbnailPath(privateDir, docRef string, attachmentIndex int, size ThumbnailSize) string {
	ext := ".jpg"
	filename := fmt.Sprintf("%s_%d_%s%s", docRef, attachmentIndex, size.String(), ext)
	return filepath.Join(privateDir, "pictures", filename)
}

func GetAttachmentPath(privateDir, docRef string, attachmentIndex int, originalExt string) string {
	if originalExt == "" {
		originalExt = ".jpg"
	}
	if !strings.HasPrefix(originalExt, ".") {
		originalExt = "." + originalExt
	}
	filename := fmt.Sprintf("%s_%d%s", docRef, attachmentIndex, originalExt)
	return filepath.Join(privateDir, "pictures", filename)
}
