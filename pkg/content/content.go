package content

import (
	"fmt"
	"net/http"
	"strings"
)

const MAX_SIZE = 5 * 1024 * 1024

func CheckImage(fileBytes []byte) (string, error){
	if (len(fileBytes) > MAX_SIZE) {
		return "", fmt.Errorf("file too big")
	}

	contentType := http.DetectContentType(fileBytes)

	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("file is not an image")
	}

	validFormats := map[string]string {"image/png": "png", "image/jpeg": "jpeg"}
	ext, ok := validFormats[contentType]
	if !ok {
		return "", fmt.Errorf("not supported file type")
	}

	return ext, nil
}