package shortcode

import (
      "crypto/sha256"
      "encoding/base64"
      "strings"
)

func Generate(input string) string {
      hash := sha256.Sum256([]byte(input))
      return strings.TrimRight(base64.URLEncoding.EncodeToString(hash[:6]), "=")
}