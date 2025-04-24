package content

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"
)

const MAX_SIZE = 5 * 1024 * 1024

func CheckImage(fileBytes []byte) (string, error) {
	if len(fileBytes) > MAX_SIZE {
		return "", fmt.Errorf("файл слишком большой")
	}

	contentType := http.DetectContentType(fileBytes)

	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("неверный формат файла")
	}

	validFormats := map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
	}
	if !validFormats[contentType] {
		return "", fmt.Errorf("не поддерживаемый формат изображения")
	}

	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return "", fmt.Errorf("не удалось декодировать изображение: %w", err)
	}
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	if width < 100 || height < 100 {
		return "", fmt.Errorf("изображение слишком маленькое (минимум 100x100)")
	}
	if width > 4000 || height > 4000 {
		return "", fmt.Errorf("изображение слишком большое (максимум 4000x4000)")
	}

	return contentType, nil
}
