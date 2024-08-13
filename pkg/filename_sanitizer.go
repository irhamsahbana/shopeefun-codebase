package pkg

import (
	"strings"

	"github.com/oklog/ulid/v2"
)

func SanitizeFilename(s string, makeUnique bool) (result string) {
	// Trim leading and trailing spaces
	s = strings.TrimSpace(s)

	// Replace spaces with underscores
	result = strings.ReplaceAll(s, " ", "_")

	// Replace other characters with underscores
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "#", "%", "&", "{", "}", "^", "~", "[", "]", "(", ")", "`"}
	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}

	// Make the filename unique
	if makeUnique {
		// check if the filename already has an extension
		if strings.Contains(result, ".") {
			// split the filename and extension
			parts := strings.Split(result, ".")

			ext := parts[len(parts)-1]
			// remove the extension from the filename
			parts = parts[:len(parts)-1]

			// append a unique string to the filename
			result = strings.Join(parts, ".") + "_" + ulid.Make().String() + "." + ext
		} else {
			// append a unique string to the filename
			result = result + "_" + ulid.Make().String()
		}
	}

	return result
}
