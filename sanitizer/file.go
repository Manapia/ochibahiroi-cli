package sanitizer

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func SanitizeFilePath(filePath string) string {
	if runtime.GOOS == "windows" {
		return onWindows(filePath)
	}
	return filePath
}

func onWindows(filePath string) string {
	filePath = filepath.FromSlash(filePath)

	// For using git bash on Windows
	if filePath[0] == filepath.Separator {
		driveLetter := filePath[1]
		filePath = string(driveLetter) + ":" + filePath[2:]
	}

	var result string
	parts := strings.Split(filePath, string(filepath.Separator))

	if len(parts[0]) == 2 && parts[0][0] >= 65 && parts[0][0] <= 122 && parts[0][1] == ':' {
		result = parts[0] + string(filepath.Separator)
		parts = parts[1:]
	}

	for i, part := range parts {
		part = strings.ReplaceAll(part, `/`, "_")
		part = strings.ReplaceAll(part, `?`, "_")
		part = strings.ReplaceAll(part, `<`, "_")
		part = strings.ReplaceAll(part, `>`, "_")
		part = strings.ReplaceAll(part, `\`, "_")
		if i != 0 {
			part = strings.ReplaceAll(part, `:`, "_")
		}
		part = strings.ReplaceAll(part, `*`, "_")
		part = strings.ReplaceAll(part, `|`, "_")
		part = strings.ReplaceAll(part, `"`, "_")
		part = sanitizeReservedWord(part, "con")
		part = sanitizeReservedWord(part, "nul")
		part = sanitizeReservedWord(part, "prn")
		for j := 0; j <= 9; j++ {
			js := strconv.Itoa(j)
			part = sanitizeReservedWord(part, "com"+js)
			part = sanitizeReservedWord(part, "lpt"+js)
		}

		result = filepath.Join(result, part)
	}

	return result
}

func sanitizeReservedWord(src, word string) string {
	if strings.ToLower(src) != strings.ToLower(word) {
		return src
	}
	return "_" + src
}
